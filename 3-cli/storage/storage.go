package storage

import (
	"cli/bins"
	"encoding/json"
	"fmt"
	"os"
)

// SaveBinListToFile сохраняет список bin в JSON файл
func SaveBinListToFile(binList bins.BinList, filename string) error {
	// Сериализуем BinList в JSON
	jsonData, err := json.MarshalIndent(binList, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal bin list: %w", err)
	}

	// Записываем JSON данные в файл
	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filename, err)
	}

	return nil
}

// LoadBinListFromFile загружает список bin из JSON файла
func LoadBinListFromFile(filename string) (bins.BinList, error) {
	var binList bins.BinList

	// Читаем файл
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Десериализуем JSON в BinList
	err = json.Unmarshal(data, &binList)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from file %s: %w", filename, err)
	}

	return binList, nil
}
