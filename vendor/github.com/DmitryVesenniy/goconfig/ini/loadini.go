package ini

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	defaultIni = "config.ini"
)

func Load(filenames ...string) (map[string]string, error) {
	if len(filenames) == 0 {
		filenames = []string{defaultIni}
	}

	for _, _file := range filenames {
		fieldsIni, err := loadFile(_file)
		if err == nil {
			return fieldsIni, nil
		}
	}

	return map[string]string{}, fmt.Errorf("not found config file")
}

func loadFile(fileName string) (map[string]string, error) {
	file, err := os.Open(fileName)
	fields := make(map[string]string)
	if err != nil {
		return fields, fmt.Errorf("error open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		keyValueSplitted := strings.Split(line, "=")
		if len(keyValueSplitted) == 0 {
			keyValueSplitted = strings.Fields(line)
		}

		fieldName := strings.TrimSpace(keyValueSplitted[0])
		if fieldName == "" {
			continue
		}

		if len(keyValueSplitted) == 1 {
			fields[fieldName] = ""
		} else {
			fields[fieldName] = strings.TrimSpace(keyValueSplitted[1])
		}
	}

	return fields, scanner.Err()
}
