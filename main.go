package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// LOGGER

const hotFolderFlag = "in"
const backupFolderFlag = "out"

func main() {
	dir := flag.String(hotFolderFlag, "", "path to hot directory")
	backupDir := flag.String(backupFolderFlag, "", "path for backup directory")
	flag.Parse()

	if !isFlagPassed(hotFolderFlag) ||
		!isFlagPassed(backupFolderFlag) {
		fmt.Println("required arguments not provided")
		return
	}

	log.Println("application is starting")

	backupper, err := NewBackupper(*backupDir)
	if err != nil {
		log.Println(err)
		return
	}

	notifier, err := NewNotifier(*dir, backupper)
	if err != nil {
		log.Println(err)
		return
	}
	defer notifier.Close()

	go notifier.Watch()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt

	log.Println("gracefully closing application")
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
			return
		}
	})
	return found
}
