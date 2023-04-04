package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"
)

const backupExtension = ".bak"
const deletePrefix = "delete_"

type Backupper struct {
	directory string
}

var _ Consumer = (*Backupper)(nil)

func NewBackupper(dir string) (*Backupper, error) {
	if err := utils.CreateFolder(dir); err != nil {
		return nil, fmt.Errorf("can't create backup directory: %v", err)
	} else {
		log.Println("backup folder created: ", dir)
	}

	return &Backupper{
		directory: dir,
	}, nil
}

func (b *Backupper) FileModified(file string) {
	if !utils.IsFile(file) {
		return
	}

	if strings.HasPrefix(filepath.Base(file), deletePrefix) {
		b.delete(file)
		return
	}

	backupFile := path.Join(b.directory, filepath.Base(file)+backupExtension)

	if err := utils.CopyFile(file, backupFile); err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("backup file:", filepath.Base(file))
}

func (b *Backupper) delete(file string) {
	originalName := strings.TrimPrefix(filepath.Base(file), deletePrefix)
	backupFile := path.Join(b.directory, originalName+backupExtension)

	if err := utils.DeleteFile(file); err != nil {
		log.Println("can't delete file: ", file, err)
	}

	if err := utils.DeleteFile(backupFile); err != nil {
		log.Println("can't delete file: ", backupFile, err)
	}

	log.Println("deleted file:", originalName)
}
