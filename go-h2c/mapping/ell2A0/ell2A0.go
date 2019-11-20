package ell2A0

import (
	"fmt"
	"math/big"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

type ell2A0 struct {
	E      C.M
	Sgn0   func(GF.Elt) int
	AA, BB GF.Elt
	cte    struct {
		C1 *big.Int
	}
}

func (m ell2A0) String() string { return fmt.Sprintf("Elligator2A0 for E: %v", m.E) }

// New is
func New(e C.EllCurve, sgn0 GF.Sgn0ID) M.Map {
	if s := (&ell2A0{E: e.(C.M)}); s.verify() {
		s.precmp(sgn0)
		return s
	}
	panic(fmt.Errorf("Failed restrictions for ell2A0"))
}

func (m *ell2A0) verify() bool {
	F := m.E.F
	q := F.Order()
	precond1 := q.Mod(q, big.NewInt(4)).Int64() == int64(3)
	precond2 := F.IsZero(m.E.A)  // A == 0
	precond3 := !F.IsZero(m.E.B) // B != 0

	return precond1 && precond2 && precond3
}

func (m *ell2A0) precmp(sgn0 GF.Sgn0ID) {
	q := m.E.F.Order()
	m.Sgn0 = m.E.F.GetSgn0(sgn0)

	t0 := big.NewInt(1)
	t0.Add(q, t0) // (q+1)
	t0.Rsh(q, 2)  // (q+1)/4
	m.cte.C1 = t0

	F := m.E.F
	t1 := F.Inv(m.E.B)      // 1/B
	m.AA = F.Mul(t1, m.E.A) // A/B
	m.BB = F.Sqr(t1)        // 1/B^2
}

func (m *ell2A0) MapToCurve(u GF.Elt) C.Point {
	F := m.E.F
	var x1, x2, gx1, x, y GF.Elt
	var e1, e2 bool

	x1 = u                         // 1.  x1 = u
	x2 = F.Neg(x1)                 // 2.  x2 = -x1
	gx1 = F.Sqr(x1)                // 3. gx1 = x1^2
	gx1 = F.Add(gx1, m.BB)         // 4. gx1 = gx1 + B
	gx1 = F.Mul(gx1, x1)           // 5. gx1 = gx1 * x1           // gx1 = x1^3 + B * x1
	y = F.Sqrt(gx1)                // 6.   y = gx1^c1             // This is either sqrt(gx1) or sqrt(gx2)
	e1 = F.AreEqual(F.Sqr(y), gx1) // 7.  e1 = (y^2) == gx1
	x = F.CMov(x2, x1, e1)         // 8.   x = CMOV(x2, x1, e1)
	e2 = m.Sgn0(u) == m.Sgn0(y)    // 9.  e2 = sgn0(u) == sgn0(y)
	y = F.CMov(F.Neg(y), y, e2)    // 10.  y = CMOV(-y, y, e2)
	x = F.Mul(x, m.E.B)
	y = F.Mul(y, m.E.B)
	return m.E.NewPoint(x, y)
}
