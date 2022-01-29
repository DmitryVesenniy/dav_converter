package main

import (
	"fmt"
	"runtime"
	"strings"

	"dav_converter/configs"
	"dav_converter/internal/prestart"
	"dav_converter/pkg/dav"
	"dav_converter/pkg/repository/files"
)

func main() {
	config, err := configs.Get("settings.txt")
	if err != nil {
		fmt.Println("Не был найден файл настроек settings.txt")
		fmt.Println("Для выхода нажмите Enter.")
		fmt.Scanln()
		return
	}

	countCPU := runtime.NumCPU()
	pathList := strings.Split(config.PathList, ",")

	for i, path := range pathList {
		pathList[i] = strings.TrimSpace(path)
	}

	if len(pathList) == 0 {
		fmt.Println("Не была указана ни одна папка с файлами .dav")
		fmt.Println("Для выхода нажмите Enter.")
		fmt.Scanln()
		return
	}

	// анализируем указанные папки, есть ли в них нужные нам файлы
	for _, _path := range pathList {
		files, err := prestart.SearchFilesFromDir(_path, ".dav")
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		if len(files) == 0 {
			fmt.Printf("В папке %s не найдены файлы с расширением .dav\n", _path)
		} else {
			fmt.Printf("В папке %s найдены файлы:\n", _path)
			for _, _fName := range files {
				fmt.Println("    ...", _fName)
			}
		}
	}

	fmt.Println("Для продолжения нажмите Enter.")
	fmt.Scanln()

	davPath := files.NewDavPathFiles(config.PathOut, pathList)
	defer davPath.Close()

	err = dav.Converter(config, davPath, countCPU)

	if err != nil {
		fmt.Printf("Во время работы программы возникла ошибка: %v\n", err)
		fmt.Println("Для выхода нажмите Enter.")
		fmt.Scanln()
		return
	}

	fmt.Println("Все файлы перекодированы!")
	fmt.Println("Для выхода нажмите Enter.")
	fmt.Scanln()
}
