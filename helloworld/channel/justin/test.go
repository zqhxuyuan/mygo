package main

import (
	"fmt"
)

func main() {
	c := make(chan int)

	select {
	case c <- 1:
		fmt.Println("put 1")
		// default:
		// 	fmt.Println("default")
	}

	fmt.Println("after put")
}
