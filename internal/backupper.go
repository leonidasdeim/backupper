package internal

import (
	"errors"
	"fmt"
	pathUtils "path"
	"path/filepath"
	"strings"
)

const (
	extension    = ".bak"
	deletePrefix = "delete_"
)

type Backupper struct {
	directory string
	log       Log
}

var _ Consumer = (*Backupper)(nil)

// Return new Backupper instance. It implements Consumer interface.
// Receives path to backup directory and logger instance
func NewBackupper(path string, log Log) (*Backupper, error) {
	if log == nil {
		return nil, errors.New("logger object is not provided")
	}

	if err := Utils.CreateFolder(path); err != nil {
		return nil, fmt.Errorf("can't create backup directory: %v", err)
	}

	return &Backupper{
		directory: path,
		log:       log,
	}, nil
}

// Callback function for 'file created' event. Argument - path to the file
func (b *Backupper) FileCreated(path string) {
	if !Utils.IsFile(path) {
		return
	}
	filename := filepath.Base(path)

	if strings.HasPrefix(filepath.Base(path), deletePrefix) {
		if err := b.delete(path); err != nil {
			b.log.Error(filename)
		}
		return
	}

	b.log.Created(filename)
	if err := b.backup(path); err != nil {
		b.log.Error(filename)
	}
}

// Callback function for 'file modified' event. Argument - path to the file
func (b *Backupper) FileModified(path string) {
	if !Utils.IsFile(path) {
		return
	}
	filename := filepath.Base(path)

	b.log.Modified(filename)
	if err := b.backup(path); err != nil {
		b.log.Error(filename)
	}
}

func (b *Backupper) backup(path string) error {
	backupFile := pathUtils.Join(b.directory, filepath.Base(path)+extension)

	if err := Utils.CopyFile(path, backupFile); err != nil {
		return errors.New("copy error")
	}

	b.log.Backup(filepath.Base(path))
	return nil
}

func (b *Backupper) delete(path string) error {
	originalName := strings.TrimPrefix(filepath.Base(path), deletePrefix)
	backupFile := pathUtils.Join(b.directory, originalName+extension)

	if err := Utils.DeleteFile(path); err != nil {
		return errors.New("delete error")
	}

	if err := Utils.DeleteFile(backupFile); err != nil {
		return errors.New("backup delete error")
	}

	b.log.Deleted(originalName)
	return nil
}
