package field

import (
	"fmt"
	"io"
	"math/big"
)

// Field is
type Field interface {
	Elt(interface{}) Elt  // Constructor of Elements
	Rand(r io.Reader) Elt // Constructor of Elements at Random
	Zero() Elt            // Returns the one element
	One() Elt             // Returns the one element
	Ext() uint            // Extension degree
	BitLen() int          // Bit length of modulus
	P() *big.Int
	hasArith
	hasPredicates
	// hasCmov
	// hasSgn0
	// hasInv0
	// hasSqrt
}

type hasPredicates interface {
	AreEqual(x, y Elt) bool
	IsZero(Elt) bool
}

type hasArith interface {
	Neg(x Elt) Elt
	Add(x, y Elt) Elt
	Sub(x, y Elt) Elt
	Mul(x, y Elt) Elt
	Sqr(x Elt) Elt
	Inv(x Elt) Elt
}

type hasCmov interface{ CMov(x, y Elt, b bool) Elt }
type hasSgn0 interface{ Sgn0(x Elt) int }
type hasInv0 interface{ Inv0(x Elt) Elt }
type hasSqrt interface {
	Sqrt(x Elt) Elt
	IsSquare(x Elt) bool
}

// Elt is
type Elt interface {
	Copy() Elt
}

// NewFromID is
func NewFromID(id Prime) Field { return getFromID(id) }

// NewGF is
func NewGF(p interface{}, m uint, name string) Field {
	if !(m == 1 || m == 2) {
		panic("not implemented")
	}
	modulus := modulus{
		name: name,
		p:    fromType(p),
	}
	if !modulus.p.ProbablyPrime(5) {
		panic(fmt.Errorf("p= %v is not prime", p))
	}
	switch m {
	case 1:
		return newFp(modulus)
	case 2:
		return newFp2(modulus)
	}
	return nil
}
