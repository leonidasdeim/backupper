package internal

import "os"

var (
	testDir        = "testtmp/"
	testFile       = testDir + "testfile"
	testBackupDir  = testDir + "backupDir/"
	testMockDir    = testDir + "mockDir/"
	deleteTestFile = testDir + "delete_testfile"
)

func cleanupTest() {
	os.RemoveAll(testDir)
}

type LogStub struct{}

func (l *LogStub) Created(string) {
	// empty implementation
}

func (l *LogStub) Modified(string) {
	// empty implementation
}

func (l *LogStub) Deleted(string) {
	// empty implementation
}

func (l *LogStub) Backup(string) {
	// empty implementation
}

func (l *LogStub) Error(string) {
	// empty implementation
}

func (l *LogStub) GetFile() *os.File {
	return nil
}

func (l *LogStub) RealTimeLog() chan string {
	return nil
}

func (l *LogStub) Close() {
	// empty implementation
}

type ConsumerStub struct {
	filesCreatedCount  int
	filesModifierCount int
}

func (c *ConsumerStub) FileCreated(string) {
	c.filesCreatedCount++
}

func (c *ConsumerStub) FileModified(string) {
	c.filesModifierCount++
}
