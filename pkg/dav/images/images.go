package images

import (
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"os"
)

type ImageFrame struct {
	filePath string
}

func (i *ImageFrame) SaveDecodingImg(data []byte) error {
	r := bytes.NewReader(data)
	img, err := jpeg.Decode(r)
	if err != nil {
		return i.SaveImg(data)
		// return fmt.Errorf("error decode jpeg from bytes: %w", err)
	}

	var imageBuf bytes.Buffer
	err = jpeg.Encode(&imageBuf, img, nil)
	if err != nil {
		return i.SaveImg(data)
		// return fmt.Errorf("error encode jpeg from bytes: %w", err)
	}

	file, err := os.Create(i.filePath)
	if err != nil {
		return fmt.Errorf("error create file: %w", err)
	}
	defer file.Close()
	fw := bufio.NewWriter(file)

	_, err = fw.Write(imageBuf.Bytes())

	return err
}

func (i *ImageFrame) SaveImg(data []byte) error {
	file, err := os.Create(i.filePath)
	if err != nil {
		return fmt.Errorf("error create file: %w", err)
	}
	defer file.Close()
	_, err = file.Write(data)

	return err
}

func NewImageFrame(filePath string) *ImageFrame {
	return &ImageFrame{
		filePath: filePath,
	}
}
