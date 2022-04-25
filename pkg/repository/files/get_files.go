package files

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"dav_converter/pkg/dav/exception"
	"dav_converter/pkg/repository"
	"dav_converter/pkg/utils"
)

type DavPathFiles struct {
	Dirs         []string
	Path         string
	currentIndex int
	files        []*os.File
}

type DavFile struct {
	Name      string
	PathFrame string
	BasePath  string
	File      io.ReadSeeker
	fileObj   *os.File
}

func (df *DavFile) GetPathFrame() string {
	return df.PathFrame
}

func (df *DavFile) GetReader() io.ReadSeeker {
	return df.File
}

func (df *DavFile) GetBasePath() string {
	return df.BasePath
}

func (df *DavFile) GetName() string {
	return df.Name
}

func (df *DavFile) Close() error {
	return df.fileObj.Close()
}

// SetDavPath Exports
func (dpf *DavPathFiles) SetDavPath(path string) {
	splitPath := strings.Split(path, ",")
	dpf.Dirs = splitPath
}

func (dpf *DavPathFiles) Next() ([]repository.IDavFile, error) {
	davFiles, err := dpf.GetDavList(dpf.currentIndex)
	if err != nil {
		return nil, err
	}
	dpf.currentIndex++
	return davFiles, nil
}

func (dpf *DavPathFiles) Close() {
	for _, file := range dpf.files {
		file.Close()
	}
}

func (dpf *DavPathFiles) GetDavList(index int) ([]repository.IDavFile, error) {
	if len(dpf.Dirs) <= index {
		return nil, exception.ErrorStopIterator
	}

	path := dpf.Dirs[index]
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error read dir dav %w", err)
	}

	fileDavList := make([]repository.IDavFile, 0, 4)

	for _, file := range files {
		// fmt.Println(file.Name(), file.IsDir())
		if !file.IsDir() {
			ext := filepath.Ext(file.Name())
			if ext == ".dav" {
				fileName := strings.Split(file.Name(), ".")[0]
				rootPath := utils.GetFirsRootPath(path)
				pathFrames := filepath.Join(dpf.Path, rootPath, fileName)
				// pathFrames := filepath.Join(path, fileName)

				_fileObj, err := os.Open(filepath.Join(path, file.Name()))
				if err != nil {
					return nil, fmt.Errorf("error open dav file %w", err)
				}
				dpf.files = append(dpf.files, _fileObj)
				_fileDav := &DavFile{
					Name:      file.Name(),
					PathFrame: pathFrames,
					BasePath:  filepath.Join(dpf.Path, rootPath),
					File:      io.ReadSeeker(_fileObj),
				}
				fileDavList = append(fileDavList, _fileDav)
			}
		}
	}

	return fileDavList, nil
}

func NewDavPathFiles(path string, pathList []string) *DavPathFiles {
	return &DavPathFiles{
		Dirs:         pathList,
		Path:         path,
		currentIndex: 0,
	}
}

func NewDavFile(filePath string) (repository.IDavFile, error) {
	ext := filepath.Ext(filePath)
	if ext != ".dav" {
		return nil, fmt.Errorf("Это не dav файл")
	}

	_fileObj, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error open dav file %w", err)
	}

	fileDav := &DavFile{
		Name:    filePath,
		File:    io.ReadSeeker(_fileObj),
		fileObj: _fileObj,
	}

	return fileDav, nil
}
