package dav

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"dav_converter/configs"
	"dav_converter/pkg/dav/converter"
	"dav_converter/pkg/dav/images"
	"dav_converter/pkg/repository"
	"dav_converter/pkg/subprocess"
	"dav_converter/pkg/utils"
	"dav_converter/pkg/worker"
)

const (
	maxErrorsConverter = 5
)

func Converter(
	cfg *configs.ConfigApp,
	dav repository.IDavPath,
	countProcess int,
) error {
	davFiles := make([]repository.IDavFile, 0, 0)

	// получаем все файлы .dav
	for {
		_files, err := dav.Next()
		if err != nil {
			break
		}
		davFiles = append(davFiles, _files...)
	}

	defs := make([]func() error, 0, len(davFiles))
	for _, _davFile := range davFiles {
		if stat, err := os.Stat(_davFile.GetPathFrame()); err == nil && stat.IsDir() {
			files, err := ioutil.ReadDir(_davFile.GetPathFrame())
			if err != nil {
				return fmt.Errorf("error read path frame: %w", err)
			}

			if len(files) > 0 {
				// ранее уже были тут и сохраняли все фреймы
				if cfg.SkipExist {
					continue
				} else {
					for _, _f := range files {
						err := os.Remove(filepath.Join(_davFile.GetPathFrame(), _f.Name()))
						if err != nil {
							return fmt.Errorf("error delete file from path frame: %w", err)
						}
					}

					baseName := strings.Split(_davFile.GetName(), ".")[0]
					mp4 := filepath.Join(_davFile.GetBasePath(), fmt.Sprintf("%s.mp4", baseName))
					isExist, err := utils.Exists(mp4)
					if err != nil {
						return fmt.Errorf("error from utils.Exists: %w", err)
					}
					if isExist {
						err := os.Remove(mp4)
						if err != nil {
							return fmt.Errorf("failed to delete file: %w", err)
						}
					}
				}
			}
		} else {
			err = os.Mkdir(_davFile.GetPathFrame(), 0755)
			if err != nil {
				return fmt.Errorf("error create path frame: %w", err)
			}
		}

		df := _davFile
		defs = append(defs, func() error {
			return runConverter(df)
		})
	}

	stop := make(chan struct{})
	done := make(chan struct{})
	errorCh := make(chan error, 2)
	go worker.Worker(defs, countProcess, stop, done, errorCh)

	countErrors := 0
LOOP:
	for {
		select {
		case err := <-errorCh:
			fmt.Println(err)
			countErrors++

		case <-done:
			return nil
		}

		if countErrors >= maxErrorsConverter {
			stop <- struct{}{}
			break LOOP
		}
		runtime.Gosched()
	}

	return fmt.Errorf("max count error")
}

func runConverter(_davFile repository.IDavFile) error {
	converter, err := converter.NewConverter(_davFile.GetReader())
	if err != nil {
		return fmt.Errorf("error create new converter: %w", err)
	}

	fmt.Printf("[!] %s: Сохранение кадров...\n", _davFile.GetPathFrame())
	count := 0
	// сохраняем все кадры в папку
	for {
		tag, i, err := converter.Next()
		if err != nil {
			break
		}

		imgBytes, err := converter.GetImagesOnTagIDX(tag)
		if err != nil {
			return fmt.Errorf("error get image from frame: %w", err)
		}

		img := images.NewImageFrame(filepath.Join(_davFile.GetPathFrame(), fmt.Sprintf("%v.jpg", i+1)))

		err = img.SaveImg(imgBytes)
		if err != nil {
			return fmt.Errorf("error save image: %w", err)
		}
		count++
	}
	fmt.Printf("[!] %s: Всего сохранено кадров: %d\n", _davFile.GetPathFrame(), count)
	fmt.Printf("[!] %s: Конвертация кадров в mp4...\n", _davFile.GetPathFrame())
	pathFromConvert := fmt.Sprintf("%s/%%d.jpg", _davFile.GetPathFrame())
	outMp4 := fmt.Sprintf("%s.mp4", _davFile.GetPathFrame())

	// конвертируем кадры в mp4
	commands := map[string]subprocess.Command{
		"linux": {
			NameCommand: "ffmpeg",
			Args: []string{"-i", pathFromConvert, "-c:v", "libx264",
				"-vf", "fps=25", "-pix_fmt", "yuv420p", outMp4},
		},
		"windows": {
			NameCommand: "ffmpeg.exe",
			Args: []string{"-i", pathFromConvert, "-c:v", "libx264",
				"-vf", "fps=25", "-pix_fmt", "yuv420p", outMp4},
		},
	}

	process, err := subprocess.NewSysProcess(commands)
	if err != nil {
		return fmt.Errorf("error init process: %w", err)
	}

	err = process.Run()
	if err != nil {
		return fmt.Errorf("error run process: %w", err)
	}

	return nil
}
