package main

import (
	"fmt"
	"unsafe"
)

func main() {
	i := unsafe.Sizeof(int32(0))
	fmt.Println("int32: ", i)

	rawi := unsafe.Sizeof(0)
	fmt.Println("raw int32:", rawi)

	i64 := unsafe.Sizeof(int64(0))
	fmt.Println("int64: ", i64)

	str := unsafe.Sizeof(string("HelloWorld"))
	fmt.Println("string: ", str)

	rawStr := unsafe.Sizeof("Hello")
	fmt.Println("raw string: ", rawStr)
}
