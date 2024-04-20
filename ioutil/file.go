package ioutil

import "io"

func Append(wr Write, s io.Seeker, w io.Writer) error {
	if _, err := s.Seek(0, io.SeekEnd); err != nil {
		return err
	}
	return wr.Write(w)
}
