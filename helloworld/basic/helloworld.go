package main

import (
	"fmt"
	"math"
	"math/cmplx"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

var (
	ToBe   bool       = false
	MaxInt uint64     = 1<<64 - 1
	z      complex128 = cmplx.Sqrt(-1 + 12i)
)

type Vertex struct {
	Name string
	Pass string
}

const (
	Big   = 1 << 100
	Small = Big >> 99
)

type MyFloat float64

func (f MyFloat) abs() float64 {
	if f < 0 {
		return float64(-f)
	}
	return float64(f)
}

func main() {
	defer fmt.Println("This is a test...")
	fmt.Println("hello go world!!!")
	fmt.Println("random number: ", rand.Intn(10))

	fmt.Println("add result:", add(1, 2))
	a, b := swap("hello", "world")
	fmt.Println("swap:", a, b)

	x, y := split(100)
	fmt.Println("split:", x, y)

	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)

	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println("sum:", sum)

	fmt.Println("pow:", pow(3, 2, 5))
	fmt.Println("sqr:", sqrt(4))

	useSwitch()
	deferStack()

	var p *int
	i, j, z := 1, 10, 100
	p = &i
	fmt.Println("pointer:", *p)
	*p = 2
	fmt.Println("new pointer:", i)
	p = &j
	fmt.Println(*p)
	pp := &z
	fmt.Println(pp)

	vertex := Vertex{
		Name: "zqh",
		Pass: "xxx",
	}
	vertex.Name = "Zqh"
	fmt.Println("vertex1:", vertex.Name, ":", vertex.Pass)
	fmt.Println(Vertex{"z", "qh"})

	//结构体指针，指针pv指向了一个结构体。可以通过(*pv).Name访问字段Name
	pv := &vertex
	pv.Name = "ZQH"
	(*pv).Name = "ZZZ"
	fmt.Println("vertex2:", vertex.Name, ":", vertex.Pass)

	pv3 := &Vertex{}
	pv3.Name = "zqh"
	fmt.Println("vertex3:", pv3, *pv3)

	useArray()
	useMap()
	useFuncValue()
	useMethod()

	f := MyFloat(-1)
	fmt.Println("myfloat:", f, f.abs())
}

func add(x, y int) int {
	return x + y
}

func swap(x, y string) (string, string) {
	return y, x
}

func split(sum int) (x, y int) {
	x = sum * 4 / 9
	y = sum - x
	return
}

func pow(x, n, limit float64) float64 {
	if v := math.Pow(x, n); v < limit {
		return v
	} else {
		fmt.Printf("%g > %g\n", v, limit)
	}
	return limit
}

func sqrt(x float64) float64 {
	z := float64(1)
	for i := 0; i < 10; i++ {
		z -= (z*z - x) / (2 * z)
		fmt.Printf("iterator %d, result:%g\n", i, z)
		if x/z == 0 {
			break
		}
	}
	return z
}

func useSwitch() {
	switch os := runtime.GOOS; os {
	case "darwin":
		fmt.Println("go on osx.")
		//fallthrough // auto break. use fallthrough to continue
	case "linux":
		fmt.Println("go on linux.")
	default:
		fmt.Println("go on ", os)
	}

	today := time.Now().Weekday()
	switch time.Sunday {
	case today + 0:
		fmt.Println("Today.")
	case today + 1:
		fmt.Println("Tomorrow.")
	case today + 2:
		fmt.Println("In two days.")
	default:
		fmt.Println("Too far away.")
	}

	t := time.Now()
	switch {
	case t.Hour() > 12:
		fmt.Println("Afternoon")
	case t.Hour() < 12:
		fmt.Println("Morning")
	default:
		fmt.Println("hey...")
	}
}

func deferStack() {
	fmt.Println("Counting...")

	for i := 0; i < 10; i++ {
		defer fmt.Println("defer ", i)
	}

	fmt.Println("Done...")
}

