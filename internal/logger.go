package internal

import (
	"fmt"
	"log"
	"os"
	"time"
)

const chanBufferLen = 50

type Log interface {
	Created(string)
	Modified(string)
	Deleted(string)
	Backup(string)
	Error(string)
	GetFile() *os.File
	ReadActiveLog() chan string
	Close()
}

var _ Log = (*logger)(nil)

type logger struct {
	file *os.File
	log  *log.Logger
	ch   chan string
}

// Create new Logger instance with log file path provided
func NewLogger(path string) (*logger, error) {
	file, err := Utils.OpenFile(path)
	if err != nil {
		return nil, err
	}

	return &logger{
		file: file,
		log:  log.New(file, "", 0),
		ch:   make(chan string, chanBufferLen),
	}, nil
}

// Closes logger instance
func (l *logger) Close() {
	l.file.Close()
	close(l.ch)
}

// Logs 'created' message to log file
func (l *logger) Created(name string) {
	l.putLog(fmt.Sprintf("created  | %s", name))
}

// Logs 'modified' message to log file
func (l *logger) Modified(name string) {
	l.putLog(fmt.Sprintf("modified | %s", name))
}

// Logs 'deleted' message to log file
func (l *logger) Deleted(name string) {
	l.putLog(fmt.Sprintf("deleted  | %s", name))
}

// Logs 'backuped' message to log file
func (l *logger) Backup(name string) {
	l.putLog(fmt.Sprintf("backedup | %s", name))
}

// Logs 'error' message to log file
func (l *logger) Error(name string) {
	l.putLog(fmt.Sprintf("error    | %s", name))
}

// Returns log file instance
func (l *logger) GetFile() *os.File {
	return l.file
}

// Returns channel for active log messages
func (l *logger) ReadActiveLog() chan string {
	return l.ch
}

func (l *logger) putLog(message string) {
	if l.log == nil || l.file == nil {
		return
	}

	t := time.Now().Format(time.DateTime)
	m := fmt.Sprintf("%s | %s", t, message)

	l.log.Printf(m)
	select {
	case l.ch <- m:
	default:
		return
	}
}
