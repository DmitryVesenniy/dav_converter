package analysis

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	maximumDeviation = 0
)

func Analysis(filePath string) error {
	errorFolders := make(map[string]struct{}, 0)
	folders, err := ioutil.ReadDir(filePath)
	if err != nil {
		return fmt.Errorf("не удалось получить файлы из папки %w", err)
	}

	for _, path := range folders {
		if path.IsDir() {
			folder := filepath.Join(filePath, path.Name())

			flePaths, err := ioutil.ReadDir(folder)
			if err != nil {
				return fmt.Errorf("не удалось получить файлы из папки %w", err)
			}

			filesMp4 := make([]*AnalysisMP4, 0)

			for _, _file := range flePaths {
				_ext := filepath.Ext(_file.Name())
				if !_file.IsDir() && _ext == ".mp4" {
					fileObj, err := os.Open(filepath.Join(folder, _file.Name()))
					if err != nil {
						return fmt.Errorf("не удалось открыть файл mp4 %w", err)
					}
					filesMp4 = append(filesMp4, NewAnalysisMP4(fileObj))
				}
			}

			for _, fileAnalysis := range filesMp4 {
				name := fileAnalysis.file.Name()
				folderName := strings.TrimSuffix(name, ".mp4")
				// filePath := filepath.Join(folder, folderName)

				frames, err := ioutil.ReadDir(folderName)
				if err != nil {
					return fmt.Errorf("не удалось получить файлы из папки %w", err)
				}

				countFramesFromPath := uint(len(frames))
				countFramesFromMP4, err := fileAnalysis.Analysis()
				if err != nil {
					return fmt.Errorf("не удалось получить количество кадров из mp4 %w", err)
				}

				if (countFramesFromPath - countFramesFromMP4) > maximumDeviation {
					// fmt.Printf("[!] Внимание кадры сконвертированны с ошибкой! Папка: %s\n", folder)
					errorFolders[folder] = struct{}{}
				}
			}
		}
	}

	if len(errorFolders) > 0 {
		errorsFoldersList := make([]string, 0, len(errorFolders))
		for _folder := range errorFolders {
			fmt.Printf("[!] Внимание кадры сконвертированны с ошибкой! Папка: %s\n", _folder)
			errorsFoldersList = append(errorsFoldersList, _folder)
		}
		ReportAnalysis(errorsFoldersList)
	}

	return nil
}
