// Package curve provides definitions of several models of elliptic curves
// defined over finite fields of large prime characteristic.
package curve

import (
	"math/big"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// Point represents an elliptic curve point.
type Point interface {
	Copy() Point
	IsIdentity() bool
	IsEqual(Point) bool
	IsTwoTorsion() bool
	X() GF.Elt
	Y() GF.Elt
}

// EllCurve represents an elliptic curve group.
type EllCurve interface {
	Field() GF.Field
	Order() *big.Int
	Cofactor() *big.Int
	NewPoint(x, y GF.Elt) Point
	// Predicates
	IsOnCurve(Point) bool
	IsEqual(EllCurve) bool
	// Arithmetic operations
	Identity() Point
	Neg(Point) Point
	Add(Point, Point) Point
	Double(Point) Point
	ClearCofactor(Point) Point
}

// RationalMap represents a birational map between two elliptic curves.
type RationalMap interface {
	Domain() EllCurve
	Codomain() EllCurve
	Push(Point) Point
	Pull(Point) Point
}

// Isogeny represents an isogeny between two elliptic curves.
type Isogeny interface {
	Domain() EllCurve
	Codomain() EllCurve
	Push(Point) Point
}
