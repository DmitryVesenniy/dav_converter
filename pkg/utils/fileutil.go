package utils

import (
	"bufio"
	"errors"
	"fmt"
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

func BuffWriteFile(b []byte, fileName string, chunk int) error {
	file, err := os.Create(fileName)
	writer := bufio.NewWriter(file)
	if err != nil {
		return fmt.Errorf("write preview: %w", err)
	}
	defer file.Close()

	for len(b) > 0 {
		writer.Write(b[:chunk]) // запись строки
		b = b[chunk:]
	}

	writer.Flush()
	return nil
}

func WriteFile(b []byte, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("write preview: %w", err)
	}
	defer file.Close()
	file.Write(b)
	return nil
}
