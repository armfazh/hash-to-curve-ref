package sswuAB0

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
)

type sswuAB0 struct {
	E   C.W
	iso C.Isogeny
	M.Map
}

func (m sswuAB0) String() string { return fmt.Sprintf("Simple SWU AB==0 for E: %v", m.E) }

// New is
func New(e C.EllCurve, z GF.Elt, sgn0 GF.Sgn0ID, iso C.Isogeny) M.Map {
	E := e.(C.W)
	F := E.F
	cond1 := !F.IsZero(E.A)
	cond2 := !F.IsZero(E.B)
	if cond1 && cond2 {
		return &sswuAB0{e.(C.W), iso, sswu.New(iso.Domain(), z, sgn0)}
	}
	panic(fmt.Errorf("Failed restrictions for sswuAB0"))
}

func (m *sswuAB0) MapToCurve(u GF.Elt) C.Point { return m.iso.Push(m.Map.MapToCurve(u)) }
