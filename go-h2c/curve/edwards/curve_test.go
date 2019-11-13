package edwards_test

import (
	"testing"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/curve/edwards"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

func TestOne(t *testing.T) {
	t.Logf("Montgomery curves")

	F0 := field.NewGF("103", 2, "p103")
	a0, b0 := F0.Zero(), F0.Elt([]interface{}{-3, uint64(0xff)})
	E0 := edwards.NewCurve(F0, a0, b0, 8)
	t.Logf("E: %v\n", E0)

	F1 := field.NewFromID(field.P25519)
	a1, b1 := F1.Elt(-1), F1.Elt(-2)
	E1 := edwards.NewCurve(F1, a1, b1, 8)
	t.Logf("E: %v\n", E1)

	P := E1.NewPoint(F1.Elt(9), F1.Elt(3))
	Q := E1.NewPoint(F1.Elt(79), F1.Elt(3))
	R := E1.Add(P, Q)
	t.Logf("P: %v\n", P)
	t.Logf("Q: %v\n", Q)
	t.Logf("R: %v\n", R)
}

func TestAdd(t *testing.T) {
	F := field.NewGF("53", 1, "p53")
	E := edwards.NewCurve(F, F.Elt(4), F.Elt(3), 4)
	G := E.NewPoint(F.Elt(16), F.Elt(4))
	order := 44
	T := make([]C.Point, order)
	T[0] = E.Identity()
	for i := 1; i < len(T); i++ {
		T[i] = E.Add(T[i-1], G)
		if !E.IsOnCurve(T[i]) {
			t.Errorf("point not in the curve")
		}
	}
	for _, P := range T {
		for _, Q := range T {
			R := E.Add(P, Q)
			if !E.IsOnCurve(R) {
				t.Errorf("point not in the curve")
			}
		}
	}
}
