package analysis

import (
	"fmt"
	"os"

	"github.com/alfg/mp4"
)

type AnalysisMP4 struct {
	FrameCount  uint
	DurationSec float64
	file        *os.File
}

func NewAnalysisMP4(file *os.File) *AnalysisMP4 {
	return &AnalysisMP4{
		file: file,
	}
}

func (analysis *AnalysisMP4) Analysis() (uint, error) {
	info, err := analysis.file.Stat()
	if err != nil {
		return 0, fmt.Errorf("Analysis Stat: %w", err)
	}
	size := info.Size()

	mp4, err := mp4.OpenFromReader(analysis.file, size)
	if err != nil {
		return 0, fmt.Errorf("Analysis mp4.OpenFromReader: %w", err)
	}

	analysis.DurationSec = float64(mp4.Moov.Traks[0].Tkhd.Duration) / 1000
	analysis.FrameCount = uint(analysis.DurationSec * 25)
	// fmt.Printf("Len: %d\n", len(mp4.Moov.Traks))
	// fmt.Printf("Duration: %+v, ModificationTime: %+v\n", mp4.Moov.Traks[len(mp4.Moov.Traks)-1].Tkhd.Duration, mp4.Moov.Traks[len(mp4.Moov.Traks)-1].Tkhd.Duration)

	return analysis.FrameCount, nil
}
