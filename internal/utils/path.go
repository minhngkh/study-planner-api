package utils

import (
	"path/filepath"
	"runtime"
)

func CurrentFileDir() string {
	_, fileName, _, ok := runtime.Caller(1)
	if !ok {
		panic("No caller information")
	}

	curDir := filepath.Dir(fileName)

	return curDir
}
