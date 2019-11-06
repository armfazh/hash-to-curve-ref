package field

import (
	"io"
	"math/big"
)

// Field is
type Field interface {
	E() Elt // Constructor of Elements
	EltFromList([]*big.Int) Elt
	Rand(r io.Reader) Elt // Constructor of Elements at Random
	Zero() Elt            // Returns the zero element
	One() Elt             // Returns the one element
	Ext() uint            // Extension degree
	BitLen() int          // Bit length of modulus
	P() *big.Int
	hasArith
	hasCmov
	hasSqrt
	hasSign
}

type hasCmov interface{ CMov(x, y Elt, b bool) Elt }
type hasSign interface{ Sgn0(x Elt) int }
type hasSqrt interface {
	Sqrt(x Elt) Elt
	IsSquare(x Elt) bool
}
type hasArith interface {
	Neg(x Elt) Elt
	Add(x, y Elt) Elt
	Sub(x, y Elt) Elt
	Mul(x, y Elt) Elt
	Sqr(x Elt) Elt
	Inv0(x Elt) Elt
}

// Elt is
type Elt interface {
	// IsZero() bool
}

// NewFromID is
func NewFromID(id Prime) Field { return getFromID(id) }

// NewGF is
func NewGF(p string, m uint, name string) Field {
	if !(m == 1 || m == 2) {
		panic("not implemented")
	}
	modulus := prime{
		name: name,
		p:    bigFromString(p),
	}
	if !modulus.p.ProbablyPrime(5) {
		panic("p=" + p + " is not prime")
	}
	switch m {
	case 1:
		return newFp(modulus)
	case 2:
		return newFp2(modulus)
	}
	return nil
}
