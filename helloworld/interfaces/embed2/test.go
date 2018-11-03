package main

import (
	"fmt"
)

type a struct {
	y int
	B // 匿名定义
}

type B struct {
	z int
	C // 匿名定义
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
	d := &D{x: 3}

	b2 := B{z: 1}
	b2.C = d

	a2 := a{y: 2}
	a2.B = b2

	a2.Tick()
}
