package main

import "fmt"

func main() {
	fmt.Println("Welcome to the Go Calculator")
	sum := add(10, 5)
	sub := sub(10, 4)
	div := div(10, 2)
	mul := mul(2, 5)
	
	fmt.Println("10 + 5 =", sum)
	fmt.Println("10 - 4 =", sub)
	fmt.Println("10 / 2 =", div)
	fmt.Println("2 * 5 =", mul)
}
