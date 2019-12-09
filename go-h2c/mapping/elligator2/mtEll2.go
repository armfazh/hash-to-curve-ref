package elligator2

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type mtEll2 struct {
	E   C.M
	rat C.RationalMap
	wcEll2
}

func (m mtEll2) String() string { return fmt.Sprintf("Montgomery Elligator2 for E: %v", m.E) }

func (m *mtEll2) verify() bool {
	if m.rat == nil {
		m.rat = m.E.ToWeierstrassC()
		cDest := m.rat.Codomain()
		m.wcEll2.E = cDest.(C.WC)
	}
	return m.wcEll2.verify()
}

func (m *mtEll2) precmp(sgn0 GF.Sgn0ID)       { m.wcEll2.precmp(sgn0) }
func (m *mtEll2) MapToCurve(u GF.Elt) C.Point { return m.rat.Pull(m.wcEll2.MapToCurve(u)) }
