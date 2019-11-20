package toy

import (
	"math/big"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// EC is an elliptic curve
type EC struct {
	E   C.EllCurve
	P   C.Point
	Iso *C.Isogeny
}

// ToyCurves is
var ToyCurves map[string]EC

func initCurves() {
	ToyCurves = make(map[string]EC)

	var f53 = GF.NewFp("p53", 53) // 1mod4, 2mod3
	var f59 = GF.NewFp("p59", 59) // 3mod4, 2mod3

	RegisterToyCurve("W0",
		C.NewWeierstrass(f53, f53.Elt(3), f53.Elt(2), big.NewInt(51), big.NewInt(3)),
		f53.Elt(46), f53.Elt(3))

	RegisterToyCurve("W1",
		C.NewWeierstrass(f53, f53.Zero(), f53.One(), big.NewInt(54), big.NewInt(2)),
		f53.Elt(13), f53.Elt(5))

	RegisterToyCurve("W1iso",
		C.NewWeierstrass(f53, f53.Elt(38), f53.Elt(22), big.NewInt(54), big.NewInt(2)),
		f53.Elt(41), f53.Elt(45))

	RegisterToyCurve("W2",
		C.NewWeierstrass(f53, f53.Zero(), f53.Elt(2), big.NewInt(51), big.NewInt(3)),
		f53.Elt(37), f53.Elt(27))

	RegisterToyCurve("M0",
		C.NewMontgomery(f53, f53.Elt(4), f53.Elt(3), big.NewInt(44), big.NewInt(4)),
		f53.Elt(16), f53.Elt(4))

	RegisterToyCurve("M1",
		C.NewMontgomery(f53, f53.Elt(3), f53.Elt(1), big.NewInt(48), big.NewInt(4)),
		f53.Elt(14), f53.Elt(22))

	RegisterToyCurve("M2",
		C.NewMontgomery(f59, f59.Zero(), f59.Elt(16), big.NewInt(60), big.NewInt(4)),
		f59.Elt(31), f59.Elt(36))

	RegisterToyCurve("E0",
		C.NewEdwards(f53, f53.Elt(1), f53.Elt(3), big.NewInt(44), big.NewInt(4)),
		f53.Elt(17), f53.Elt(49))

}

// RegisterToyCurve is
func RegisterToyCurve(name string, e C.EllCurve, x, y GF.Elt) {
	n := EC{E: e, P: e.NewPoint(x, y)}
	ToyCurves[name] = n
}
