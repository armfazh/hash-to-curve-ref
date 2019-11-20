package ell2_test

import (
	"testing"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/ell2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

func TestMap(t *testing.T) {
	for _, EC := range toy.MtCurves {
		F := EC.E.Field()
		n := F.Order().Int64()
		for _, m := range []mapping.Map{
			ell2.New(EC.E, EC.Z, GF.SignLE),
			ell2.New(EC.E, EC.Z, GF.SignBE),
		} {
			for i := int64(0); i < n; i++ {
				u := F.Elt(i)
				P := m.MapToCurve(u)
				if !EC.E.IsOnCurve(P) {
					t.Fatalf("u: %v got P: %v\n", u, P)
				}
			}
		}
	}
}
