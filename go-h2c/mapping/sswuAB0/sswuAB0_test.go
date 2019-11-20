package sswuAB0_test

import (
	"testing"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

func TestMap(t *testing.T) {
	trivialMap := func(e0, e1 C.EllCurve, p C.Point) C.Point { return p }
	isogeny := &C.Isogeny{
		E0:  toy.ToyCurves["W1iso"].E,
		E1:  toy.ToyCurves["W1iso"].E,
		Map: trivialMap}
	E := isogeny.Codomain()
	F := E.Field()
	n := F.Order().Int64()
	Z := F.Elt(3)
	for _, m := range []mapping.MapToCurve{
		sswuAB0.New(E, Z, "be", isogeny),
		sswuAB0.New(E, Z, "le", isogeny),
	} {
		for i := int64(0); i < n; i++ {
			u := F.Elt(i)
			P := m.Map(u)
			if !E.IsOnCurve(P) {
				t.Fatalf("u: %v got P: %v\n", u, P)
			}
		}
	}
}
