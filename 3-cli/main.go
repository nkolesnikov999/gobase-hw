package main

import (
	"cli/api"
	"cli/bins"
	"cli/config"
	"cli/file"
	"cli/storage"
	"fmt"
	"log"
)

func main() {
	// Демонстрация работы с конфигурацией
	fmt.Println("=== Демонстрация Config и API ===")
	cfg := config.LoadFromEnvFile(".env")

	fmt.Println("Загруженная конфигурация:")
	fmt.Printf("KEY: %s\n", cfg.GetByKey("KEY"))

	// Создаем API сервис с конфигурацией
	apiService := api.NewAPIService(cfg)

	// Запускаем API сервис
	if err := apiService.Start(); err != nil {
		log.Fatalf("Ошибка запуска API: %v", err)
	}

	// Демонстрация использования API
	response := apiService.HandleRequest()
	fmt.Printf("API Response: %s\n", response)

	fmt.Println("\n=== Демонстрация работы с Bins ===")
	// Создаем пример данных
	bin1 := bins.NewBin("bin-001", "My First Bin", false)
	bin2 := bins.NewBin("bin-002", "My Second Bin", true)

	binList := bins.NewBinList()
	binList.AddBin(bin1)
	binList.AddBin(bin2)

	filename := "bins_data.json"

	// Демонстрация работы с JSON storage через dependency injection
	fmt.Println("=== Демонстрация JSON Storage (через DI) ===")
	jsonStorage := storage.NewStorage()
	jsonService := bins.NewService(jsonStorage)

	fmt.Println("Сохраняем bin list в JSON файл:", filename)
	err := jsonService.SaveBins(binList, filename)
	if err != nil {
		log.Fatalf("Ошибка при сохранении: %v", err)
	}
	fmt.Println("Данные успешно сохранены!")

	fmt.Println("Загружаем bin list из JSON файла:", filename)
	loadedBinList, err := jsonService.LoadBins(filename)
	if err != nil {
		log.Fatalf("Ошибка при загрузке: %v", err)
	}
	fmt.Printf("Загружено %d bins:\n", len(loadedBinList))
	for i, bin := range loadedBinList {
		fmt.Printf("  %d. ID: %s, Name: %s, Private: %t\n", i+1, bin.ID, bin.Name, bin.Private)
	}

	// Демонстрация работы с file storage через dependency injection
	fmt.Println("\n=== Демонстрация File Storage (заглушка через DI) ===")
	fileStorage := file.NewFileHandler()
	fileService := bins.NewService(fileStorage)

	fmt.Println("Тестируем файловое хранилище (заглушка):")
	err = fileService.SaveBins(binList, "test_file.dat")
	if err != nil {
		log.Printf("Ошибка при сохранении в file storage: %v", err)
	}

	_, err = fileService.LoadBins("test_file.dat")
	if err != nil {
		log.Printf("Ошибка при загрузке из file storage: %v", err)
	}

	// Демонстрация дополнительных функций file пакета
	fmt.Println("\n=== Дополнительные функции file пакета ===")
	fmt.Printf("Файл '%s' является JSON файлом: %t\n", filename, file.IsJSONFile(filename))
	fmt.Printf("Файл 'test.txt' является JSON файлом: %t\n", file.IsJSONFile("test.txt"))

	// Читаем содержимое JSON файла
	fmt.Println("\nЧитаем содержимое файла:")
	content, err := file.ReadFile(filename)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
	} else {
		fmt.Printf("Содержимое файла %s:\n%s\n", filename, content)
	}
}
