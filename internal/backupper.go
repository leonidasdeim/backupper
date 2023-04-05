package internal

import (
	"errors"
	"fmt"
	"path"
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
	if err := Utils.CreateFolder(path); err != nil {
		return nil, fmt.Errorf("can't create backup directory: %v", err)
	}

	return &Backupper{
		directory: path,
		log:       log,
	}, nil
}

func (b *Backupper) FileCreated(file string) {
	if !Utils.IsFile(file) {
		return
	}
	filename := filepath.Base(file)

	if strings.HasPrefix(filepath.Base(file), deletePrefix) {
		if err := b.delete(file); err != nil {
			b.log.Error(filename)
		}
		return
	}

	b.log.Created(filename)
	if err := b.backup(file); err != nil {
		b.log.Error(filename)
	}
}

func (b *Backupper) FileModified(file string) {
	if !Utils.IsFile(file) {
		return
	}
	filename := filepath.Base(file)

	b.log.Modified(filename)
	if err := b.backup(file); err != nil {
		b.log.Error(filename)
	}
}

func (b *Backupper) backup(file string) error {
	backupFile := path.Join(b.directory, filepath.Base(file)+extension)

	if err := Utils.CopyFile(file, backupFile); err != nil {
		return errors.New("copy error")
	}

	b.log.Backup(filepath.Base(file))
	return nil
}

func (b *Backupper) delete(file string) error {
	originalName := strings.TrimPrefix(filepath.Base(file), deletePrefix)
	backupFile := path.Join(b.directory, originalName+extension)

	if err := Utils.DeleteFile(file); err != nil {
		return errors.New("delete error")
	}

	if err := Utils.DeleteFile(backupFile); err != nil {
		return errors.New("backup delete error")
	}

	b.log.Deleted(originalName)
	return nil
}
