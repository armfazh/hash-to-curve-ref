package math

import (
	"math/big"
)

// Elt is a prime field element.
type Elt struct{ *big.Int }

func (e Elt) String() string { return "0x" + e.Text(16) }

// Fp implements prime field arithmetic.
type Fp struct{ *big.Int }

// NewElt returns an element reduced modulo p.
func (f Fp) NewElt(b []byte) Elt {
	i := new(big.Int).SetBytes(b)
	i.Mod(i, f.Int)
	return Elt{i}
}

// Size returns the size in bits of the prime modulus.
func (f Fp) Size() uint { return uint(f.Int.BitLen()) }

// Add performs modular addition.
func (f Fp) Add(z, x, y Elt) { z.Int.Add(x.Int, y.Int); z.Int.Mod(z.Int, f.Int) }

// Sub performs modular subtraction.
func (f Fp) Sub(z, x, y Elt) { z.Int.Sub(x.Int, y.Int); z.Int.Mod(z.Int, f.Int) }

// Mul performs modular multiplication.
func (f Fp) Mul(z, x, y Elt) { z.Int.Mul(x.Int, y.Int); z.Int.Mod(z.Int, f.Int) }

// Sqr performs modular squaring.
func (f Fp) Sqr(z, x Elt) { z.Int.Mul(x.Int, x.Int); z.Int.Mod(z.Int, f.Int) }

// Inv performs modular inversion.
func (f Fp) Inv(z, x Elt) { z.Int.ModInverse(x.Int, f.Int) }

// EltFq is a prime field element.
type EltFq []Elt

func (e EltFq) String() string {
	str := "[\n"
	for _, ei := range e {
		str += ei.String() + ",\n"
	}
	return str + "]"
}

func gfp(prime string) Fp {
	p, _ := new(big.Int).SetString(prime, 0)
	if !p.ProbablyPrime(5) {
		panic("p is not a prime number")
	}
	return Fp{p}
}

// Fq implements arithmetic of an extension of a prime field.
type Fq struct {
	P Fp   // Modulus
	M uint // Extension degree
}

// IsPrime returns true if f is not an extension field
func (f Fq) IsPrime() bool { return f.M == 1 }

// GF creates a finite field
func GF(p string, m uint) Fq {
	if m == 0 {
		panic("m must be larger than zero")
	}
	return Fq{
		P: gfp(p),
		M: m,
	}
}
