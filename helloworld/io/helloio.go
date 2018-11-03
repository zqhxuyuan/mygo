package main

import (
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("hello, world!")
	b := make([]byte, 8)

	for {
		n, error := r.Read(b)
		if error == io.EOF {
			break
		} else {
			fmt.Printf("n = %v err = %v b = %v, read:%q\n", n, error, b, b[:n])
		}
	}
	fmt.Println("complete.")

}
