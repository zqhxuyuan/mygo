package main

import (
	"fmt"
)

type a struct {
	y int
	b B
}

type B struct {
	z int
	C
}

type C interface {
	Tick()
}

type D struct {
	x int
}

func (d *D) Tick() {
	fmt.Println("Tick")
}

func main() {
	fmt.Println("Hello")
	b1 := B{z: 1}
	a1 := a{y: 2, b: b1}
	fmt.Println(a1)
	//panic: runtime error: invalid memory address or nil pointer dereference
	//a1.b.Tick()

	d := &D{x: 3}
	b2 := B{z: 1}
	b2.C = d
	a2 := a{y: 2, b: b2}
	fmt.Println(a2)
	a2.b.Tick()

	//a2.Tick()
}
