package testpackage

import (
	"fmt"
)

func Test() {
	fmt.Println("test method")
}

func Add() int {
	i := 1
	i += 1
	return i
}
