package elligator2

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

// New is
func New(e C.EllCurve, sgn0 GF.Sgn0ID) M.Map {
	switch curve := e.(type) {
	case C.W:
		return newWA0Ell2(curve, sgn0)
	case C.WC:
		return newWCEll2(curve, sgn0)
	case C.M:
		return newMTEll2(curve, sgn0)
	case C.T:
		return newTEEll2(curve, sgn0)
	default:
		panic(fmt.Errorf("Curve doesn't support an elligator2 mapping"))
	}
}
