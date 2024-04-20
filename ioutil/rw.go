package ioutil

import "io"

type Read interface {
	Read(io.Reader) error
}

type Write interface {
	Write(io.Writer) error
}

type ReadWrite interface {
	Read
	Write
}
