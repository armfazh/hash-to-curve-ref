package curve

import GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"

// Point is an elliptic curve point
type Point interface {
	IsIdentity() bool
	IsEqual(Point) bool
	Copy() Point
}

// EllCurve is an elliptic curve
type EllCurve interface {
	NewPoint(x, y GF.Elt) Point
	Identity() Point
	IsOnCurve(Point) bool
	Neg(p Point) Point
	Add(p, q Point) Point
	Double(p Point) Point
}
