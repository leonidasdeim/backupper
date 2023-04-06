package internal

import (
	"strings"
	"testing"
)

func setupLoggerTest() {
	Utils.CreateFolder(testDir)
}

func Test_NewLogger_BadPath_ReturnsError(t *testing.T) {
	l, err := NewLogger("")

	if l != nil {
		t.Log("logger instance should be nil")
		t.FailNow()
	}

	if err == nil || !strings.Contains(err.Error(), "can't open log file") {
		t.Logf("bad error received: %v", err)
		t.FailNow()
	}
}

func Test_NewLogger_CorrectPath_ShouldCreateFile(t *testing.T) {
	setupLoggerTest()
	defer cleanupTest()

	l, err := NewLogger(testFile)
	if l == nil {
		t.Log("logger instance should not be nil")
		t.FailNow()
	}
	defer l.Close()

	if err != nil {
		t.Logf("error received: %v", err)
		t.FailNow()
	}

	if !Utils.IsFile(testFile) {
		t.Log("should create log file")
		t.FailNow()
	}
}

func Test_NewLogger_SendLog_ShouldLogToFile(t *testing.T) {
	setupLoggerTest()
	defer cleanupTest()

	l, _ := NewLogger(testFile)
	if l == nil {
		t.Log("logger instance should not be nil")
		t.FailNow()
	}
	defer l.Close()

	l.Created("filename")
	l.Modified("filename")
	l.Backup("filename")
	l.Deleted("filename")
	l.Error("filename")

	numberOfLines := 0
	file, _ := Utils.OpenFile(testFile)
	Utils.FileScanner(file, func(s string) {
		numberOfLines++
	})

	if numberOfLines != 5 {
		t.Logf("log file contain %d lines, should be 5", numberOfLines)
		t.FailNow()
	}
}

func Test_NewLogger_SendLog_ShouldPutToChannel(t *testing.T) {
	setupLoggerTest()
	defer cleanupTest()

	l, _ := NewLogger(testFile)
	if l == nil {
		t.Log("logger instance should not be nil")
		t.FailNow()
	}
	defer l.Close()

	l.Created("filename")
	l.Modified("filename")
	l.Backup("filename")
	l.Deleted("filename")
	l.Error("filename")

	numberOfMessages := len(l.RealTimeLog())
	if numberOfMessages != 5 {
		t.Logf("channel contain %d messages, should be 5", numberOfMessages)
		t.FailNow()
	}
}
