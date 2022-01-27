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

func (i *ImageFrame) SaveImg(data []byte) error {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error decode jpeg from bytes: %w", err)
	}

	var imageBuf bytes.Buffer
	err = jpeg.Encode(&imageBuf, img, nil)
	if err != nil {
		return fmt.Errorf("error encode jpeg from bytes: %w", err)
	}

	file, err := os.Create(i.filePath)
	if err != nil {
		return fmt.Errorf("error create file: %w", err)
	}
	fw := bufio.NewWriter(file)

	_, err = fw.Write(imageBuf.Bytes())

	return err
}

func NewImageFrame(filePath string) *ImageFrame {
	return &ImageFrame{
		filePath: filePath,
	}
}