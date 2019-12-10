package sswuAB0_test

import (
	"testing"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

var curves = []string{"W1iso"}

func TestMap(t *testing.T) {
	for _, id := range curves {
		E := toy.ToyCurves[id].E
		isogeny := identityMap{E}
		F := E.Field()
		n := F.Order().Int64()
		Z := F.Elt(3)
		for _, m := range []mapping.Map{
			sswuAB0.New(E, Z, GF.SignLE, isogeny),
			sswuAB0.New(E, Z, GF.SignBE, isogeny),
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

type identityMap struct{ E C.EllCurve }

func (t identityMap) Domain() C.EllCurve     { return t.E }
func (t identityMap) Codomain() C.EllCurve   { return t.E }
func (t identityMap) Push(p C.Point) C.Point { return p }
func (t identityMap) Pull(p C.Point) C.Point { return p }
