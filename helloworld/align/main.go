package main

import (
	"fmt"
	"unsafe"
)

func init() {
	var w *W = new(W)
	fmt.Printf("size=%d\n", unsafe.Sizeof(*w))
	fmt.Println("b:", unsafe.Alignof(w.b))
	fmt.Println("i1:", unsafe.Alignof(w.i1))
	fmt.Println("i2:", unsafe.Alignof(w.i2))

	w2 := &W{}
	fmt.Printf("size=%d\n", unsafe.Sizeof(w2))
}

func main() {
	w := &W{}
	fmt.Println("b:", unsafe.Alignof(w.b))
	fmt.Println("i1:", unsafe.Alignof(w.i1))
	fmt.Println("i2:", unsafe.Alignof(w.i2))

	var w2 *W = new(W)
	fmt.Println("b:", unsafe.Alignof(w2.b))
	fmt.Println("i1:", unsafe.Alignof(w2.i1))
	fmt.Println("i2:", unsafe.Alignof(w2.i2))
}

type W struct {
	b  byte
	i1 int32
	i2 int64
}
