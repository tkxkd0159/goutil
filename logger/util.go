package logger

import (
	"log"
	"os"
)

func MustMkDir(path string, dirperm os.FileMode) {
	err := os.MkdirAll(path, dirperm)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}
}
