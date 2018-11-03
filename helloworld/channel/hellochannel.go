package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	c := make(chan int)
	go printer(c)
	wg.Add(1)

	for i := 1; i <= 5; i++ {
		c <- i
	}

	close(c)
	wg.Wait()
	fmt.Println("channel Done...")

	//对切片中的数进行求和，将任务分配给两个 Go 程。一旦两个 Go 程完成了它们的计算，它就能算出最终的结果。
	fmt.Println("===============")
	arr := []int{1, 2, 3, 3, 4, 54, 5, 65, 6, 6, 7}
	sc := make(chan int)
	go sum(arr[:len(arr)/2], sc)
	go sum(arr[len(arr)/2:], sc)
	x, y := <-sc, <-sc
	fmt.Println("x+y:", x, y, x+y)

	fmt.Println("===============")
	testFibByChanel()

	testMutext()

	testPipie()
}

// channel
var wg sync.WaitGroup

func printer(ch chan int) {
	for i := range ch {
		fmt.Printf("channel Received: %d\n", i)
	}
	wg.Done()
}

func sum(s []int, c chan int) {
	sum := 0
	for _, i := range s {
		sum += i
	}
	c <- sum
}

func fib(c chan int, q chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x: //将x放入c channel中。如果调用者没有从c中读取，就会阻塞
			x, y = y, x+y
		case <-q:
			fmt.Println("quit")
			return
		}
	}
}

func testFibByChanel() {
	c := make(chan int)
	q := make(chan int)

	go func() {
		for i := 0; i < 10; i++ {
			//从c channel中读取10次数据
			fmt.Println(<-c)
		}
		//往q channel中放入数据
		q <- 0
	}()

	fib(c, q)
}

type SafeCounter struct {
	v   map[string]int
	mux sync.Mutex
}

func (c *SafeCounter) inc(key string) {
	c.mux.Lock()
	c.v[key]++
	c.mux.Unlock()
}

func (c *SafeCounter) value(key string) int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.v[key]
}

func testMutext() {
	c := SafeCounter{v: make(map[string]int)}
	for i := 0; i < 1000; i++ {
		go c.inc("key")
	}
	time.Sleep(time.Second)
	fmt.Println(c.value("key"))
}

// https://blog.golang.org/pipelines
func testPipie() {
	in := gen(1, 2, 3)
	out := sq(in)
	for o := range out {
		fmt.Println(o)
	}

	// Set up the pipeline and consume the output.
	for n := range sq(sq(gen(1, 2, 3))) {
		fmt.Println(n) // 16 then 81
	}
	fmt.Println("merge...")

	in1 := gen(1, 2, 3)
	in2 := gen(4, 5, 6)
	m1 := merge(in1, in2)
	for m := range m1 {
		fmt.Println(m)
	}
}

func gen(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	sq := make(chan int)
	go func() {
		for n := range in {
			sq <- n * n
		}
		close(sq)
	}()
	return sq
}

// merge multi input channels to one outpout channel
// c1 <- 1,2,3
// c2 <- 4,5,6
// merge: <- 1,2,3,4,5,6
func merge(cs ...<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	// each input channel, put to output channel
	// when this input channel complete, it's Done!
	output := func(c <-chan int) {
		for i := range c { // loop from channel, only value
			out <- i
		}
		wg.Done()
	}
	wg.Add(len(cs))

	// all input channel(this is an array, so for range return k,v)
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
