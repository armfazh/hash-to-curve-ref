package mapping

import (
	"github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// Map is
type Map interface {
	MapToCurve(field.Elt) curve.Point
}
