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
	// if data[len(data)-1] == '\x00' && data[len(data)-2] == '\xff' {
	// 	data[len(data)-1] = '\xd9'
	// }

	r := bytes.NewReader(data)
	img, err := jpeg.Decode(r)
	// fmt.Printf("[!] %x [%x]: [%x]\n", data[len(data)-1], data[len(data)-2], data[len(data)-10])
	// fmt.Printf("%+v\n", data[len(data)-10:])
	if err != nil {
		return i.SaveImg(data)
		// return fmt.Errorf("error decode jpeg from bytes: %w", err)
	}

	var imageBuf bytes.Buffer
	err = jpeg.Encode(&imageBuf, img, nil)
	if err != nil {
		return i.SaveImg(data)
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
