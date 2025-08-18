package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// parseNumbers разбирает строку чисел, разделённых запятыми
func parseNumbers(numbersStr string) ([]float64, error) {
	parts := strings.Split(numbersStr, ",")
	numbers := make([]float64, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		num, err := strconv.ParseFloat(trimmed, 64)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга числа '%s': %w", trimmed, err)
		}
		numbers = append(numbers, num)
	}

	if len(numbers) == 0 {
		return nil, fmt.Errorf("не введено ни одного числа")
	}

	return numbers, nil
}

// calculateSum вычисляет сумму чисел
func calculateSum(numbers []float64) float64 {
	sum := 0.0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// calculateAverage вычисляет среднее значение
func calculateAverage(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	return calculateSum(numbers) / float64(len(numbers))
}

// calculateMedian вычисляет медиану
func calculateMedian(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0
	}

	// Создаём копию для сортировки, чтобы не изменять исходный слайс
	sorted := make([]float64, len(numbers))
	copy(sorted, numbers)
	sort.Float64s(sorted)

	n := len(sorted)
	if n%2 == 1 {
		// Нечётное количество элементов
		return sorted[n/2]
	}

	// Чётное количество элементов
	return (sorted[n/2-1] + sorted[n/2]) / 2
}

// performCalculation выполняет вычисление в зависимости от операции
func performCalculation(operation string, numbers []float64) (float64, error) {
	switch strings.ToUpper(operation) {
	case "SUM":
		return calculateSum(numbers), nil
	case "AVG":
		return calculateAverage(numbers), nil
	case "MED":
		return calculateMedian(numbers), nil
	default:
		return 0, fmt.Errorf("неизвестная операция: %s. Доступные операции: SUM, AVG, MED", operation)
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	// Запрашиваем операцию
	fmt.Print("Введите операцию (SUM, AVG, MED): ")
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "Ошибка чтения операции")
		os.Exit(1)
	}
	operation := strings.TrimSpace(scanner.Text())

	// Запрашиваем числа
	fmt.Print("Введите числа через запятую: ")
	if !scanner.Scan() {
		fmt.Fprintln(os.Stderr, "Ошибка чтения чисел")
		os.Exit(1)
	}
	numbersStr := strings.TrimSpace(scanner.Text())

	// Проверяем на ошибки сканирования
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения ввода: %v\n", err)
		os.Exit(1)
	}

	// Парсим числа
	numbers, err := parseNumbers(numbersStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка парсинга чисел: %v\n", err)
		os.Exit(1)
	}

	// Выполняем вычисление
	result, err := performCalculation(operation, numbers)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка вычисления: %v\n", err)
		os.Exit(1)
	}

	// Выводим результат
	fmt.Printf("Результат: %.2f\n", result)
}
