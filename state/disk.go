package state

import (
	"os"
	"path/filepath"
)

type Disk interface {
	Open(path string) (File, error)
	Scan() ([]string, error)
}

type LocalDisk struct {
	root string
}

type RemoteDisk struct{}

func NewLocalDisk(root string) *LocalDisk {
	return &LocalDisk{root}
}

func (d *LocalDisk) Open(path string) (LocalFile, error) {
	f, err := os.OpenFile(d.path(path), os.O_RDWR|os.O_CREATE, 0666)
	return LocalFile{f}, err
}

func (d *LocalDisk) Scan() ([]string, error) {
	var files []string
	walkfun := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	}

	if err := filepath.Walk(d.path(d.root), walkfun); err != nil {
		return nil, err
	}

	return files, nil
}

func (d *LocalDisk) path(path string) string {
	return filepath.Join(d.root, path)
}
