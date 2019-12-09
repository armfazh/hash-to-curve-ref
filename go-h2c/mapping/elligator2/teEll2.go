package elligator2

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type teEll2 struct {
	E   C.T
	rat C.RationalMap
	wcEll2
}

func (m teEll2) String() string { return fmt.Sprintf("Edwards Elligator2 for E: %v", m.E) }

func (m *teEll2) verify() bool {
	if m.rat == nil {
		m.rat = m.E.ToWeierstrassC()
		cDest := m.rat.Codomain()
		m.wcEll2.E = cDest.(C.WC)
	}
	return m.wcEll2.verify()
}

func (m *teEll2) precmp(sgn0 GF.Sgn0ID)       { m.wcEll2.precmp(sgn0) }
func (m *teEll2) MapToCurve(u GF.Elt) C.Point { return m.rat.Pull(m.wcEll2.MapToCurve(u)) }
