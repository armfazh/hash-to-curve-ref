package sswuAB0_test

import (
	"testing"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

type trivialMap struct{ E C.EllCurve }

func (t trivialMap) Domain() C.EllCurve     { return t.E }
func (t trivialMap) Codomain() C.EllCurve   { return t.E }
func (t trivialMap) Push(p C.Point) C.Point { return p }
func (t trivialMap) Pull(p C.Point) C.Point { return p }

func TestMap(t *testing.T) {
	E := toy.ToyCurves["W1iso"].E
	isogeny := trivialMap{E}
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
