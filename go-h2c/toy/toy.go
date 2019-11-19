package toy

import (
	"math/big"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// EC is an elliptic curve
type EC struct {
	E C.EllCurve
	P C.Point
}

// ToyCurves is
var ToyCurves map[string]EC

var WeCurves, EdCurves, MtCurves []EC

func init() {
	ToyCurves = make(map[string]EC)

	var f53 = GF.NewFp("p53", 53)
	RegisterToyCurve("W0", C.Weierstrass, &C.Params{
		F: f53, A: f53.Elt(3), B: f53.Elt(2),
		R: big.NewInt(51), H: big.NewInt(3)},
		f53.Elt(46), f53.Elt(3))

	RegisterToyCurve("W1", C.Weierstrass, &C.Params{
		F: f53, A: f53.Zero(), B: f53.One(),
		R: big.NewInt(54), H: big.NewInt(2)},
		f53.Elt(13), f53.Elt(5))

	RegisterToyCurve("W1iso", C.Weierstrass, &C.Params{
		F: f53, A: f53.Elt(38), B: f53.Elt(22),
		R: big.NewInt(54), H: big.NewInt(2)},
		f53.Elt(41), f53.Elt(45))

	RegisterToyCurve("M0", C.Montgomery, &C.Params{
		F: f53, A: f53.Elt(4), B: f53.Elt(3),
		R: big.NewInt(44), H: big.NewInt(4)},
		f53.Elt(16), f53.Elt(4))

	RegisterToyCurve("E0", C.Edwards, &C.Params{
		F: f53, A: f53.Elt(1), D: f53.Elt(3),
		R: big.NewInt(44), H: big.NewInt(4)},
		f53.Elt(17), f53.Elt(49))

}

// RegisterToyCurve is
func RegisterToyCurve(name string, model C.Model, params *C.Params, x, y GF.Elt) {
	var E C.EllCurve
	var list *[]EC
	switch model {
	case C.Weierstrass:
		E = C.NewWeierstrass(params)
		list = &WeCurves
	case C.Edwards:
		E = C.NewEdwards(params)
		list = &EdCurves
	case C.Montgomery:
		E = C.NewMontgomery(params)
		list = &MtCurves
	default:
		panic("model not supported")
	}
	P := E.NewPoint(x, y)
	n := EC{E, P}
	ToyCurves[name] = n
	*list = append(*list, n)
}
