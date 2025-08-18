package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type exchangeMap = map[string]map[string]float64

// Константы для валют
const (
	USD = "USD"
	EUR = "EUR"
	RUB = "RUB"
)

func main() {

	// Курсы валют в виде map с двумя ключами
	var exchangeRates = exchangeMap{
		USD: {
			EUR: 0.85,
			RUB: 80.10,
		},
		EUR: {
			USD: 1.18,
			RUB: 94.24,
		},
		RUB: {
			USD: 0.0125,
			EUR: 0.0106,
		},
	}
	num, origCur, targetCur := inputCur()
	calculation(&exchangeRates, num, origCur, targetCur)
}

func inputCur() (float64, string, string) {
	fmt.Println("=== Конвертер валют ===")

	// Ввод исходной валюты
	origCur := inputSourceCurrency()

	// Ввод суммы
	num := inputAmount()

	// Ввод целевой валюты
	targetCur := inputTargetCurrency(origCur)

	return num, origCur, targetCur
}

// getCurrencyName возвращает имя валюты по номеру
func getCurrencyName(choice int) string {
	switch choice {
	case 1:
		return USD
	case 2:
		return EUR
	case 3:
		return RUB
	default:
		return ""
	}
}

// isWholeNumber проверяет, является ли число целым
func isWholeNumber(num float64) bool {
	return num == float64(int(num))
}

// inputSourceCurrency запрашивает у пользователя исходную валюту
func inputSourceCurrency() string {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Выберите исходную валюту:")
		fmt.Println("1 - USD")
		fmt.Println("2 - EUR")
		fmt.Println("3 - RUB")
		fmt.Print("Введите номер (1-3): ")

		if !scanner.Scan() {
			fmt.Println("Ошибка чтения ввода. Попробуйте снова.")
			continue
		}

		input := strings.TrimSpace(scanner.Text())
		choice, err := strconv.Atoi(input)

		if err != nil || choice < 1 || choice > 3 {
			fmt.Println("Ошибка: введите число от 1 до 3")
			continue
		}

		return getCurrencyName(choice)
	}
}

// inputAmount запрашивает у пользователя сумму для конвертации
func inputAmount() float64 {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите сумму (целое положительное число): ")

		if !scanner.Scan() {
			fmt.Println("Ошибка чтения ввода. Попробуйте снова.")
			continue
		}

		input := strings.TrimSpace(scanner.Text())
		amount, err := strconv.ParseFloat(input, 64)

		if err != nil {
			fmt.Println("Ошибка: введите корректное число")
			continue
		}

		if amount <= 0 {
			fmt.Println("Ошибка: сумма должна быть положительным числом")
			continue
		}

		// Проверяем, что это целое число
		if !isWholeNumber(amount) {
			fmt.Println("Ошибка: введите целое число")
			continue
		}

		return amount
	}
}

// inputTargetCurrency запрашивает у пользователя целевую валюту
func inputTargetCurrency(sourceCurrency string) string {
	scanner := bufio.NewScanner(os.Stdin)

	// Создаем список доступных валют (исключая исходную)
	var availableCurrencies []string
	var currencyMap = make(map[int]string)
	index := 1

	currencies := []string{USD, EUR, RUB}
	for _, currency := range currencies {
		if currency != sourceCurrency {
			availableCurrencies = append(availableCurrencies, currency)
			currencyMap[index] = currency
			index++
		}
	}

	for {
		fmt.Println("Выберите целевую валюту:")
		for i, currency := range availableCurrencies {
			fmt.Printf("%d - %s\n", i+1, currency)
		}
		fmt.Printf("Введите номер (1-%d): ", len(availableCurrencies))

		if !scanner.Scan() {
			fmt.Println("Ошибка чтения ввода. Попробуйте снова.")
			continue
		}

		input := strings.TrimSpace(scanner.Text())
		choice, err := strconv.Atoi(input)

		if err != nil || choice < 1 || choice > len(availableCurrencies) {
			fmt.Printf("Ошибка: введите число от 1 до %d\n", len(availableCurrencies))
			continue
		}

		return currencyMap[choice]
	}
}

func calculation(exchangeRates *exchangeMap, num float64, origCur string, targetCur string) {
	// Получаем курс обмена из map
	rate, exists := (*exchangeRates)[origCur][targetCur]
	if !exists {
		fmt.Printf("Ошибка: неподдерживаемое направление конвертации %s -> %s\n", origCur, targetCur)
		return
	}

	// Выполняем конвертацию
	result := num * rate

	// Выводим результат с точностью до 2 знаков после запятой
	fmt.Printf("Результат конвертации: %.0f %s = %.2f %s\n", num, origCur, result, targetCur)
}
