package bf_test

import (
	"testing"

	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/bf"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

var curves = []string{"W2"}

func TestMap(t *testing.T) {
	for _, id := range curves {
		E := toy.ToyCurves[id].E
		F := E.Field()
		n := F.Order().Int64()
		m := bf.New(E)
		for i := int64(0); i < n; i++ {
			u := F.Elt(i)
			P := m.MapToCurve(u)
			if !E.IsOnCurve(P) {
				t.Fatalf("u: %v got P: %v\n", u, P)
			}
		}
	}
}
