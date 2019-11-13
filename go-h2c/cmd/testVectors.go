package main

import (
	crand "crypto/rand"
	"crypto/sha256"
	"fmt"
	mrand "math/rand"

	"github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

func main() {
	_ = crand.Reader
	_ = mrand.New(mrand.NewSource(5))
	fmt.Println("Test Vectors")

	H := sha256.New
	msg := "hello"
	dst := "world"
	// fmt.Printf("msg: %v\n", msg)
	// fmt.Printf("dst: %v\n", dst)
	//
	// F1 := field.NewGF("103", 1, "2^7-25")
	// a1 := field.HashToField(H, []byte(msg), []byte(dst), 0, 3, F1)
	// fmt.Println(F1)
	// fmt.Printf("H(msg): %v\n", a1)
	//
	F2 := field.NewGF("103", 2, "2^7-25")
	a2 := field.HashToField([]byte(msg), []byte(dst), 0, H, F2, 3)
	fmt.Println(F2)
	fmt.Printf("H(msg): %v\n", a2)

	F3 := field.NewFromID(field.P25519)
	a3 := field.HashToField([]byte(msg), []byte(dst), 0, H, F3, 3)
	fmt.Println(F3)
	fmt.Printf("H(msg): %v\n", a3)
}
