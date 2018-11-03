package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	//c = nil

	go func() {
		for {
			<-c
			fmt.Println("get 1")
		}
	}()

	go func() {
		for {
			select {
			case c <- 1:
				fmt.Println("put 1")
			default:
				//fmt.Println("default")
			}
		}
	}()

	time.Sleep(400 * time.Microsecond)

	c = nil
	fmt.Println("nill channel")

	time.Sleep(10 * time.Millisecond)
}
