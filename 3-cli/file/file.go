package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ReadFile читает содержимое любого файла и возвращает его как строку
func ReadFile(filename string) (string, error) {
	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	return string(data), nil
}

// IsJSONFile проверяет, имеет ли файл расширение .json
func IsJSONFile(filename string) bool {
	// Получаем расширение файла и приводим к нижнему регистру
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".json"
}
