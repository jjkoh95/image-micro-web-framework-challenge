package fileutil

import (
	"io/ioutil"
	"os"
)

// SaveToStorage saves byte array locally
func SaveToStorage(byteArr []byte, filename string) error {
	err := ioutil.WriteFile(filename, byteArr, 0644) // 0644 stands for permission
	return err
}

// HasFile checks if path exists
func HasFile(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
