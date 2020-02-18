package fileutil_test

import (
	"os"
	"testing"

	fileutil "github.com/jjkoh95/image-micro-web-framework/file-utils"
)

const testFile = "test"

func TestMain(m *testing.M) {
	code := m.Run()
	// tear down
	tearDownFile()
	os.Exit(code)
}

func TestSaveToStorage(t *testing.T) {
	byteArray := []byte("to be written to a file\n")
	err := fileutil.SaveToStorage(byteArray, testFile)
	if err != nil {
		t.Error("Expect to save without error")
	}
}

func TestHasFile(t *testing.T) {
	f1 := "fileUtils.go"
	hasFile := fileutil.HasFile(f1)
	if !hasFile {
		t.Error("Expect haveFile to return true here")
	}

	f2 := "does_not_exist"
	hasFile = fileutil.HasFile(f2)
	if hasFile {
		t.Error("Expect haveFile to return false here")
	}
}

func tearDownFile() {
	os.Remove(testFile)
}
