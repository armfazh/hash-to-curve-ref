package curve

import (
	"crypto/sha256"
	"testing"

	"github.com/armfazh/hash-to-curve-ref/h2c"
	"github.com/armfazh/hash-to-curve-ref/h2c/field"
)

func TestElligator2(t *testing.T) {
	var e Elligator2
	e.A = math.Elt{}
	e.B = math.Elt{}
	e.Z = math.Elt{}
	F := field.GF("103", 1)

	H := sha256.New
	msg := []byte("Lorem ipsum dolor sit amet")
	DST := []byte("QUUX-V01-CS01")
	ctr := uint(0)
	u := h2c.HashToField(H, msg, DST, ctr, 128, F)

	P := e.MapToCurve(u)
	t.Logf("Point: %v", P)
}
