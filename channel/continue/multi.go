package main

import (
	"fmt"
)

func main() {
	c1 := make(chan int, 2)
	c2 := make(chan int, 2)
	defer close(c1)
	defer close(c2)

	for {
		fmt.Println("loop..")
		select {
		case c1 <- 1:
			i := <-c1
			fmt.Println("get c1", i)
			c2 <- i
		case i2 := <-c2:
			fmt.Println("get c2", i2)
			return
		}
	}
}
