package elligator2_test

import (
	"testing"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/elligator2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

var curves = []string{
	"M0",
	"M1",
	"E0",
	"W3",
}

func TestMap(t *testing.T) {
	for _, id := range curves {
		E := toy.ToyCurves[id].E
		F := E.Field()
		n := F.Order().Int64()
		for _, m := range []mapping.Map{
			elligator2.New(E, GF.SignLE, nil),
			elligator2.New(E, GF.SignBE, nil),
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
