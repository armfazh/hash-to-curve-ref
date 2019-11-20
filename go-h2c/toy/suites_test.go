package toy_test

import (
	"strings"
	"testing"

	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

func TestSuites(t *testing.T) {
	msg := []byte("hello world")
	dst := []byte("QUUX-V01-CS01")

	for suiteID, h2c := range toy.ToySuites {
		v := strings.Split(suiteID, "-")
		E := toy.ToyCurves[v[0]].E
		P := h2c.Hash(msg, dst)
		if !E.IsOnCurve(P) {
			t.Fatalf("point not on curve")
		}
	}
}
