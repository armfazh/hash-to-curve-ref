// Package mapping contains a set of deterministic functions that take an
// element of the field F and return a point on an elliptic curve E over F.
// Certain mappings restrict the form of the curve or its parameters.
//
// Choosing a mapping function
//
// If the target elliptic curve is
//  - a supersingular curve, then use either the Boneh-Franklin method (package bf) or the Elligator 2 method for A == 0 (package elligator2);
//  - a Montgomery curve, then use the Elligator 2 (package elligator2);
//  - a twisted Edwards curve, then use Elligator 2 (package elligator2);
//  - a Weierstrass curve, then use either the Simplified SWU (package sswu). But if either A or B is zero, then use the special case of Simplified SWU (package sswuAB0);
//  - if none of the above applies, then use the Shallue-van de Woestijne method (package svdw).
//
// Map is the generic interface shared by all mappings. To instantiate a mapping
// use the New function provided in each package.
// Note: the mappings must not be used standalone, since its correct usage is
// determined by EncodeToCurve or HashToCurve high-level interfaces.
package mapping

import (
	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// MapToCurve maps a field element into a elliptic curve point.
type MapToCurve interface {
	Map(GF.Elt) C.Point
}

// ID is an identifier of a mapping.
type ID uint

const (
	// SSWU is the Simplified SWU method.
	SSWU ID = iota
	// SVDW is Shallue-Woestijne method.
	SVDW
	// ELL2 is Elligator2 method for Montgomery curves.
	ELL2
	// EDELL2 is Elligator2 method for twisted Edwards curves.
	EDELL2
	// BF is the Boneh-Franklin method
	BF
)

// Get returns a MapToCurve instance based on MappingID.
func (id ID) Get(e C.EllCurve, z GF.Elt, sgn0 GF.Sgn0ID, iso C.Isogeny) MapToCurve {
	switch id {
	case BF:
		return NewBF(e)
	case SSWU:
		return NewSSWU(e, z, sgn0, iso)
	case SVDW:
		return NewSVDW(e, sgn0)
	case ELL2, EDELL2:
		return NewElligator2(e, sgn0)
	default:
		panic("Mapping not supported")
	}
}
