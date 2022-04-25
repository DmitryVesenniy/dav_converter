package analysis

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Analysis(filePath string) error {
	errorFolders := make([]string, 0)
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
				folderName := strings.Split(name, ".")[0]
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

				if countFramesFromPath != countFramesFromMP4 {
					fmt.Printf("[!] Внимание кадры сконвертированны с ошибкой! Папка: %s\n", folder)
					errorFolders = append(errorFolders, folder)
				}
			}
		}
	}

	if len(errorFolders) > 0 {
		ReportAnalysis(errorFolders)
	}

	return nil
}
