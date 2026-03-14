package main

import "fmt"

func main() {
	fmt.Println("Welcome to the Go Calculator")
	sum := add(10, 5)
	difference := sub(10, 4)
	quotient, err := div(10, 2)
	product := mul(2, 5)
	
	fmt.Println("10 + 5 =", sum)
	fmt.Println("10 - 4 =", difference)
	fmt.Println("2 * 5 =", product)
	if err != nil{
		fmt.Println(err)
	} else {
		fmt.Println("10 / 2 =", quotient)
	}
	quotient, err = div(10, 0)
	if err != nil{
		fmt.Println(err)
	} else {
		fmt.Println("10 / 0 =", quotient)
	}
}
