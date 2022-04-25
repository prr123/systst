package main

import "fmt"

func main() {
	c := make(chan int)

	fmt.Printf("type of `c` is %T\n", c)
	fmt.Printf("value of `c` is %v\n", c)
}
