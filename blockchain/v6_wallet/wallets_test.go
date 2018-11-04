package main

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestChecksum(t *testing.T) {
	data := []byte("hello")
	res := checksum2(data)
	fmt.Println("checksum:", res)
}

func checksum2(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	fmt.Println("length:", len(secondSHA))
	fmt.Println("sha2:", secondSHA)

	return secondSHA[:addressChecksumLen]
}

func TestRangeMap(t *testing.T) {
	test := make(map[string]int)
	test["A"] = 1
	test["B"] = 2

	// just key
	for data := range test {
		fmt.Println(data)
	}

	// both key and value
	for i, data := range test {
		fmt.Println(i, data)
	}
}
