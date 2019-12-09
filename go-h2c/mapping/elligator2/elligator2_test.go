package elligator2_test

import (
	"testing"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/elligator2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

var curves = []struct {
	Name string
	Z    uint
}{
	{"M0", 3},
	{"M1", 5},
	{"E0", 2},
}

func TestMap(t *testing.T) {
	for _, id := range curves {
		E := toy.ToyCurves[id.Name].E
		F := E.Field()
		n := F.Order().Int64()
		Z := F.Elt(id.Z)
		for _, m := range []mapping.Map{
			elligator2.New(E, Z, GF.SignLE, nil),
			elligator2.New(E, Z, GF.SignBE, nil),
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
