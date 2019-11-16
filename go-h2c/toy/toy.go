package toy

import (
	"math/big"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

var f53 = GF.NewFp("p53", 53)

// ToyCurves is
var ToyCurves = []C.EllCurve{
	C.NewWeierstrass(&C.Params{F: f53, A: f53.Elt(3), B: f53.Elt(2), R: big.NewInt(51), H: big.NewInt(3)}),
	C.NewMontgomery(&C.Params{F: f53, A: f53.Elt(4), B: f53.Elt(3), R: big.NewInt(44), H: big.NewInt(4)}),
	C.NewEdwards(&C.Params{F: f53, A: f53.Elt(1), D: f53.Elt(3), R: big.NewInt(44), H: big.NewInt(4)}),
}

// ToyPoints is
var ToyPoints = []C.Point{
	ToyCurves[0].NewPoint(f53.Elt(46), f53.Elt(3)),
	ToyCurves[1].NewPoint(f53.Elt(16), f53.Elt(4)),
	ToyCurves[2].NewPoint(f53.Elt(17), f53.Elt(49)),
}
