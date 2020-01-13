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
	hasCMov
	hasSgn0
	hasInv0
	hasSqrt
	hasExp
	hasAdvanced
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

type hasExp interface{ Exp(Elt, *big.Int) Elt }
type hasCMov interface{ CMov(x, y Elt, b bool) Elt }
type hasInv0 interface{ Inv0(Elt) Elt }
type hasSqrt interface{ Sqrt(Elt) Elt }
type hasSgn0 interface{ GetSgn0(Sgn0ID) func(Elt) int }
type hasAdvanced interface{ Generator() Elt }

type Sgn0ID int

const (
	// SignLE is
	SignLE Sgn0ID = iota
	// SignBE is
	SignBE
)
