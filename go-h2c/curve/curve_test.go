package curve_test

import (
	"testing"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

func TestCurves(t *testing.T) {
	for i := range toy.ToyCurves {
		t.Run("", func(t *testing.T) {
			testAdd(t, toy.ToyCurves[i], toy.ToyPoints[i], toy.ToyCurves[i].Order().Uint64())
		})
	}
}

func testAdd(t *testing.T, ec C.EllCurve, g C.Point, order uint64) {
	T := make([]C.Point, order)
	T[0] = ec.Identity()
	for i := 1; i < len(T); i++ {
		T[i] = ec.Add(T[i-1], g)
		if !ec.IsOnCurve(T[i]) {
			t.Fatalf("point not in the curve: %v\n", T[i])
		}
	}
	for _, P := range T {
		for _, Q := range T {
			R := ec.Add(P, Q)
			if !ec.IsOnCurve(R) {
				t.Fatalf("point not in the curve: %v\n", R)
			}
		}
	}
}

func BenchmarkCurve(b *testing.B) {
	ec := toy.ToyCurves[0]
	P := toy.ToyPoints[0]
	Q := ec.Double(P)
	b.Run("double", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P = ec.Double(P)
		}
	})

	b.Run("add", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			P = ec.Add(P, Q)
		}
	})
}
