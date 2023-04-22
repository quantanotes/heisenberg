package main

import (
	"os"
	"path/filepath"
)

// Get full dir from local dir
func GetDir(path string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(wd, path), nil
}
