package elligator2

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

type mtEll2 struct {
	E C.M
	C.RationalMap
	M.Map
}

func (m mtEll2) String() string { return fmt.Sprintf("Montgomery Elligator2 for E: %v", m.E) }

func newMTEll2(e C.M, sgn0 GF.Sgn0ID, rat C.RationalMap) M.Map {
	if rat == nil {
		rat = e.ToWeierstrassC()
	}
	cDest := rat.Codomain()
	eDest := cDest.(C.WC)
	return &mtEll2{e, rat, newWCEll2(eDest, sgn0)}
}

func (m *mtEll2) MapToCurve(u GF.Elt) C.Point { return m.Pull(m.Map.MapToCurve(u)) }
