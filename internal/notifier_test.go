package internal

import (
	"os"
	"testing"
	"time"
)

func setupNotifierTest() {
	Utils.CreateFolder(testDir)
}

func Test_NewNotifier_BadConsumer_ReturnsError(t *testing.T) {
	n, err := NewNotifier("", nil)

	if n != nil {
		t.Log("notifier instance should be nil")
		t.FailNow()
	}

	if err == nil || err.Error() != "consumer object is not provided" {
		t.Logf("bad error received: %v", err)
		t.FailNow()
	}
}

func Test_NewNotifier_CorrectArguments_ReturnsInstance(t *testing.T) {
	c := &ConsumerStub{}
	n, err := NewNotifier("", c)
	if n == nil {
		t.Log("notifier instance should not be nil")
		t.FailNow()
	}
	defer n.Close()

	if err != nil {
		t.Logf("should not return an error, received: %v", err)
		t.FailNow()
	}
}

func Test_NewNotifier_FilesCreated_NotifiesConsumer(t *testing.T) {
	setupNotifierTest()
	defer cleanupTest()

	c := &ConsumerStub{}
	n, _ := NewNotifier(testDir, c)
	if n == nil {
		t.Log("notifier instance should not be nil")
		t.FailNow()
	}
	defer n.Close()

	go n.Watch()

	Utils.OpenFile(testFile + "0")
	Utils.OpenFile(testFile + "1")
	Utils.OpenFile(testFile + "2")
	Utils.OpenFile(deleteTestFile + "1")
	Utils.OpenFile(deleteTestFile + "2")

	// wait for events to be created in notifiers goroutine
	time.Sleep(100 * time.Millisecond)

	if c.filesCreatedCount != 5 {
		t.Logf("counter = %d, should be 5", c.filesCreatedCount)
		t.FailNow()
	}
}

func Test_NewNotifier_FilesModified_NotifiesConsumer(t *testing.T) {
	setupNotifierTest()
	defer cleanupTest()

	c := &ConsumerStub{}
	n, _ := NewNotifier(testDir, c)
	if n == nil {
		t.Log("notifier instance should not be nil")
		t.FailNow()
	}
	defer n.Close()

	go n.Watch()

	Utils.OpenFile(testFile)
	Utils.OpenFile(deleteTestFile)

	os.WriteFile(testFile, []byte{1}, filePerms)
	os.WriteFile(testFile, []byte{2}, filePerms)
	os.WriteFile(testFile, []byte{3}, filePerms)
	os.WriteFile(deleteTestFile, []byte{1}, filePerms)
	os.WriteFile(deleteTestFile, []byte{2}, filePerms)

	// wait for events to be created in notifiers goroutine
	time.Sleep(100 * time.Millisecond)

	if c.filesModifierCount != 5 {
		t.Logf("counter = %d, should be 5", c.filesModifierCount)
		t.FailNow()
	}
}
