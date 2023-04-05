package main

import (
	"flag"
	"fmt"

	"github.com/leonidasdeim/backupper/cmd"
)

// TESTS

const (
	dirFlag    = "in"
	backupFlag = "out"
)

func main() {
	dirPath := flag.String(dirFlag, "", "path to hot directory")
	backupPath := flag.String(backupFlag, "", "path for backup directory")
	flag.Parse()

	if !isFlagPresent(dirFlag) ||
		!isFlagPresent(backupFlag) {
		fmt.Println("Required arguments not provided: run with --help for usage")
		return
	}

	cmd.RunApp(cmd.AppProps{
		HotDir:    *dirPath,
		BackupDir: *backupPath,
	})
}

func isFlagPresent(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
			return
		}
	})
	return found
}
