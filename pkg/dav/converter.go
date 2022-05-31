package dav

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"dav_converter/configs"
	"dav_converter/pkg/dav/converter"
	"dav_converter/pkg/dav/exception"
	"dav_converter/pkg/dav/images"
	"dav_converter/pkg/repository"
	"dav_converter/pkg/subprocess"
	"dav_converter/pkg/utils"
	"dav_converter/pkg/worker"
)

const (
	maxErrorsConverter = 5
)

func PreviewConverter(
	cfg *configs.ConfigApp,
	dav repository.IDavPath,
	countProcess int,
) error {
	davFiles := make([]repository.IDavFile, 0, 0)
	// получаем все файлы .dav
	for {
		_files, err := dav.Next()
		if err != nil {
			if errors.Is(err, exception.ErrorStopIterator) {
				break
			}
		}
		davFiles = append(davFiles, _files...)
	}

	for _, _dav := range davFiles {
		// fmt.Println("**** _dav: ", _dav.GetBasePath())
		// fmt.Println("**** _dav: ", _dav.GetName())
		// fmt.Println("**** _dav: ", _dav.GetPathFrame())
		err := runConverterPreview(_dav)
		if err != nil {
			return err
		}
	}

	return nil
}

func Converter(
	cfg *configs.ConfigApp,
	dav repository.IDavPath,
	countProcess int,
) []error {
	davFiles := make([]repository.IDavFile, 0, 0)
	errorsConverter := make([]error, 0)
	// получаем все файлы .dav
	for {
		_files, err := dav.Next()
		if err != nil {
			if errors.Is(err, exception.ErrorStopIterator) {
				break
			}
			errorsConverter = append(errorsConverter, err)
			return errorsConverter
		}
		davFiles = append(davFiles, _files...)
	}

	defs := make([]func() error, 0, len(davFiles))
	for _, _davFile := range davFiles {
		if stat, err := os.Stat(_davFile.GetPathFrame()); err == nil && stat.IsDir() {
			files, err := ioutil.ReadDir(_davFile.GetPathFrame())
			if err != nil {
				errorsConverter = append(errorsConverter, fmt.Errorf("error read path frame: %w", err))
				return errorsConverter
			}

			if len(files) > 0 {
				// ранее уже были тут и сохраняли все фреймы
				if cfg.SkipExist {
					continue
				} else {
					for _, _f := range files {
						err := os.Remove(filepath.Join(_davFile.GetPathFrame(), _f.Name()))
						if err != nil {
							errorsConverter = append(errorsConverter, fmt.Errorf("error delete file from path frame: %w", err))
							return errorsConverter
						}
					}

					baseName := strings.Split(_davFile.GetName(), ".")[0]
					mp4 := filepath.Join(_davFile.GetBasePath(), fmt.Sprintf("%s.mp4", baseName))
					isExist, err := utils.Exists(mp4)
					if err != nil {
						errorsConverter = append(errorsConverter, fmt.Errorf("error from utils.Exists: %w", err))
						return errorsConverter
					}
					if isExist {
						err := os.Remove(mp4)
						if err != nil {
							errorsConverter = append(errorsConverter, fmt.Errorf("failed to delete file: %w", err))
							return errorsConverter
						}
					}
				}
			}
		} else {
			isExistRoot, err := utils.Exists(_davFile.GetBasePath())
			if err != nil {
				errorsConverter = append(errorsConverter, fmt.Errorf("error get stat path: %w", err))
				return errorsConverter
			}

			if !isExistRoot {
				err = os.Mkdir(_davFile.GetBasePath(), 0755)
				if err != nil {
					errorsConverter = append(errorsConverter, fmt.Errorf("не удалось создать папку: %s; error: %w", _davFile.GetBasePath(), err))
					return errorsConverter
				}
			}

			err = os.Mkdir(_davFile.GetPathFrame(), 0755)
			if err != nil {
				errorsConverter = append(errorsConverter,
					fmt.Errorf("не удалось создать папку: %s; error: %w", _davFile.GetPathFrame(), err))
				return errorsConverter
			}
		}

		df := _davFile
		defs = append(defs, func() error {
			return runConverter(df, cfg)
		})
	}

	stop := make(chan struct{}, 1)
	done := make(chan struct{})
	errorCh := make(chan error, 2)
	go worker.Worker(defs, countProcess, stop, done, errorCh)

	countErrors := 0

	// LOOP:
	for {
		select {
		case err := <-errorCh:
			errorsConverter = append(errorsConverter, err)
			countErrors++

		case <-done:
			return errorsConverter
		}

		if countErrors >= maxErrorsConverter {
			stop <- struct{}{}
			// break LOOP
		}
		runtime.Gosched()
	}

	errorsConverter = append(errorsConverter,
		fmt.Errorf("max count error"))
	return errorsConverter
}

