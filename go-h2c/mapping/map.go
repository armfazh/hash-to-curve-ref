package mapping

import (
	"github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// MapToCurve is
type MapToCurve interface {
	Map(field.Elt) curve.Point
}
