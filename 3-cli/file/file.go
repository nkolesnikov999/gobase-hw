package file

import (
	"cli/bins"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileHandler реализует интерфейс bins.Repository для файлового хранилища
type FileHandler struct{}

// NewFileHandler создает новый экземпляр FileHandler
func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

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

// Save сохраняет список bin в файл (заглушка, реализация интерфейса bins.Repository)
func (fh *FileHandler) Save(binList bins.BinList, filename string) error {
	// TODO: Реализовать сохранение bin list в файл
	fmt.Printf("File storage: saving %d bins to %s\n", len(binList), filename)
	return nil
}

// Load загружает список bin из файла (заглушка, реализация интерфейса bins.Repository)
func (fh *FileHandler) Load(filename string) (bins.BinList, error) {
	// TODO: Реализовать загрузку bin list из файла
	fmt.Printf("File storage: loading bins from %s\n", filename)
	return bins.NewBinList(), nil
}
