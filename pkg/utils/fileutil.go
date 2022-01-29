package utils

import (
	"errors"
	"os"
	"strings"
)

// Exists determine the existence of a file
func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func GetListPath(path string) []string {
	sep := string(os.PathSeparator)
	return strings.Split(strings.Trim(path, sep), sep)
}

func GetFirsRootPath(path string) string {
	pathList := GetListPath(path)
	return pathList[len(pathList)-1]
}
