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

func newTEEll2(e C.T, sgn0 GF.Sgn0ID) M.Map {
	var rat C.RationalMap
	var ell2Map M.Map
	switch e.Id {
	case C.Edwards25519:
		rat = C.FromTe2Mt25519()
		ell2Map = newMTEll2(rat.Codomain().(C.M), sgn0)
	case C.Edwards448:
		rat = C.FromTe2Mt4ISO448()
		ell2Map = newMTEll2(rat.Codomain().(C.M), sgn0)
	default:
		rat = e.ToWeierstrassC()
		ell2Map = newWCEll2(rat.Codomain().(C.WC), sgn0)
	}
	return &teEll2{e, rat, ell2Map}
}

func (m *teEll2) MapToCurve(u GF.Elt) C.Point { return m.Pull(m.Map.MapToCurve(u)) }
