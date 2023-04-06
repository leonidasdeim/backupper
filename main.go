package main

import (
	"flag"

	"github.com/leonidasdeim/backupper/cmd"
	"github.com/leonidasdeim/backupper/internal"
)

const (
	dirFlag    = "in"
	backupFlag = "out"
)

var _ cmd.AppState = (*internal.State)(nil)
var _ cmd.FilterState = (*internal.State)(nil)

func main() {
	hot := flag.String(dirFlag, "", "path to hot directory")
	backup := flag.String(backupFlag, "", "path for backup directory")
	flag.Parse()

	state := internal.LoadState(cmd.StateFile)

	if isFlagPresent(dirFlag) && isFlagPresent(backupFlag) {
		state.SetDirectories(
			internal.Directories{
				Hot:    *hot,
				Backup: *backup,
			},
		).Save()
	}

	cmd.RunApp(&state)
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
