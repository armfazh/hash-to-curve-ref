package curve_test

import (
	"fmt"
	"testing"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	TE "github.com/armfazh/hash-to-curve-ref/go-h2c/curve/edwards"
	MT "github.com/armfazh/hash-to-curve-ref/go-h2c/curve/montgomery"
	WE "github.com/armfazh/hash-to-curve-ref/go-h2c/curve/weierstrass"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type curveParams struct {
	e      C.Model           // model of curve
	coef   map[string]uint64 // curve parameters
	p      uint64            // field characteristic
	m      uint              // degree of extension field
	r      uint              // order of curve
	h      uint              // cofactor
	gx, gy uint64            // point of order r
}

var toyCurves = []curveParams{
	curveParams{
		e: C.ModelWeierstrass, coef: map[string]uint64{"A": 3, "B": 2},
		p: 53, m: 1, r: 51, h: 3, gx: 46, gy: 3,
	},
	curveParams{
		e: C.ModelMontgomery, coef: map[string]uint64{"A": 4, "B": 3},
		p: 53, m: 1, r: 44, h: 4, gx: 16, gy: 4,
	},
	curveParams{
		e: C.ModelEdwards, coef: map[string]uint64{"A": 1, "D": 3},
		p: 53, m: 1, r: 44, h: 4, gx: 17, gy: 49,
	},
}

func TestCurves(t *testing.T) {
	for _, toy := range toyCurves {
		t.Run("", func(t *testing.T) {
			F := field.NewGF(toy.p, toy.m, fmt.Sprintf("p%v", toy.p))
			var E C.EllCurve
			switch toy.e {
			case C.ModelWeierstrass:
				E = WE.NewCurve(F, F.Elt(toy.coef["A"]), F.Elt(toy.coef["B"]), toy.h)
			case C.ModelEdwards:
				E = TE.NewCurve(F, F.Elt(toy.coef["A"]), F.Elt(toy.coef["D"]), toy.h)
			case C.ModelMontgomery:
				E = MT.NewCurve(F, F.Elt(toy.coef["A"]), F.Elt(toy.coef["B"]), toy.h)
			default:
				panic("curve not supported")
			}
			G := E.NewPoint(F.Elt(toy.gx), F.Elt(toy.gy))
			testAdd(t, E, G, toy.r)
		})
	}
}

func testAdd(t *testing.T, ec C.EllCurve, g C.Point, order uint) {
	T := make([]C.Point, order)
	T[0] = ec.Identity()
	for i := 1; i < len(T); i++ {
		T[i] = ec.Add(T[i-1], g)
		if !ec.IsOnCurve(T[i]) {
			t.Fatalf("point not in the curve: %v\n", T[i])
		}
	}
	for _, P := range T {
		for _, Q := range T {
			R := ec.Add(P, Q)
			if !ec.IsOnCurve(R) {
				t.Fatalf("point not in the curve: %v\n", R)
			}
		}
	}
}
