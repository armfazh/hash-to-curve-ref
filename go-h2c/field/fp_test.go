package field_test

import (
	"testing"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

func TestSqrt(t *testing.T) {
	var primes = []int{
		607, // 3 mod 4
		613, // 5 mod 8
		// 617, // 9 mod 16
		// 641, // 1 mod 16
	}
	for _, p := range primes {
		testSqrt(t, p)
	}
}

func testSqrt(t *testing.T, p int) {
	F := GF.NewFp(GF.PrimeID(p), p)
	for i := 0; i < p; i++ {
		x := F.Elt(i)
		if F.IsSquare(x) {
			y := F.Sqrt(x)
			got := F.Sqr(y)
			want := x
			if !F.AreEqual(got, want) {
				t.Fatalf("got: %v\nwant: %v\nF:%v", got, want, F)
			}
		}
	}
}
