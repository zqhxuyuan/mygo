package main

import (
	"fmt"
	"time"
)

// 接口类型
type I interface {
	M()
}

// 结构类型
type T struct {
	S string
}

// 结构类型实现了接口类型的所有方法
func (t *T) M() {
	//即便接口内的具体值为 nil，方法仍然会被 nil 接收者调用。
	if t == nil {
		fmt.Println("<nil>")
		return
	}
	fmt.Println(t.S)
}

// 自定义类型
type F float64

// 自定义类型实现了接口类型的所有方法
func (f F) M() {
	fmt.Println(f)
}

//接口值:在内部，接口值可以看做包含值(%v)和具体类型(%T)的元组：(value, type)
//接口值保存了一个具体底层类型的具体值。接口值调用方法时会执行其底层类型的同名方法。
func desc(i I) {
	fmt.Printf("(%v, %T)\n", i, i)
}

//指定了零个方法的接口值被称为"空接口"
func desc2(i interface{}) {
	fmt.Printf("(%v, %T)\n", i, i)
}

func main() {
	var i I

	//nil接口值:既不保存值也不保存具体类型。
	//desc(i)
	//i.M()

	t := &T{"hello"}
	f := F(10)

	i = t
	desc(i)
	i.M()

	i = f
	desc(i)
	i.M()

	var t1 *T
	i = t1
	desc(i)
	i.M()

	fmt.Println("=============")

	var i2 interface{}
	desc2(i2)

	i3 := 1
	desc2(i3)

	i4 := "hello"
	desc2(i4)

	var i5 interface{} = "hello"
	si, ok := i5.(string)
	fmt.Println(si, ok)

	fi, ok := i5.(float64)
	fmt.Println(fi, ok)

	do(1)
	do("hello")
	do(1.1)
	do(true)

	p := Person{"zqh", "111"}
	fmt.Println(p)

	error := &MyError{time.Now(), "opsss..."}
	fmt.Println(error)
}

func do(i interface{}) {
	switch v := i.(type) {
	case string:
		fmt.Println("string value:", v)
	case int:
		fmt.Println("int value:", v)
	default:
		fmt.Println("unknown val:", v)
	}
}

type Person struct {
	Name string
	Pass string
}

func (p Person) String() string {
	return fmt.Sprintf("%v (%v)", p.Name, p.Pass)
}

type MyError struct {
	when time.Time
	what string
}

func (e *MyError) Error() string {
	return fmt.Sprintf("at %v, %s", e.when, e.what)
}
