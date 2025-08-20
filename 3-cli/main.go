package main

import (
	"cli/bins"
	"cli/file"
	"cli/storage"
	"fmt"
	"log"
)

func main() {
	// Создаем пример данных
	bin1 := bins.NewBin("bin-001", "My First Bin", false)
	bin2 := bins.NewBin("bin-002", "My Second Bin", true)

	binList := bins.NewBinList()
	binList.AddBin(bin1)
	binList.AddBin(bin2)

	// Пример работы с storage пакетом
	filename := "bins_data.json"

	// Сохраняем bin list в JSON файл
	fmt.Println("Сохраняем bin list в файл:", filename)
	err := storage.SaveBinListToFile(binList, filename)
	if err != nil {
		log.Fatalf("Ошибка при сохранении: %v", err)
	}
	fmt.Println("Данные успешно сохранены!")

	// Загружаем bin list из JSON файла
	fmt.Println("Загружаем bin list из файла:", filename)
	loadedBinList, err := storage.LoadBinListFromFile(filename)
	if err != nil {
		log.Fatalf("Ошибка при загрузке: %v", err)
	}
	fmt.Printf("Загружено %d bins:\n", len(loadedBinList))
	for i, bin := range loadedBinList {
		fmt.Printf("  %d. ID: %s, Name: %s, Private: %t\n", i+1, bin.ID, bin.Name, bin.Private)
	}

	// Пример работы с file пакетом
	fmt.Println("\nПроверяем расширение файла:")
	fmt.Printf("Файл '%s' является JSON файлом: %t\n", filename, file.IsJSONFile(filename))
	fmt.Printf("Файл 'test.txt' является JSON файлом: %t\n", file.IsJSONFile("test.txt"))

	// Читаем содержимое JSON файла
	fmt.Println("\nЧитаем содержимое файла:")
	content, err := file.ReadFile(filename)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}
	fmt.Printf("Содержимое файла %s:\n%s\n", filename, content)
}
