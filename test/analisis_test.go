package test_test

import (
	"os"
	"testing"

	"dav_converter/pkg/analysis"
)

func TestAnalysisMP4(t *testing.T) {
	filePath := "/run/media/dmitry/Хранилище/test/Арлюк-Васильевка_prm/Cam5.mp4" // "/run/media/dmitry/Хранилище/test/Новос-ЛК-Кемерово-Юрга 147-229/Cam1.mp4"

	file, err := os.Open(filePath)
	if err != nil {
		t.Error(err)
	}

	analisis := analysis.NewAnalysisMP4(file)

	analisis.Analysis()
}
