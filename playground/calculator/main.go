package main

import "fmt"

func add(a, b int) int{
	return a + b
}

func main() {
	fmt.Println("Welcome to the Go Calculator")
	sum := add(10, 5)
	fmt.Println("10 + 5 =", sum)
}
