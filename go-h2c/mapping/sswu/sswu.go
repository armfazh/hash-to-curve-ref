package sswu

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

type extField interface {
	GF.Field
	GF.HasInv0
	GF.HasSgn0
	GF.HasCMov
	GF.HasSqrt
}

type sswu struct {
	E      C.WECurve
	F      extField
	Z      GF.Elt
	c1, c2 GF.Elt
	Sgn0   func(GF.Elt) int
}

func (m sswu) String() string { return fmt.Sprintf("Simple SWU for E: %v", m.E) }

// New is
func New(e C.EllCurve, z GF.Elt, sgn0 string) M.MapToCurve {
	if s := (&sswu{E: e.(C.WECurve), F: e.Field().(extField), Z: z}); s.verify() {
		s.precmp(sgn0)
		return s
	}
	panic(fmt.Errorf("Failed restrictions for sswu"))
}

func (m *sswu) precmp(sgn0 string) {
	F := m.E.F
	switch sgn0 {
	case "le":
		m.Sgn0 = m.F.Sgn0LE
	case "be":
		m.Sgn0 = m.F.Sgn0BE
	}

	t0 := F.Inv(m.E.A)    // 1/A
	t0 = F.Mul(t0, m.E.B) // B/A
	m.c1 = F.Neg(t0)      // -B/A
	t0 = F.Inv(m.Z)       // 1/Z
	m.c2 = F.Neg(t0)      // -1/Z
}

func (m sswu) verify() bool {
	F := m.E.F
	precond1 := !F.IsZero(m.E.A)         // A != 0
	precond2 := !F.IsZero(m.E.B)         // B != 0
	cond1 := !F.IsSquare(m.Z)            // Z is non-square
	cond2 := !F.AreEqual(m.Z, F.Elt(-1)) // Z != -1
	t0 := F.Mul(m.Z, m.E.A)              // Z*A
	t0 = F.Inv(t0)                       // 1/(Z*A)
	t0 = F.Mul(t0, m.E.B)                // B/(Z*A)
	g := m.E.EvalRHS(t0)                 // g(B/(Z*A))
	cond4 := F.IsSquare(g)               // g(B/(Z*A)) is square
	return precond1 && precond2 && cond1 && cond2 && cond4
}

func (m *sswu) Map(u GF.Elt) C.Point {
	F := m.F
	var t1, t2 GF.Elt
	var x1, x2, gx1, gx2, y2, x, y GF.Elt
	var e1, e2, e3 bool

	t1 = F.Sqr(u)               // 0.   t1 = u^2
	t1 = F.Mul(t1, m.Z)         // 1.   t1 = Z * u^2
	t2 = F.Sqr(t1)              // 2.   t2 = t1^2
	x1 = F.Add(t1, t2)          // 3.   x1 = t1 + t2
	x1 = F.Inv0(x1)             // 4.   x1 = inv0(x1)
	e1 = F.IsZero(x1)           // 5.   e1 = x1 == 0
	x1 = F.Add(x1, F.One())     // 6.   x1 = x1 + 1
	x1 = F.CMov(x1, m.c2, e1)   // 7.   x1 = CMOV(x1, c2, e1)
	x1 = F.Mul(x1, m.c1)        // 8.   x1 = x1 * c1
	gx1 = F.Sqr(x1)             // 9.  gx1 = x1^2
	gx1 = F.Add(gx1, m.E.A)     // 10. gx1 = gx1 + A
	gx1 = F.Mul(gx1, x1)        // 11. gx1 = gx1 * x1
	gx1 = F.Add(gx1, m.E.B)     // 12. gx1 = gx1 + B
	x2 = F.Mul(t1, x1)          // 13.  x2 = t1 * x1
	t2 = F.Mul(t1, t2)          // 14.  t2 = t1 * t2
	gx2 = F.Mul(gx1, t2)        // 15. gx2 = gx1 * t2
	e2 = F.IsSquare(gx1)        // 16.  e2 = is_square(gx1)
	x = F.CMov(x2, x1, e2)      // 17.   x = CMOV(x2, x1, e2)
	y2 = F.CMov(gx2, gx1, e2)   // 18.  y2 = CMOV(gx2, gx1, e2)
	y = F.Sqrt(y2)              // 19.   y = sqrt(y2)
	e3 = m.Sgn0(u) == m.Sgn0(y) // 20.  e3 = sgn0(u) == sgn0(y)
	y = F.CMov(F.Neg(y), y, e3) // 21.   y = CMOV(-y, y, e3)
	return m.E.NewPoint(x, y)
}
