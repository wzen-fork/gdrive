package auth

import (
	"fmt"
	"os"
	"path/filepath"
)

type Closeable interface {
	Close() error
}

func mkdir(path string) error {
	dir := filepath.Dir(path)
	if fileExists(dir) {
		return nil
	}
	return os.Mkdir(dir, 0700)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

func silentClose(resource Closeable){
	err := resource.Close()
	if err != nil {
		fmt.Printf("close resource error: %v\n", err)
	}
}

func silentError(err error){
	if err != nil {
		fmt.Printf("silent error: %v\n", err)
	}
}
