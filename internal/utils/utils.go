package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	appendFileFlags = os.O_RDWR | os.O_CREATE | os.O_APPEND
	truncFileFlags  = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	FilePerms       = 0644
)

// Creates folder if it does not exist
func CreateFolder(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return os.Mkdir(path, os.ModePerm)
	}
	return nil
}

// Checks if path is a file
func IsFile(path string) bool {
	f, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !f.IsDir()
}

// Copies file from 'in' path to 'out' path
func CopyFile(in string, out string) error {
	fin, err := os.Open(in)
	if err != nil {
		return fmt.Errorf("can't open a file: %v", err)
	}
	defer fin.Close()

	err = CreateFolder(filepath.Dir(out))
	if err != nil {
		return fmt.Errorf("error while creating output directory: %v", err)
	}

	fout, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("can't create a new file: %v", err)
	}
	defer fout.Close()

	if _, err = io.Copy(fout, fin); err != nil {
		return fmt.Errorf("can't copy file contents: %v", err)
	}

	return nil
}

// Overwrites existing file or creates new one
func OverwriteFile(path string, data []byte) error {
	f, err := os.OpenFile(path, truncFileFlags, FilePerms)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}

	return nil
}

// Get file contents
func ReadFile(path string) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, FilePerms)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	res, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func OpenFile(path string) (*os.File, error) {
	return os.OpenFile(path, appendFileFlags, FilePerms)
}

func DeleteFile(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return os.Remove(path)
}

// Scans file and calls callback function for each line
func FileScanner(file *os.File, cb func(string)) {
	if file == nil || cb == nil {
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cb(scanner.Text())
	}
}

// Utility for CLI prompt messages
func Prompt(text string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("- %s: \n", text)
		if answer, _ := reader.ReadString('\n'); answer != "" {
			return strings.TrimSpace(answer)
		}
	}
}
