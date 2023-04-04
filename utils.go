package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type _utils struct{}

var utils = _utils{}

func (_utils) CreateFolder(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

func (_utils) IsFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

func (_utils) CopyFile(in string, out string) error {
	fin, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("can't open a file: %v", err)
	}
	defer fin.Close()

	err = utils.CreateFolder(filepath.Dir(out))
	if err != nil {
		return fmt.Errorf("error while creating output directory: %v", err)
	}

	fout, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("can't create a new file: %v", err)
	}
	defer fout.Close()

	if _, err = io.Copy(fout, fin); err != nil {
		return fmt.Errorf("can't copy file contents: %v", err)
	}

	return nil
}

func (_utils) DeleteFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.Remove(path)
}
