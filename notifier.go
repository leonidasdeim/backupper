package main

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

type Consumer interface {
	FileModified(string)
	FileRenamed(string)
}

type Notifier struct {
	consumer Consumer
	watcher  *fsnotify.Watcher
}

func NewNotifier(dir string, consumer Consumer) (*Notifier, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err = w.Add(dir); err != nil {
		return nil, err
	}

	return &Notifier{
		consumer: consumer,
		watcher:  w,
	}, nil
}

func (n *Notifier) Watch() {
	for {
		select {
		case event, ok := <-n.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				go n.consumer.FileModified(event.Name)
			}
		case err, ok := <-n.watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}

func (n *Notifier) Close() {
	n.watcher.Close()
}
