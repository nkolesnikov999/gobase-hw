package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	// Channel from generator to squarer
	numbers := make(chan int)
	// Channel from squarer to main
	squares := make(chan int)

	// Random source
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Goroutine 1: generate 10 random numbers [0..100] and send them one by one
	go func() {
		defer close(numbers)
		data := make([]int, 10)
		for i := 0; i < len(data); i++ {
			data[i] = r.Intn(101)
		}
		for _, v := range data {
			numbers <- v
		}
	}()

	// Goroutine 2: receive numbers, square them, and send to main
	go func() {
		defer close(squares)
		for v := range numbers {
			squares <- v * v
		}
	}()

	// Collect all squared numbers and print
	result := make([]int, 0, 10)
	for x := range squares {
		result = append(result, x)
	}
	fmt.Println(result)
}
