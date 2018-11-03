package main

import (
	"fmt"
)

func main() {
	c1 := make(chan int)
	c2 := make(chan int)
	c1 = c2

	// 往c1中放数据，可以从c2中取出数据
	go func() {
		c1 <- 1
	}()
	v1 := <-c2
	fmt.Println("v1:", v1)

	// 往c2中放数据，也可以从c1中取出数据
	go func() {
		c2 <- 2
	}()
	v2 := <-c1
	fmt.Println("v2:", v2)
}
