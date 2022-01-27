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
)

type DavPathFiles struct {
	Dirs         []string
	currentIndex int
}

type DavFile struct {
	Name      string
	PathFrame string
	BasePath  string
	File      io.ReadSeeker
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
				pathFrames := filepath.Join(path, fileName)
				_fileObj, err := os.Open(filepath.Join(path, file.Name()))
				if err != nil {
					return nil, fmt.Errorf("error open dav file %w", err)
				}
				_fileDav := &DavFile{
					Name:      file.Name(),
					PathFrame: pathFrames,
					BasePath:  path,
					File:      io.ReadSeeker(_fileObj),
				}
				fileDavList = append(fileDavList, _fileDav)
			}
		}
	}

	return fileDavList, nil
}

func NewDavPathFiles(pathList []string) *DavPathFiles {
	return &DavPathFiles{
		Dirs:         pathList,
		currentIndex: 0,
	}
}
