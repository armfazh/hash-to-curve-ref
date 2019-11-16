package sswu_test

import (
	"crypto/rand"
	"testing"

	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

func TestMap(t *testing.T) {
	E := toy.ToyCurves[0]
	F := E.Field()
	Z := F.Elt(3)
	mapping := sswu.New(E, Z, "le")
	t.Logf("E: %v\n", E)
	t.Logf("Map: %v\n", mapping)
	t.Logf("Z: %v\n", Z)

	u := F.Rand(rand.Reader)
	t.Logf("u: %v\n", u)
	P := mapping.Map(u)
	t.Logf("P: %v\n", P)
	t.Logf("P\\in E: %v\n", E.IsOnCurve(P))
}
