package main

import (
	"fmt"
	"math"
)

type Vertex struct {
	X, Y float64
}

type user struct {
	name string
	pass string
}

//指向类型值的指针
func (v *Vertex) Abs1() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

//类型值
func (v Vertex) Abs2() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

//指向类型值的指针
func (v *Vertex) Scale1(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

//类型值
func (v Vertex) Scale2(f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

//带指针参数的函数必须接受一个指针
func ScaleFunc(v *Vertex, f float64) {
	v.X = v.X * f
	v.Y = v.Y * f
}

//接受一个值作为参数的函数必须接受一个指定类型的值
func AbsFunc(v Vertex) float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func main() {
	//类型值
	v := Vertex{3, 4}
	v.Scale1(10)
	fmt.Println(v, v.Abs2)

	v.Scale2(10)
	fmt.Println(v, v.Abs2)

	v2 := &v
	(*v2).Scale1(10)
	fmt.Println(*v2, v2.Abs2)

	(*v2).Scale2(10)
	fmt.Println(*v2, v2.Abs2)

	fmt.Println("==========")

	u1 := user{"zqh", "111"}
	u2 := &user{"zzz", "222"}
	u1.notify()
	u2.notify()

	u1.changePass("1111")
	u2.changePass("2222")
	u1.notify()
	u2.notify()

	//以指针为接收者的方法被调用时，接收者既能为值又能为指针
	//以值为接收者的方法被调用时，接收者既能为值又能为指针
	//以..为接收者的方法，指的是方法的声明。接收者指的是调用者（v,&v,*v等）
	(&u1).changePass("11111")
	(&u1).notify()

	(*u2).changePass("22222")
	(*u2).notify()

	(*(&u1)).changePass("111111")
	(*(&u1)).notify()

	(&(*u2)).changePass("222222")
	(&(*u2)).notify()

	fmt.Println("==============")
	//定义一个接口类型的变量
	var a Abser
	fi := MyFloat(-1)
	vi := Vertex{1, 1}
	// 将MyFloat类型赋值给接口类型
	a = fi
	fmt.Println("MyFloat as Abser:", a.Abs())

	//将*Vertex类型赋值给接口类型，注意：如果是a = vi，则编译器会报错
	a = &vi
	fmt.Println("*Vertex as Abser:", a.Abs())
}

//值接收者
func (u user) notify() {
	fmt.Println("notify user:", u)
}

//指针接收者，有两种场景会使用指针接收者：
//1.方法能够修改其接收者指向的值。
//2.避免在每次调用方法时复制该值。若值的类型为大型结构体时，这样做会更加高效。
func (u *user) changePass(pass string) {
	u.pass = pass
}

type Abser interface {
	Abs() float64
}

type MyFloat float64

//MyFloat实现了Abser接口
func (f MyFloat) Abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

//*Vertex实现了Abser接口，但是Vertex并没有实现Abser接口
func (v *Vertex) Abs() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}