func runConverterPreview(_davFile repository.IDavFile) error {
	converter, err := converter.NewConverter(_davFile.GetReader())
	if err != nil {
		return fmt.Errorf("[!] Ошибка конвертирования <%s>: %w", _davFile.GetPathFrame(), err)
	}

	tag, _, err := converter.Next()
	imgBytes, err := converter.GetImagesOnTagIDX(tag)
	if err != nil {
		return fmt.Errorf("[!] Ошибка конвертирования <%s>: error get image from frame: %w", _davFile.GetPathFrame(), err)
	}

	outJpg := fmt.Sprintf("%s.jpg", _davFile.GetPathFrame())

	return utils.WriteFile(imgBytes, outJpg)
}

func runConverter(_davFile repository.IDavFile, cfg *configs.ConfigApp) error {
	converter, err := converter.NewConverter(_davFile.GetReader())
	if err != nil {
		return fmt.Errorf("[!] Ошибка конвертирования <%s>: %w", _davFile.GetPathFrame(), err)
	}

	var preview []byte

	fmt.Printf("[!] %s: Сохранение кадров...\n", _davFile.GetPathFrame())
	// сохраняем все кадры в папку
	i := 1
	var currentImgBytes []byte
	for {
		tag, _, err := converter.Next()
		if err != nil {
			break
		}

		imgBytes, err := converter.GetImagesOnTagIDX(tag)
		if err != nil {
			return fmt.Errorf("[!] Ошибка конвертирования <%s>: error get image from frame: %w", _davFile.GetPathFrame(), err)
		}

		if len(preview) == 0 {
			preview = imgBytes
		}

		if len(imgBytes) == 0 {
			imgBytes = currentImgBytes
		} else {
			currentImgBytes = imgBytes
		}

		if len(imgBytes) == 0 {
			continue
		}

		img := images.NewImageFrame(filepath.Join(_davFile.GetPathFrame(), fmt.Sprintf("%v.jpg", i)))
		err = img.SaveDecodingImg(imgBytes)
		if err != nil {
			return fmt.Errorf("[!] Ошибка конвертирования <%s>: error save image: %w", _davFile.GetPathFrame(), err)
		}
		i++
	}
	fmt.Printf("[!] %s: Всего сохранено кадров: %d\n", _davFile.GetPathFrame(), i)
	fmt.Printf("[!] %s: Конвертация кадров в mp4...\n", _davFile.GetPathFrame())
	pathFromConvert := filepath.Join(_davFile.GetPathFrame(), "%d.jpg") //fmt.Sprintf("%s/%%d.jpg", _davFile.GetPathFrame())
	outMp4 := fmt.Sprintf("%s.mp4", _davFile.GetPathFrame())
	outJpg := fmt.Sprintf("%s.jpg", _davFile.GetPathFrame())

	if cfg.IsDev {
		fmt.Println("[!DEV] pathFromConvert: ", pathFromConvert)
		fmt.Println("[!DEV] outMp4: ", outMp4)
	}

	// конвертируем кадры в mp4
	commands := map[string]subprocess.Command{
		"linux": {
			NameCommand: "ffmpeg",
			Args: []string{"-i", pathFromConvert, "-c:v", "libx264", "-crf", "25",
				"-vf", "fps=25", "-pix_fmt", "yuv420p", outMp4},
		},
		"windows": {
			NameCommand: "ffmpeg.exe",
			Args: []string{"-i", pathFromConvert, "-c:v", "libx264", "-crf", "25",
				"-vf", "fps=25", "-pix_fmt", "yuv420p", outMp4},
		},
	}

	process, err := subprocess.NewSysProcess(commands)
	if err != nil {
		return fmt.Errorf("[!] Ошибка конвертирования <%s>: error init process: %w", _davFile.GetPathFrame(), err)
	}

	err = process.Run()
	if err != nil {
		return fmt.Errorf("[!] Ошибка конвертирования <%s>: error run process: %w", _davFile.GetPathFrame(), err)
	}

	err = utils.WriteFile(preview, outJpg)

	return err
}
