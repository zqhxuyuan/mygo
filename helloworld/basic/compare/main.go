package main

import (
	"bytes"
	"fmt"
)

func main() {
	a := []byte("a")
	b := []byte("b")
	c := []byte("c")
	fmt.Println("", bytes.Compare(b, a)) // 1
	fmt.Println("", bytes.Compare(b, c)) // -1
}
