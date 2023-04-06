package internal

import (
	"errors"
	"fmt"

	"github.com/fsnotify/fsnotify"
)

type Consumer interface {
	FileCreated(string)
	FileModified(string)
}

type Notifier struct {
	consumer Consumer
	watcher  *fsnotify.Watcher
}

// Initializes Notifier object for provided directory path
// Provided consumer will be notified on file events
func NewNotifier(path string, consumer Consumer) (*Notifier, error) {
	if consumer == nil {
		return nil, errors.New("consumer object is not provided")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("can't create new watcher instance: %v", err)
	}

	if err = w.Add(path); err != nil {
		return nil, err
	}

	return &Notifier{
		consumer: consumer,
		watcher:  w,
	}, nil
}

// Runs files events watcher.
// This function is blocking, should be run from goroutine.
func (n *Notifier) Watch() {
	for {
		select {
		case event, ok := <-n.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Create) {
				go n.consumer.FileCreated(event.Name)
			} else if event.Has(fsnotify.Write) {
				go n.consumer.FileModified(event.Name)
			}
		case err, ok := <-n.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Notifier error: %v\n", err)
		}
	}
}

// Closes Notifier instance
func (n *Notifier) Close() {
	n.watcher.Close()
}
