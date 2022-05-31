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

	// fmt.Printf(">> %+v\n", c.HeaderDav)

	// Заполняем таблицу индексов
	c.ReaderDav.Seek(c.HeaderDav.OffsetIndex, io.SeekStart)
	for {
		tagIndex, err := getTagIdx(c.ReaderDav)

		if tagIndex.KadrOffset == 0 {
			return fmt.Errorf("error get tagIndex dav: index table read error")
		}

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

	// fmt.Printf("[!] headerFrame: %+v\n", headerFrame)

	kadrBytes := make([]byte, headerFrame.SizeKadr)

	_, err = c.ReaderDav.Read(kadrBytes)
	if err != nil {
		return nil, fmt.Errorf("error reader kadr: %w", err)
	}
	// fmt.Printf("[!] kadrBytes: %+v\n", kadrBytes)

	return kadrBytes, nil
}

func (c *Converter) GetFrameIndex(index int) (usecase.TagFrameIDX, error) {
	if index >= len(c.IndexTable) {
		return usecase.TagFrameIDX{}, exception.FrameBounds{}
	}

	return c.IndexTable[index], nil
}

func (c *Converter) Next() (usecase.TagFrameIDX, int, error) {
	currentIndex := c.currentIndex
	// nextIndex := usecase.TagFrameIDX{}

	if currentIndex < len(c.IndexTable) {
		tagIndex := c.IndexTable[currentIndex]

		// if currentIndex < len(c.IndexTable)-1 {
		// 	nextIndex = c.IndexTable[currentIndex+1]
		// }

		// if nextIndex.TimeTM == tagIndex.TimeTM {
		// 	tagIndex = nextIndex
		// 	currentIndex++
		// }
		c.currentIndex = currentIndex + 1
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
