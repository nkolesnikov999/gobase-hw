package main

import "fmt"

func main() {
	const USD_EUR = 0.85
	const USD_RUB = 80.10
	const EUR_RUB = USD_RUB / USD_EUR
	num, origCur, targetCur := inputCur()
	calculation(num, origCur, targetCur)
}

func inputCur() (float64, string, string) {
	var num float64
	var origCur string
	var targetCur string
	fmt.Print("Введите сумму: ")
	fmt.Scan(&num)
	fmt.Print("Введите исходную валюту: ")
	fmt.Scan(&origCur)
	fmt.Print("Введите целевую: ")
	fmt.Scan(&targetCur)
	return num, origCur, targetCur
}

func calculation(num float64, origCur string, targetCur string) {
	fmt.Println(num, origCur, targetCur)
}
