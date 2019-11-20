package toy

import (
	"math/big"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// EC is an elliptic curve
type EC struct {
	E    C.EllCurve
	P    C.Point
	Z    GF.Elt
	H    *big.Int
	Iso  *C.Isogeny
	Sgn0 GF.Sgn0ID
}

// ToyCurves is
var ToyCurves map[string]EC

var WeCurves, TeCurves, MtCurves []EC

func initCurves() {
	ToyCurves = make(map[string]EC)

	var f53 = GF.NewFp("p53", 53)
	RegisterToyCurve("W0",
		C.NewWeierstrass(f53, f53.Elt(3), f53.Elt(2), big.NewInt(51), big.NewInt(3)),
		f53.Elt(46), f53.Elt(3), f53.Elt(3), GF.SignLE)

	RegisterToyCurve("W1",
		C.NewWeierstrass(f53, f53.Zero(), f53.One(), big.NewInt(54), big.NewInt(2)),
		f53.Elt(13), f53.Elt(5), f53.Elt(3), GF.SignBE)

	RegisterToyCurve("W1iso",
		C.NewWeierstrass(f53, f53.Elt(38), f53.Elt(22), big.NewInt(54), big.NewInt(2)),
		f53.Elt(41), f53.Elt(45), f53.Elt(3), GF.SignLE)

	RegisterToyCurve("M0",
		C.NewMontgomery(f53, f53.Elt(4), f53.Elt(3), big.NewInt(44), big.NewInt(4)),
		f53.Elt(16), f53.Elt(4), f53.Elt(3), GF.SignLE)

	RegisterToyCurve("E0",
		C.NewEdwards(f53, f53.Elt(1), f53.Elt(3), big.NewInt(44), big.NewInt(4)),
		f53.Elt(17), f53.Elt(49), f53.Elt(3), GF.SignLE)

}

// RegisterToyCurve is
func RegisterToyCurve(name string, e C.EllCurve, x, y, z GF.Elt, sgn0 GF.Sgn0ID) {
	var list *[]EC
	switch e.(type) {
	case *C.WECurve:
		list = &WeCurves
	case *C.TECurve:
		list = &TeCurves
	case *C.MTCurve:
		list = &MtCurves
	default:
		panic("model not supported")
	}
	n := EC{E: e, P: e.NewPoint(x, y), Z: z, Sgn0: sgn0}
	ToyCurves[name] = n
	*list = append(*list, n)
}
