package main

import (
	"fmt"
	"log"
	"runtime"
	"strings"

	"dav_converter/configs"
	"dav_converter/pkg/dav"
	"dav_converter/pkg/repository/files"
)

func main() {
	config, err := configs.Get("settings.txt")
	if err != nil {
		log.Fatal("Не был найден файл настроек settings.txt\n")
	}

	countCPU := runtime.NumCPU()
	pathList := strings.Split(config.PathList, ",")

	for i, path := range pathList {
		pathList[i] = strings.TrimSpace(path)
	}

	if len(pathList) == 0 {
		log.Fatal("Не была указана ни одна папка с файлами .dav\n")
	}

	davPath := files.NewDavPathFiles(pathList)

	err = dav.Converter(config, davPath, countCPU)

	if err != nil {
		log.Fatalf("Во время работы программы возникла ошибка: %v\n", err)
	}

	fmt.Println("Для выхода нажмите Enter.")
	fmt.Scanln()
}
