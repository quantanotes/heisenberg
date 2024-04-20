package state

import (
	"io"
	"os"
)

type File interface {
	io.Reader
	io.Writer
	io.Seeker
	io.Closer
	Sync() error
	Size() (int, error)
}

type LocalFile struct {
	*os.File
}

type RemoteFile struct{}

func (f *LocalFile) Size() (int, error) {
	stats, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return int(stats.Size()), nil
}
