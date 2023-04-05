package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/leonidasdeim/backupper/internal"
)

const (
	appName   = "backupper"
	logFile   = appName + ".log"
	stateFile = appName + ".tmp"
)

type AppProps struct {
	HotDir    string
	BackupDir string
}

func RunApp(props AppProps) {
	logger, err := internal.NewLogger(logFile)
	if err != nil {
		fmt.Println("Logger initialization error: ", err)
		return
	}
	defer logger.Close()

	backupper, err := internal.NewBackupper(props.BackupDir, logger)
	if err != nil {
		fmt.Println("Backupper initialization error: ", err)
		return
	}

	notifier, err := internal.NewNotifier(props.HotDir, backupper)
	if err != nil {
		fmt.Println("Notifier initialization error: ", err)
		return
	}
	defer notifier.Close()

	go RunLogFilter(logger, stateFile)
	go notifier.Watch()
	fmt.Printf("Initialization successful. Backing up directory: [%s], to: [%s]\n",
		props.HotDir,
		props.BackupDir,
	)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	fmt.Println("Gracefully closing application")
}
