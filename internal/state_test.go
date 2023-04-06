package internal

import (
	"reflect"
	"testing"
)

const (
	testDate        = "2000-01-02"
	testName        = "regex"
	testHotPath     = "hot"
	testBackupPath  = "backup"
	testStateString = `{"filter":{"date":"` +
		testDate +
		`","name":"` +
		testName +
		`"},"dirs":{"hot":"` +
		testHotPath +
		`","backup":"` +
		testBackupPath +
		`"}}`
)

func setupStateTest() {
	Utils.CreateFolder(testDir)
	Utils.OpenFile(testFile)
}

func Test_LoadState_ProvidePath_SetsPath(t *testing.T) {
	state := LoadState(testFile)

	if state.path != testFile {
		t.Logf("wrong state path: %s", state.path)
		t.FailNow()
	}

	if state.GetDirectories() != nil {
		t.Logf("directories state should be nil")
		t.FailNow()
	}

	if state.GetFilter() != nil {
		t.Logf("filter state should be nil")
		t.FailNow()
	}
}

func Test_LoadState_ProvideStateFile_ReturnsState(t *testing.T) {
	setupStateTest()
	defer cleanupTest()

	Utils.OverwriteFile(testFile, []byte(testStateString))

	state := LoadState(testFile)

	expectedDirState := Directories{
		Hot:    testHotPath,
		Backup: testBackupPath,
	}

	expectedFilterState := Filter{
		Name: testName,
		Date: testDate,
	}

	if dirState := state.GetDirectories(); dirState != nil {
		if !reflect.DeepEqual(*dirState, expectedDirState) {
			t.Logf("bad directories state: %v - %v", *dirState, expectedDirState)
			t.FailNow()
		}
	} else {
		t.Logf("directories state = nil")
		t.FailNow()
	}

	if filterState := state.GetFilter(); filterState != nil {
		if !reflect.DeepEqual(*filterState, expectedFilterState) {
			t.Logf("bad filter state: %v - %v", *filterState, expectedFilterState)
			t.FailNow()
		}
	} else {
		t.Logf("filter state = nil")
		t.FailNow()
	}
}
