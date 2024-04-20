package master

import (
	"bufio"
	"encoding/binary"
	"heisenberg/ioutil"
	"heisenberg/model"
	"heisenberg/state"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

type wal struct {
	mu      sync.Mutex
	file    state.File
	writer  *bufio.Writer
	logger  *slog.Logger
	syncer  ioutil.Syncer
	size    int
	maxSize int
}

type walEntry struct {
	method model.Method
	key    []byte
	value  []byte
}

func newWAL(file state.File) (*wal, error) {
	size, err := file.Size()
	if err != nil {
		return nil, err
	}

	h := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(h)

	w := &wal{
		file:   file,
		writer: bufio.NewWriter(file),
		logger: logger,
		size:   size,
	}

	w.syncer = ioutil.NewSyncer(w, time.Millisecond, logger)
	w.syncer.Start()

	return w, nil

}

func (w *wal) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.file.Close()
	w.syncer.Stop()
}

func (w *wal) Write(method model.Method, key, value []byte) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	we := newWALEntry(method, key, value)
	if err := ioutil.Append(&we, w.file, w.writer); err != nil {
		return err
	}
	w.size += we.Size()

	return nil
}

func (w *wal) Full() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.size >= w.maxSize
}

func (w *wal) Rotate() error {
	return nil
}

func (w *wal) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.sync()
}

func (w *wal) sync() error {
	if err := w.writer.Flush(); err != nil {
		return err
	}
	return w.file.Sync()
}

func newWALEntry(method model.Method, key, value []byte) walEntry {
	return walEntry{method, key, value}
}

func (we *walEntry) Read(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &we.method); err != nil {
		return err
	}

	if err := ioutil.ReadBytes(r, we.key); err != nil {
		return err
	}

	if err := ioutil.ReadBytes(r, we.value); err != nil {
		return err
	}

	return nil

}

func (we *walEntry) Write(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, we.method); err != nil {
		return err
	}

	if err := ioutil.WriteBytes(w, we.key); err != nil {
		return err
	}

	if err := ioutil.WriteBytes(w, we.value); err != nil {
		return err
	}

	return nil
}

func (we *walEntry) Size() int {
	return 12 + len(we.key) + len(we.value)
}
