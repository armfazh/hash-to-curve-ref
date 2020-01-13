package svdw_test

import (
	"testing"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/svdw"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

var curves = []string{
	"W0",
}

func TestMap(t *testing.T) {
	for _, c := range curves {
		E := toy.ToyCurves[c].E
		F := E.Field()
		n := F.Order().Int64()
		for _, m := range []mapping.Map{
			svdw.New(E, GF.SignLE),
			svdw.New(E, GF.SignBE),
		} {
			for i := int64(0); i < n; i++ {
				u := F.Elt(i)
				P := m.MapToCurve(u)
				if !E.IsOnCurve(P) {
					t.Fatalf("%vu: %v\nP: %v not on curve.", m, u, P)
				}
			}
		}
	}
}
