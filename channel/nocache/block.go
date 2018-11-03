package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	defer close(c)

	go func() {
		fmt.Println("prepare channel in")
		c <- 1
		fmt.Println("put channel ok")
	}()

	time.Sleep(3 * time.Second)

	fmt.Println("prepare get channel")
	i := <-c
	fmt.Println("get channel", i)

	time.Sleep(1 * time.Second)
}
