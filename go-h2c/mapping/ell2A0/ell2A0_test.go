package ell2A0_test

import (
	"testing"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/ell2A0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

var curves = []string{
	"M2",
}

func TestMap(t *testing.T) {
	for _, id := range curves {
		E := toy.ToyCurves[id].E
		F := E.Field()
		n := F.Order().Int64()
		for _, m := range []mapping.Map{
			ell2A0.New(E, GF.SignLE),
			ell2A0.New(E, GF.SignBE),
		} {
			for i := int64(0); i < n; i++ {
				u := F.Elt(i)
				P := m.MapToCurve(u)
				if !E.IsOnCurve(P) {
					t.Fatalf("u: %v got P: %v\n", u, P)
				}
			}
		}
	}
}
