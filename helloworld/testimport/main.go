package main

import (
	"fmt"
	tp "mygo/helloworld/testpackage"
)

func main() {
	fmt.Println("test import")
	tp.Test()

	i := tp.Add()
	fmt.Println("add:", i)
}
