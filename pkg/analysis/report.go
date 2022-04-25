package analysis

import (
	"bufio"
	"fmt"
	"os"
)

func ReportAnalysis(folders []string) error {
	file, err := os.Create("report.txt")
	if err != nil {
		return fmt.Errorf("ошибка создания файла отчета: %w", err)
	}

	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, filePath := range folders {
		writer.WriteString(filePath)
		writer.WriteString("\n")
	}

	writer.Flush()
	return nil
}
