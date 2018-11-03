package main

import (
	"fmt"
)

func main() {
	im := &meta{}
	fmt.Println(im)

	m0 := &meta{
		1,
		"meta0",
	}
	fmt.Println(m0)
	fmt.Println("----------------------")

	db := &db{
		"db",
		100,
		m0,
	}
	fmt.Println("db meta:", db.meta)
	fmt.Printf("db meta: %p\n", db.meta)
	fmt.Println("----------------------")

	tx := &tx{}
	tx.meta = &meta{}

	db.meta.copy(tx.meta)
	fmt.Println("tx meta:", tx.meta)
	fmt.Printf("tx meta: %p\n", tx.meta)
	fmt.Println("----------------------")
	tx.meta = &meta{
		2,
		"meta1",
	}
	fmt.Println("db meta:", db.meta)
	fmt.Println("tx meta:", tx.meta)
}

func (m *meta) copy(dest *meta) {
	*dest = *m
}

type db struct {
	name string
	size int
	meta *meta
}

type meta struct {
	index int
	name  string
}

type tx struct {
	txid int
	meta *meta
}
