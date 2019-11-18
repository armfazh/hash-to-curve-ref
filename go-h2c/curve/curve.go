package curve

import (
	"fmt"
	"math/big"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// Point is an elliptic curve point
type Point interface {
	IsIdentity() bool
	IsEqual(Point) bool
	Copy() Point
}

// EllCurve is an elliptic curve
type EllCurve interface {
	NewPoint(x, y GF.Elt) Point
	hasArith
	hasParams
}

type hasArith interface {
	Identity() Point
	IsOnCurve(Point) bool
	Neg(Point) Point
	Double(Point) Point
	Add(Point, Point) Point
}

type Model int

const (
	Weierstrass Model = iota
	Montgomery
	Edwards
)

// Params is
type Params struct {
	F       GF.Field
	A, B, D GF.Elt
	H       *big.Int
	R       *big.Int
}

func (e *Params) String() string {
	return fmt.Sprintf("F: %v\nA: %v\nB: %v\n", e.F, e.A, e.B)
}
func (e *Params) Order() *big.Int    { return e.R }
func (e *Params) Cofactor() *big.Int { return e.H }
func (e *Params) Field() GF.Field    { return e.F }

type hasParams interface {
	Field() GF.Field
	Order() *big.Int
	Cofactor() *big.Int
}
