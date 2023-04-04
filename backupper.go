package main

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
)

const backupExtension = ".bak"

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
	backupFile := path.Join(b.directory, filepath.Base(file)+backupExtension)

	if err := utils.CopyFile(file, backupFile); err != nil {
		log.Println(err.Error())
		return
	}
	log.Println("created/modified file:", backupFile)
}

func (b *Backupper) FileRenamed(file string) {
	log.Println("renamed file:", file)
}
