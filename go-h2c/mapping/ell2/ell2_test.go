package ell2_test

import (
	"testing"

	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/ell2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

func TestMap(t *testing.T) {
	for _, EC := range toy.MtCurves {
		E := EC.E
		F := E.Field()
		n := F.Order().Int64()
		Z := F.Elt(3)
		for _, m := range []mapping.MapToCurve{
			ell2.New(E, Z, "be"),
			ell2.New(E, Z, "le"),
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
}
