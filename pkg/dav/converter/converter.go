package converter

import (
	"fmt"
	"io"

	"dav_converter/pkg/dav/exception"
	"dav_converter/pkg/dav/usecase"
)

type Converter struct {
	ReaderDav    io.ReadSeeker
	currentIndex int
	HeaderDav    usecase.HeaderDav
	IndexTable   []usecase.TagFrameIDX
}

func (c *Converter) convert() error {
	var err error
	c.HeaderDav, err = getHeaderDav(c.ReaderDav)

	if err != nil {
		return fmt.Errorf("error get header dav: %w", err)
	}

	// Заполняем таблицу индексов
	c.ReaderDav.Seek(c.HeaderDav.OffsetIndex, io.SeekStart)
	for {
		tagIndex, err := getTagIdx(c.ReaderDav)

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error get tagIndex dav: %w", err)
		}

		c.IndexTable = append(c.IndexTable, tagIndex)

		// _, err = c.ReaderDav.Seek(20, io.SeekCurrent)
		// if err != nil {
		// 	return fmt.Errorf("error seek file: %w", err)
		// }
	}
	return nil
}

func (c *Converter) GetImagesOnTagIDX(t usecase.TagFrameIDX) ([]byte, error) {
	c.ReaderDav.Seek(t.KadrOffset, io.SeekStart)

	headerFrame, err := getHeaderFrame(c.ReaderDav)
	if err != nil {
		return nil, fmt.Errorf("error get header frame: %w", err)
	}

	kadrBytes := make([]byte, headerFrame.SizeKadr)

	_, err = c.ReaderDav.Read(kadrBytes)
	if err != nil {
		return nil, fmt.Errorf("error reader kadr: %w", err)
	}

	return kadrBytes, nil
}

func (c *Converter) GetFrameIndex(index int) (usecase.TagFrameIDX, error) {
	if index >= len(c.IndexTable) {
		return usecase.TagFrameIDX{}, exception.FrameBounds{}
	}

	return c.IndexTable[index], nil
}

func (c *Converter) Next() (usecase.TagFrameIDX, int, error) {
	if c.currentIndex < len(c.IndexTable) {
		tagIndex := c.IndexTable[c.currentIndex]
		currentIndex := c.currentIndex
		c.currentIndex++
		return tagIndex, currentIndex, nil
	}
	return usecase.TagFrameIDX{}, 0, exception.FrameBounds{}
}

func NewConverter(f io.ReadSeeker) (*Converter, error) {
	converter := &Converter{
		ReaderDav:    f,
		currentIndex: 0,
	}

	err := converter.convert()

	if err != nil {
		return nil, err
	}

	return converter, nil
}
