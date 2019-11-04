package math

import (
	"math/big"
)

// FiniteField is
type FiniteField interface {
	secureMove
	squareRoot
	Arith
	Char() *big.Int // Field characteristic
	IsPrime() bool  // True if field is a prime field
	Elt() Element   // Constructor of elements
	Zero() Element
	One() Element
}

type secureMove interface {
	CMov(x, y Element, b bool) Element
}

type squareRoot interface {
	Sqrt(x Element) Element
	IsSquare(x Element) bool
}

// Arith is
type Arith interface {
	Neg(x Element) Element
	Add(x, y Element) Element
	Sub(x, y Element) Element
	Mul(x, y Element) Element
	Sqr(x Element) Element
	Inv0(x Element) Element
	Sgn0(x Element) int
}

// Element is
type Element interface {
	IsZero() bool
}

// NEWGF is
func NEWGF(p string, m uint) FiniteField {
	if m == 0 || m > 2 {
		panic("not implemented")
	}
	prime, _ := new(big.Int).SetString(p, 0)
	if !prime.ProbablyPrime(5) {
		panic("p is not a prime number")
	}
	modulus := ff{p: prime}
	switch m {
	case 1:
		return newFp(modulus)
	case 2:
		return fquadratic{modulus}
	}
	return nil
}

type ff struct {
	p *big.Int
}

func (f ff) String() string { return "0x" + f.p.Text(16) }
func (f ff) Char() *big.Int { return new(big.Int).Set(f.p) }