func useArray() {
	var i [4]int
	i[0] = 1
	fmt.Println(i)

	primes := [5]int{2, 3, 5, 7, 11}
	fmt.Println(primes)

	slice := primes[0:4]
	fmt.Println(slice)

	names := [4]string{"Spring", "Summer", "Fall", "Winter"}
	a := names[0:2]
	b := names[2:4]
	fmt.Println(names, len(names), cap(names))
	fmt.Println(a, len(a), cap(a))
	fmt.Println(b, len(b), cap(b))
	//更改切片的元素会修改其底层数组中对应的元素
	a[0] = "spring"
	b[0] = "fall"
	fmt.Println(names)

	ss := []struct {
		b int
		c bool
	}{
		{1, true},
		{2, false},
	}
	fmt.Println(ss)

	//[] 0 0
	var snil []string
	fmt.Println(snil, len(snil), cap(snil))

	//[0 0 0 0 0] 5 5
	sm := make([]int, 5)
	fmt.Println(sm, len(sm), cap(sm))

	//[] 0 5
	sm = make([]int, 0, 5)
	fmt.Println(sm, len(sm), cap(sm))

	sm = []int{1, 2, 3, 4, 5}

	//[1 2 3] 3 5
	sm = sm[0:3]
	fmt.Println(sm, len(sm), cap(sm))

	//[3] 1 3
	sm = sm[2:3]
	fmt.Println(sm, len(sm), cap(sm))

	// slice of slice
	board := [][]string{
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
		[]string{"_", "_", "_"},
	}
	board[0][0] = "X"
	board[2][2] = "O"
	board[1][2] = "X"
	board[1][0] = "O"
	board[0][2] = "X"
	for i := 0; i < len(board); i++ {
		fmt.Printf("%s\n", strings.Join(board[i], " "))
	}

	var s []int
	s = append(s, 0)
	fmt.Println("append 0:", s, len(s), cap(s))

	s = append(s, 1)
	fmt.Println("append 1:", s, len(s), cap(s))

	s = append(s, 2, 3, 4)
	fmt.Println("append *:", s, len(s), cap(s))

	// append two slice
	var s1 []int
	var s2 []int
	var s3 []int
	s1 = append(s1, 0, 1, 2)
	s2 = append(s2, 3, 4, 5)
	s3 = append(s1, s2...)
	fmt.Println(s1, s2, s3)
}

func useMap() {
	var m1 map[string]string
	m1 = make(map[string]string)
	m1["1"] = "one"
	fmt.Println(m1)

	m2 := make(map[string]string)
	m2["2"] = "two"
	fmt.Println(m2)

	m3 := make(map[string]Vertex)
	m3["3"] = Vertex{"zqh", "111"}

	m45 := map[string]Vertex{
		"4": {
			"z", "q",
		},
		"5": {
			"x", "y",
		},
	}

	mm := make(map[string]interface{})
	mm["1"] = m1
	mm["2"] = m2
	mm["3"] = m3
	mm["45"] = m45
	fmt.Println(mm)

	m := make(map[string]int)
	m["Answer"] = 42
	fmt.Println("The value:", m["Answer"])
	m["Answer"] = 48
	fmt.Println("The value:", m["Answer"])
	delete(m, "Answer")
	fmt.Println("The value:", m["Answer"])
	v, ok := m["Answer"]
	fmt.Println("The value:", v, "Present?", ok)

	str := "hello world hello go go"
	fmt.Println(wordCount(str))
}

func wordCount(s string) map[string]int {
	arr := strings.Split(s, " ")
	wc := make(map[string]int)
	for _, word := range arr {
		cnt, ok := wc[word]
		if ok {
			wc[word] = cnt + 1
		} else {
			wc[word] = 1
		}
	}
	return wc
}

func useFuncValue() {
	myF1 := func(x, y float64) float64 {
		return x*x + y*y
	}
	fmt.Println(myF1(1, 2))

	myF2 := func(x, y float64) float64 {
		return math.Sqrt(x*x + y*y)
	}

	fmt.Println(compute(myF1))
	fmt.Println(compute(myF2))

	pos := adder()
	fmt.Println(pos(1))
	fmt.Println(pos(2))

	neg := adder()
	fmt.Println(neg(-1))
	fmt.Println(neg(-2))

	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

func compute(fn func(float64, float64) float64) float64 {
	return fn(3, 4)
}

func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

func fibonacci() func() int {
	x, y := 0, 1
	return func() int {
		f := x
		x, y = y, x+y
		return f
	}
}

func useMethod() {
	v := Vertex{"zqh", "pass"}
	fmt.Println(printUser(v))

	fmt.Println(v.printUser2())
}

func printUser(v Vertex) string {
	return v.Name + "." + v.Pass
}
func (v Vertex) printUser2() string {
	return v.Name + "." + v.Pass
}
