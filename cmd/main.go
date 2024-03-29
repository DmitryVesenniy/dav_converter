package main

import (
	"flag"
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"dav_converter/configs"
	"dav_converter/internal/prestart"
	"dav_converter/pkg/analysis"
	"dav_converter/pkg/dav"
	"dav_converter/pkg/dav/converter"
	"dav_converter/pkg/repository/files"
)

var (
	pathPattern = regexp.MustCompile(`\"(.*?)\"`)
)

func main() {
	file := flag.String("file", "", "путь до файла dav")
	analysisPath := flag.String("analysis", "", "путь до папки с отконвертированными папками заездов")
	isDev := flag.Bool("dev", false, "включить режим разработчика")
	flag.Parse()

	if *file != "" {
		davFile, err := files.NewDavFile(*file)
		if err != nil {
			fmt.Printf("Непредвиденная ошибка: %s\n", err)
			return
		}

		converter, err := converter.NewConverter(davFile.GetReader())
		if err != nil {
			fmt.Printf("Непредвиденная ошибка: %s\n", err)
			return
		}

		fmt.Printf("Количество кадров: %d\n", len(converter.IndexTable))

		davFile.Close()
		return
	}

	if *analysisPath != "" {
		fmt.Printf("[!] Анализ конвертированных папок по адресу: %s\n", *analysisPath)
		err := analysis.Analysis(*analysisPath)
		if err != nil {
			fmt.Printf("во время анализа возникла ошибка: %s\n", err.Error())
		}
		return
	}

	config, err := configs.Get("settings.txt")
	if err != nil {
		fmt.Println("Не был найден файл настроек settings.txt")
		fmt.Println("Для выхода нажмите Enter.")
		fmt.Scanln()
		return
	}

	config.IsDev = *isDev

	if config.IsDev {
		fmt.Println("[!DEV] Включен режим разработчика")
	}

	countCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(countCPU)
	pathList := make([]string, 0)
	previewPathList := make([]string, 0)

	for _, path := range pathPattern.FindAllString(config.PathList, -1) {
		_path := strings.TrimSpace(path)
		if _path != "" {
			pathList = append(pathList, _path)
		}
	}

	fmt.Println("pathList: ", pathList)

	for _, path := range pathPattern.FindAllString(config.PreviewPath, -1) {
		_path := strings.TrimSpace(path)
		if _path != "" {
			previewPathList = append(previewPathList, _path)
		}
	}

	if len(pathList) == 0 && len(previewPathList) == 0 {
		fmt.Println("Не была указана ни одна папка с файлами .dav")
		fmt.Println("Для выхода нажмите Enter.")
		fmt.Scanln()
		return
	}

	var errList []error

	if len(pathList) > 0 {
		fmt.Println("[!] Конвертирую")
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

		errList = dav.Converter(config, davPath, countCPU)
	}

	if len(previewPathList) > 0 {
		fmt.Println("[!] Создаю превью")
		davPath := files.NewDavPathFiles(config.PathOut, previewPathList)
		defer davPath.Close()

		err := dav.PreviewConverter(config, davPath, countCPU)
		if err != nil {
			errList = append(errList, err)
		}
	}

	if len(errList) > 0 {
		fmt.Println("[!] Во время работы программы возникли следующие ошибки:")

		for _, _err := range errList {
			fmt.Printf("    -- %v\n", _err)
		}
		fmt.Println("Для выхода нажмите Enter.")
		fmt.Scanln()
		return
	}

	fmt.Println("Все файлы перекодированы!")
	fmt.Println("Для выхода нажмите Enter.")
	fmt.Scanln()
}
