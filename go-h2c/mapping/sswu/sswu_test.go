package sswu_test

import (
	"testing"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

var curves = []struct {
	Name string
	Z    uint
}{
	{"W0", 3},
	{"W1iso", 3},
}

func TestMap(t *testing.T) {
	for _, c := range curves {
		E := toy.ToyCurves[c.Name].E
		F := E.Field()
		n := F.Order().Int64()
		Z := F.Elt(c.Z)
		for _, m := range []mapping.Map{
			sswu.New(E, Z, GF.SignLE),
			sswu.New(E, Z, GF.SignBE),
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
