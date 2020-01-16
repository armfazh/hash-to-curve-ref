package curve

import (
	"errors"
	"fmt"
	"math/big"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// TECurve is a twisted Edwards curve
type TECurve struct{ *params }

type T = *TECurve

func (e *TECurve) String() string {
	return fmt.Sprintf("Ax^2+y^2=1+Dx^2y^2\nF: %v\nA: %v\nD: %v\n", e.F, e.A, e.D)
}

// NewEdwards returns a twisted Edwards curve
func NewEdwards(id CurveID, f GF.Field, a, d GF.Elt, r, h *big.Int) *TECurve {
	if e := (&TECurve{&params{
		Id: id, F: f, A: a, D: d, R: r, H: h,
	}}); e.IsValid() {
		return e
	}
	panic(errors.New("can't instantiate a twisted Edwards curve"))
}

func (e *TECurve) NewPoint(x, y GF.Elt) (P Point) {
	if P = (&ptTe{e, &afPoint{x: x, y: y}}); e.IsOnCurve(P) {
		return P
	}
	panic(fmt.Errorf("p:%v not on %v", P, e))
}
func (e *TECurve) IsValid() bool {
	F := e.F
	cond1 := !F.AreEqual(e.A, e.D) // A != D
	cond2 := !F.IsZero(e.A)        // A != 0
	cond3 := !F.IsZero(e.D)        // D != 0
	return cond1 && cond2 && cond3
}
func (e *TECurve) IsComplete() bool {
	F := e.F
	return F.IsSquare(e.A) && !F.IsSquare(e.D) // A != D
}
func (e *TECurve) IsOnCurve(p Point) bool {
	P := p.(*ptTe)
	F := e.F
	var t0, t1, t2 GF.Elt
	t0 = F.Sqr(P.x)         // x^2
	t1 = F.Sqr(P.y)         // y^2
	t2 = F.Mul(t0, t1)      // x^2y^2
	t2 = F.Mul(t2, e.D)     // Dx^2y^2
	t2 = F.Add(t2, F.One()) // 1+Dx^2y^2
	t0 = F.Mul(t0, e.A)     // Ax^2
	t0 = F.Add(t0, t1)      // Ax^2+y^2
	return F.AreEqual(t0, t2)
}
func (e *TECurve) Identity() Point { return e.NewPoint(e.F.Zero(), e.F.One()) }
func (e *TECurve) Add(p, q Point) Point {
	P := p.(*ptTe)
	Q := q.(*ptTe)
	F := e.F

	var t0, t1, t2, t3 GF.Elt
	t0 = F.Mul(e.D, P.x)    // Dx1
	t0 = F.Mul(t0, P.y)     // Dx1y1
	t0 = F.Mul(t0, Q.x)     // Dx1y1x2
	t0 = F.Mul(t0, Q.y)     // Dx1y1x2y2
	t2 = F.Add(F.One(), t0) // 1+Dx1y1x2y2
	t3 = F.Sub(F.One(), t0) // 1-Dx1y1x2y2
	t2 = F.Inv(t2)          // 1/(1+Dx1y1x2y2)
	t3 = F.Inv(t3)          // 1/(1-Dx1y1x2y2)

	t0 = F.Mul(P.x, Q.y) // x1y2
	t1 = F.Mul(Q.x, P.y) // x2y1
	t0 = F.Add(t0, t1)   // x1y2+x2y1
	x := F.Mul(t0, t2)   // (x1y2+x2y1)/(1+Dx1y1x2y2)

	t0 = F.Mul(P.y, Q.y) // y1y2
	t1 = F.Mul(P.x, Q.x) // x1x2
	t1 = F.Mul(t1, e.A)  // Ax1x2
	t0 = F.Sub(t0, t1)   // y1y2-Ax1x2
	y := F.Mul(t0, t3)   // (y1y2-Ax1x2)/(1-Dx1y1x2y2)

	return &ptTe{e, &afPoint{x: x, y: y}}
}
func (e *TECurve) Neg(p Point) Point {
	P := p.(*ptTe)
	return &ptTe{e, &afPoint{x: e.F.Neg(P.x), y: P.y.Copy()}}
}
func (e *TECurve) Double(p Point) Point { return e.Add(p, p) }
func (e *TECurve) ScalarMult(p Point, k *big.Int) Point {
	Q := e.Identity()
	for i := k.BitLen() - 1; i >= 0; i-- {
		Q = e.Double(Q)
		if k.Bit(i) != 0 {
			Q = e.Add(Q, p)
		}
	}
	return Q
}
func (e *TECurve) ClearCofactor(p Point) Point { return e.ScalarMult(p, e.H) }

type ptTe struct {
	*TECurve
	*afPoint
}

func (p *ptTe) String() string { return p.afPoint.String() }
func (p *ptTe) Copy() Point    { return &ptTe{p.TECurve, p.copy()} }
func (p *ptTe) IsEqual(q Point) bool {
	qq := q.(*ptTe)
	return p.TECurve == qq.TECurve && p.isEqual(p.F, qq.afPoint)
}
func (p *ptTe) IsIdentity() bool   { return p.F.IsZero(p.x) && p.F.AreEqual(p.y, p.F.One()) }
func (p *ptTe) IsTwoTorsion() bool { return p.F.IsZero(p.x) && p.F.AreEqual(p.y, p.F.Elt(-1)) }
