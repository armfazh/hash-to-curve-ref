// Package mapping contains a set of deterministic functions that take an
// element of the field F and return a point on an elliptic curve E over F.
// Certain mappings restrict the form of the curve or its parameters.
//
// Choosing a mapping function
//
// If the target elliptic curve is
// - a supersingular curve, then use either the Boneh-Franklin method
//   (package bf) or the Elligator 2 method for A == 0 (package ell2A0);
// - a Montgomery curve, then use the Elligator 2 (package ell2);
// - a twisted Edwards curve, then use Elligator 2 (package edell2);
// - a Weierstrass curve, then use either the Simplified SWU (package sswu).
//   But if either A or B is zero, then use the special case of Simplified SWU
//   (package sswuAB0);
// - if none of the above applies, then use the Shallue-van de Woestijne method
//   (package svdw).
//
// Map is the generic interface shared by all mappings. To instantiate a mapping
// use the New function provided in each package.
// Note: the mappings cannot be used standalone, its usage is determined by
// either EncodeToCurve or HashToCurve interfaces.
package mapping

import (
	"github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// Map is
type Map interface {
	MapToCurve(field.Elt) curve.Point
}
