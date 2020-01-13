package elligator2

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

type teEll2 struct {
	E C.T
	C.RationalMap
	M.Map
}

func (m teEll2) String() string { return fmt.Sprintf("Edwards Elligator2 for E: %v", m.E) }

func newTEEll2(e C.T, sgn0 GF.Sgn0ID, rat C.RationalMap) M.Map {
	if rat == nil {
		rat = e.ToWeierstrassC()
	}
	cDest := rat.Codomain()
	eDest := cDest.(C.WC)
	return &teEll2{e, rat, newWCEll2(eDest, sgn0)}
}

func (m *teEll2) MapToCurve(u GF.Elt) C.Point { return m.Pull(m.Map.MapToCurve(u)) }
