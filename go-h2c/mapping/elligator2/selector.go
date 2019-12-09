package elligator2

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

// New is
func New(e C.EllCurve, z GF.Elt, sgn0 GF.Sgn0ID, rat C.RationalMap) M.Map {
	switch curve := e.(type) {
	case C.WC:
		if s := (&wcEll2{E: curve, Z: z}); s.verify() {
			s.precmp(sgn0)
			return s
		}
	case C.T:
		if s := (&teEll2{E: curve, rat: rat, wcEll2: wcEll2{Z: z}}); s.verify() {
			s.precmp(sgn0)
			return s
		}
	case C.M:
		if s := (&mtEll2{E: curve, rat: rat, wcEll2: wcEll2{Z: z}}); s.verify() {
			s.precmp(sgn0)
			return s
		}
	default:
		panic(fmt.Errorf("Curve didn't match elligator2 mapping"))
	}
	panic(fmt.Errorf("Failed restrictions for ell2"))
}
