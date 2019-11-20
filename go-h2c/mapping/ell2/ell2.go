package ell2

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

type ell2 struct {
	E      C.M
	Z      GF.Elt
	AA, BB GF.Elt
	Sgn0   func(GF.Elt) int
}

func (m ell2) String() string { return fmt.Sprintf("Elligator2 for E: %v", m.E) }

// New is
func New(e C.EllCurve, z GF.Elt, sgn0 GF.Sgn0ID) M.Map {
	if s := (&ell2{E: e.(C.M), Z: z}); s.verify() {
		s.precmp(sgn0)
		return s
	}
	panic(fmt.Errorf("Failed restrictions for ell2"))
}

func (m ell2) verify() bool {
	F := m.E.F
	precond1 := !F.IsZero(m.E.A) // A != 0
	precond2 := !F.IsZero(m.E.B) // B != 0
	cond1 := !F.IsSquare(m.Z)    // Z is non-square

	T := m.E.NewPoint(F.Elt(0), F.Elt(0))
	fmt.Println(T)
	Q := m.E.Double(T)
	fmt.Println(Q)

	return precond1 && precond2 && cond1
}

func (m *ell2) precmp(sgn0 GF.Sgn0ID) {
	F := m.E.F

	m.Sgn0 = F.GetSgn0(sgn0)

	t0 := F.Inv(m.E.B)      // 1/B
	m.AA = F.Mul(t0, m.E.A) // A/B
	m.BB = F.Sqr(t0)        // 1/B^2
}

func (m ell2) MapToCurve(u GF.Elt) C.Point {
	F := m.E.F
	var t1 GF.Elt
	var x1, x2, gx1, gx2, y2, x, y GF.Elt
	var e1, e2, e3 bool
	t1 = F.Sqr(u)                  // 1.   t1 = u^2
	t1 = F.Mul(m.Z, t1)            // 2.   t1 = Z * t1              // Z * u^2
	e1 = F.AreEqual(t1, F.Elt(-1)) // 3.   e1 = t1 == -1            // exceptional case: Z * u^2 == -1
	t1 = F.CMov(t1, F.Zero(), e1)  // 4.   t1 = CMOV(t1, 0, e1)     // if t1 == -1, set t1 = 0
	x1 = F.Add(t1, F.One())        // 5.   x1 = t1 + 1
	x1 = F.Inv0(x1)                // 6.   x1 = inv0(x1)
	x1 = F.Mul(F.Neg(m.AA), x1)    // 7.   x1 = -A * x1             // x1 = -A / (1 + Z * u^2)
	gx1 = F.Add(x1, m.AA)          // 8.  gx1 = x1 + A
	gx1 = F.Mul(gx1, x1)           // 9.  gx1 = gx1 * x1
	gx1 = F.Add(gx1, m.BB)         // 10. gx1 = gx1 + B
	gx1 = F.Mul(gx1, x1)           // 11. gx1 = gx1 * x1            // gx1 = x1^3 + A * x1^2 + B * x1
	x2 = F.Sub(F.Neg(x1), m.AA)    // 12.  x2 = -x1 - A
	gx2 = F.Mul(t1, gx1)           // 13. gx2 = t1 * gx1
	e2 = F.IsSquare(gx1)           // 14.  e2 = is_square(gx1)
	x = F.CMov(x2, x1, e2)         // 15.   x = CMOV(x2, x1, e2)    // If is_square(gx1), x = x1, else x = x2
	y2 = F.CMov(gx2, gx1, e2)      // 16.  y2 = CMOV(gx2, gx1, e2)  // If is_square(gx1), y2 = gx1, else y2 = gx2
	y = F.Sqrt(y2)                 // 17.   y = sqrt(y2)
	e3 = m.Sgn0(u) == m.Sgn0(y)    // 18.  e3 = sgn0(u) == sgn0(y)  // Fix sign of y
	y = F.CMov(F.Neg(y), y, e3)    // 19.   y = CMOV(-y, y, e3)
	x = F.Mul(x, m.BB)
	y = F.Mul(y, m.BB)
	return m.E.NewPoint(x, y)
}
