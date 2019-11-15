package curve

import (
	"fmt"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type ecTe struct{ *Params }

// TECurve is a twisted Edwards curve
type TECurve = *ecTe

// NewEdwards returns a twisted Edwards curve
func NewEdwards(ecParams *Params) EllCurve { return &ecTe{ecParams} }

func (e *ecTe) String() string {
	return fmt.Sprintf("Ax^2+y^2=1+Dx^2y^2\nF: %v\nA: %v\nD: %v\n", e.F, e.A, e.D)
}
func (e *ecTe) IsOnCurve(p Point) bool {
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

func (e *ecTe) NewPoint(x, y GF.Elt) Point { return &ptTe{e, &afPoint{x: x, y: y}} }
func (e *ecTe) Identity() Point            { return e.NewPoint(e.F.Zero(), e.F.One()) }
func (e *ecTe) Add(p, q Point) Point {
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
func (e *ecTe) Double(p Point) Point { return e.Add(p, p) }
func (e *ecTe) Neg(p Point) Point {
	P := p.(*ptTe)
	return &ptTe{e, &afPoint{x: e.F.Neg(P.x), y: P.y.Copy()}}
}

type ptTe struct {
	*ecTe
	*afPoint
}

func (p *ptTe) Copy() Point {
	return &ptTe{p.ecTe, &afPoint{x: p.x.Copy(), y: p.y.Copy()}}
}
func (p *ptTe) IsEqual(q Point) bool {
	qq := q.(*ptTe)
	return p.ecTe == qq.ecTe && p.isEqual(p.F, qq.afPoint)
}
func (p *ptTe) IsIdentity() bool {
	return p.F.IsZero(p.x) && p.F.AreEqual(p.y, p.F.One())
}
