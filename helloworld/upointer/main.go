package main

import (
	"fmt"
	"unsafe"
)

func main() {
	u := new(user)
	fmt.Println(*u)

	pName := (*string)(unsafe.Pointer(u))
	*pName = "张三"

	// use offset by age
	// pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + unsafe.Offsetof(u.age)))
	// *pAge = 20

	// use interval from name pointer
	pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + unsafe.Sizeof(*pName)))
	*pAge = 21

	fmt.Println(*u)
}

type user struct {
	name string
	age  int
}
