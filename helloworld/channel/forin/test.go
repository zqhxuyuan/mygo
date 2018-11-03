package main

import (
	"fmt"
)

func main() {
	c := make(chan int)

	go func() {
		for {
			select {
			case c <- 1:
				fmt.Println("put 1")
			}
		}
	}()

	d := <-c
	fmt.Println(d)

	for index := 0; index < 10; index++ {
		d = <-c
		fmt.Println(index, d)
	}
}
