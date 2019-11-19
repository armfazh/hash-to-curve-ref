package sswuAB0

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
)

type sswuAB0 struct {
	E   C.WECurve
	iso C.Isogeny
	mm  M.MapToCurve
}

func (m sswuAB0) String() string { return fmt.Sprintf("Simple SWU AB==0 for E: %v", m.E) }

// New is
func New(e C.EllCurve, z GF.Elt, sgn0 string, iso C.Isogeny) M.MapToCurve {
	if s := (&sswuAB0{
		E:   e.(C.WECurve),
		iso: iso,
	}); s.verify(z, sgn0) {
		return s
	}
	panic(fmt.Errorf("Failed restrictions for sswuAB0"))
}

func (m *sswuAB0) verify(z GF.Elt, sgn0 string) bool {
	cond1 := m.E == m.iso.Codomain()
	m.mm = sswu.New(m.iso.Domain(), z, sgn0)
	cond2 := m.mm != nil
	return cond1 && cond2
}

func (m *sswuAB0) Map(u GF.Elt) C.Point { return m.iso.Apply(m.mm.Map(u)) }
