package field

import (
	"io"
	"math/big"
)

// Elt is a field element
type Elt interface {
	Copy() Elt
}

// Field is
type Field interface {
	Zero() Elt            // Constructor of elements
	One() Elt             // Returns the one element
	Elt(interface{}) Elt  // Constructor of elements from a value
	Rand(r io.Reader) Elt // Constructor of elements at random
	P() *big.Int          // Characteristic of the field
	Order() *big.Int      // Size of the field
	Ext() uint            // Extension degree of field
	BitLen() int          // Bit length of modulus
	hasArith
	hasPredicates
}

type hasPredicates interface {
	AreEqual(Elt, Elt) bool
	IsZero(Elt) bool
	IsSquare(Elt) bool
}

type hasArith interface {
	Neg(x Elt) Elt
	Add(x, y Elt) Elt
	Sub(x, y Elt) Elt
	Mul(x, y Elt) Elt
	Sqr(x Elt) Elt
	Inv(x Elt) Elt
}

// HasCMov is
type HasCMov interface{ CMov(x, y Elt, b bool) Elt }

// HasSgn0 is
type HasSgn0 interface {
	Sgn0BE(Elt) int
	Sgn0LE(Elt) int
}

// HasInv0 is
type HasInv0 interface{ Inv0(Elt) Elt }

// HasSqrt is
type HasSqrt interface{ Sqrt(Elt) Elt }
