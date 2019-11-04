package main

import "fmt"
import "github.com/armfazh/hash-to-curve-ref/h2c/math"

func main() {
	fmt.Println("Test Vectors")
	Fp := math.NEWGF("103", 1)
	fmt.Println(Fp)
	a, b := Fp.Elt(), Fp.Elt()
	fmt.Println(a)
	fmt.Println(b)
	c := Fp.Add(a, b)
	fmt.Println(c)

	Fp2 := math.NEWGF("103", 2)
	fmt.Println(Fp2)
	x, y := Fp2.Elt(), Fp2.Elt()
	fmt.Println(x)
	fmt.Println(y)
	z := Fp2.Add(x, y)
	fmt.Println(z)
}
