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
	LogFile   = appName + ".log"
	StateFile = appName + ".tmp"
)

type AppState interface {
	SetState(state internal.Data)
	SetFilter(filter internal.Filter) *internal.State
	GetFilter() *internal.Filter
	SetDirectories(dirs internal.Directories) *internal.State
	GetDirectories() *internal.Directories
	Save() error
}

func RunApp(state AppState) {
	dirs := state.GetDirectories()
	if dirs == nil {
		fmt.Println("Can't start app: paths for hot/backup folders are not provided")
		return
	}

	logger, err := internal.NewLogger(LogFile)
	if err != nil {
		fmt.Println("Logger init error: ", err)
		return
	}
	defer logger.Close()

	backupper, err := internal.NewBackupper(dirs.Backup, logger)
	if err != nil {
		fmt.Println("Backupper init error: ", err)
		return
	}

	notifier, err := internal.NewNotifier(dirs.Hot, backupper)
	if err != nil {
		fmt.Println("Notifier init error: ", err)
		return
	}
	defer notifier.Close()

	go notifier.Watch()
	fmt.Printf("Backup in progress...\n- Hot folder: %s \n- Backup folder: %s\n\n",
		dirs.Hot,
		dirs.Backup,
	)
	go RunLogViewer(logger, state)

	// wait for SIGTERM or SIGINT to close application
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-interrupt

	fmt.Println("Gracefully closing application")
}
