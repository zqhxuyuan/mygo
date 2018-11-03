package main

import "fmt"

func main() {
	sliceTest()
}

func sliceTest() {
	arr := []int{1, 2, 3, 4, 5}
	s := arr[:]
	for e := range s {
		fmt.Println(s[e])
	}
	fmt.Println("------------")
	s1 := make([]int, 3)
	for e := range s1 {
		fmt.Println(s1[e])
	}
	fmt.Println("------------")
	s2 := arr[:0]
	for e := range s2 {
		fmt.Println(s[e])
	}
}

func twoDimensionArray() {
	/* 数组 - 5 行 2 列*/
	var a = [][]int{{0, 0}, {1, 2}, {2}, {3, 6}, {4, 8}}
	var i, j int

	/* 输出数组元素 */
	for i = 0; i < len(a); i++ {
		for j = 0; j < len(a[i]); j++ {
			fmt.Printf("a[%d][%d] = %d\n", i, j, a[i][j])
		}
	}
}
