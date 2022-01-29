package prestart

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// IsExistsPath  determine if such a folder exists
func IsExistsPath(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		if stat.IsDir() {
			return true, nil
		}
		return false, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// SearchFilesFromDir find in the specified folder all files with the extension ext
func SearchFilesFromDir(path string, ext string) ([]string, error) {
	existPath, err := IsExistsPath(path)
	if err != nil {
		return nil, err
	}

	if !existPath {
		return nil, fmt.Errorf("папка %s не найдена", path)
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить файлы из папки %w", err)
	}

	fileDavList := make([]string, 0, 4)

	for _, file := range files {
		if !file.IsDir() {
			_ext := filepath.Ext(file.Name())
			if _ext == ext {
				fileDavList = append(fileDavList, file.Name())
			}
		}
	}

	return fileDavList, nil
}
