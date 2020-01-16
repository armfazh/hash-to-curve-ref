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
		rat = getMontToEdw25519()
		ell2Map = newMTEll2(rat.Domain().(C.M), sgn0)
	case C.Edwards448:
		rat = getMontToEdw448()
		ell2Map = newMTEll2(rat.Domain().(C.M), sgn0)
	default:
		rat = e.ToWeierstrassC()
		ell2Map = newWCEll2(rat.Domain().(C.WC), sgn0)
	}
	return &teEll2{e, rat, ell2Map}
}

func (m *teEll2) MapToCurve(u GF.Elt) C.Point { return m.Push(m.Map.MapToCurve(u)) }

type mont2edw25519 struct {
	E0       C.M
	E1       C.T
	invSqrtD GF.Elt
}

// getMontToEdw25519 returns a birational map
func getMontToEdw25519() C.RationalMap {
	e0 := C.Curve25519.Get()
	e1 := C.Edwards25519.Get()
	F := e0.Field()
	return mont2edw25519{
		E0:       e0.(C.M),
		E1:       e1.(C.T),
		invSqrtD: F.Elt("6853475219497561581579357271197624642482790079785650197046958215289687604742"),
	}
}
func (m mont2edw25519) String() string       { return fmt.Sprintf("Rational Map from %v to\n%v", m.E0, m.E1) }
func (m mont2edw25519) Domain() C.EllCurve   { return m.E0 }
func (m mont2edw25519) Codomain() C.EllCurve { return m.E1 }
func (m mont2edw25519) Push(p C.Point) C.Point {
	if p.IsIdentity() {
		return m.E1.Identity()
	}
	F := m.E0.Field()
	t0 := F.Add(F.One(), p.Y()) // 1+y
	t1 := F.Sub(F.One(), p.Y()) // 1-y
	t1 = F.Mul(t1, m.invSqrtD)  // invSqrtD*(1-y)
	t1 = F.Inv(t1)              // 1/(invSqrtD*(1-y))
	x := F.Mul(t0, t1)          // x = (1+y)/(invSqrtD*(1-y))
	t0 = F.Inv(p.Y())           // 1/y
	y := F.Mul(x, t0)           // y = x/y
	return m.E1.NewPoint(x, y)
}
func (m mont2edw25519) Pull(p C.Point) C.Point {
	if p.IsIdentity() {
		return m.E0.Identity()
	}
	F := m.E0.Field()
	if p.IsTwoTorsion() {
		return m.E0.NewPoint(F.Zero(), F.Elt(-1))
	}
	t0 := F.Inv(p.Y())            // 1/y
	x := F.Mul(p.X(), t0)         // X = x/y
	t0 = F.Mul(m.invSqrtD, p.X()) // invSqrtD*x
	t1 := F.Add(t0, F.One())      // invSqrtD*x+1
	t2 := F.Sub(t0, F.One())      // invSqrtD*x-1
	t1 = F.Inv(t1)                // 1/(invSqrtD*x+1)
	y := F.Mul(t1, t2)            // Y = (invSqrtD*x-1)/(invSqrtD*x+1)
	return m.E0.NewPoint(x, y)
}

type mont2edw448 struct {
	E0 C.M
	E1 C.T
}

func getMontToEdw448() C.RationalMap { return nil }
