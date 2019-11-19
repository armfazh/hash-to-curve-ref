package sswuAB0_test

import (
	"testing"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

func TestMap(t *testing.T) {
	trivialIso := func(p C.Point) C.Point { return p }
	iso := C.NewIsogeny(
		toy.ToyCurves["W1iso"].E,
		toy.ToyCurves["W1iso"].E,
		uint(1), trivialIso)
	E := iso.Codomain()
	F := E.Field()
	n := F.Order().Int64()
	Z := F.Elt(3)
	for _, m := range []mapping.MapToCurve{
		sswuAB0.New(E, Z, "be", iso),
		sswuAB0.New(E, Z, "le", iso),
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
