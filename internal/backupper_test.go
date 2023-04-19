package internal

import (
	"path"
	"strings"
	"testing"

	"github.com/leonidasdeim/backupper/internal/utils"
)

func setupBackupperTest() {
	utils.CreateFolder(testDir)
	utils.CreateFolder(testMockDir)
	utils.OpenFile(testFile)
	utils.OpenFile(deleteTestFile)
}

func Test_NewBackupper_LoggerNotProvided_ReturnsError(t *testing.T) {
	b, err := NewBackupper("", nil)

	if b != nil {
		t.Log("backupper instance should be nil")
		t.FailNow()
	}

	if err == nil || err.Error() != "logger object is not provided" {
		t.Logf("bad error received: %v", err)
		t.FailNow()
	}
}

func Test_NewBackupper_BadPath_ReturnsError(t *testing.T) {
	l := &LogStub{}
	b, err := NewBackupper("", l)

	if b != nil {
		t.Log("backupper instance should be nil")
		t.FailNow()
	}

	if err == nil || !strings.Contains(err.Error(), "can't create backup directory") {
		t.Logf("bad error received: %v", err)
		t.FailNow()
	}
}

func Test_NewBackupper_Successful_ReturnsInstance(t *testing.T) {
	setupBackupperTest()
	defer cleanupTest()

	l := &LogStub{}
	b, err := NewBackupper(testBackupDir, l)

	if b == nil {
		t.Log("backupper instance should not be nil")
		t.FailNow()
	}

	if err != nil {
		t.Logf("should not return an error: %v", err)
		t.FailNow()
	}
}

func Test_FileCreated_Successful_ShouldCopyFile(t *testing.T) {
	setupBackupperTest()
	defer cleanupTest()

	l := &LogStub{}
	b, err := NewBackupper(testBackupDir, l)
	if err != nil {
		t.Logf("should not return an error: %v", err)
		t.FailNow()
	}

	b.FileCreated(testFile)
	if !utils.IsFile(testBackupDir + path.Base(testFile) + extension) {
		t.Logf("backup file should exist")
		t.FailNow()
	}
}

func Test_FileCreated_FolderPassed_ShouldNotCopyFolder(t *testing.T) {
	setupBackupperTest()
	defer cleanupTest()

	l := &LogStub{}
	b, err := NewBackupper(testBackupDir, l)
	if err != nil {
		t.Logf("should not return an error: %v", err)
		t.FailNow()
	}

	b.FileCreated(testMockDir)
	if utils.IsFile(testBackupDir + path.Base(testMockDir) + extension) {
		t.Logf("backup file should not exist")
		t.FailNow()
	}
}

func Test_FileCreated_DeletePrefix_ShouldDeleteFile(t *testing.T) {
	setupBackupperTest()
	defer cleanupTest()

	l := &LogStub{}
	b, err := NewBackupper(testBackupDir, l)
	if err != nil {
		t.Logf("should not return an error: %v", err)
		t.FailNow()
	}

	b.FileCreated(deleteTestFile)
	if utils.IsFile(deleteTestFile) {
		t.Logf("backup file should be deleted")
		t.FailNow()
	}
}

func Test_FileModified_Successful_ShouldBackupFile(t *testing.T) {
	setupBackupperTest()
	defer cleanupTest()

	l := &LogStub{}
	b, err := NewBackupper(testBackupDir, l)
	if err != nil {
		t.Logf("should not return an error: %v", err)
		t.FailNow()
	}

	b.FileModified(testFile)
	if !utils.IsFile(testBackupDir + path.Base(testFile) + extension) {
		t.Logf("backup file should exist")
		t.FailNow()
	}
}

func Test_FileModified_FolderPassed_ShouldNotCopyFolder(t *testing.T) {
	setupBackupperTest()
	defer cleanupTest()

	l := &LogStub{}
	b, err := NewBackupper(testBackupDir, l)
	if err != nil {
		t.Logf("should not return an error: %v", err)
		t.FailNow()
	}

	b.FileModified(testMockDir)
	if utils.IsFile(testBackupDir + path.Base(testMockDir) + extension) {
		t.Logf("backup file should not exist")
		t.FailNow()
	}
}
