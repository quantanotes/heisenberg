package ioutil

import (
	"encoding/binary"
	"io"
)

func ReadBytes(r io.Reader, buf []byte) error {
	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil
	}

	if _, err := io.ReadFull(r, buf); err != nil {
		return nil
	}

	return nil
}

func WriteBytes(w io.Writer, buf []byte) error {
	sizebuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizebuf, uint32(len(buf)))
	if _, err := w.Write(sizebuf); err != nil {
		return err
	}

	if _, err := w.Write(buf); err != nil {
		return err
	}

	return nil
}

func AppendBytes(w io.Writer, s io.Seeker, buf []byte) error {
	if _, err := s.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	WriteBytes(w, buf)
	return nil
}
