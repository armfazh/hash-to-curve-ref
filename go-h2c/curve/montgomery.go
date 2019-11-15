package curve

import (
	"fmt"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type ecMt struct{ *Params }

// MTCurve is a Montgomery curve
type MTCurve = *ecMt

// NewMontgomery returns a Montgomery curve
func NewMontgomery(ecParams *Params) EllCurve { return &ecMt{ecParams} }

func (e *ecMt) String() string {
	return fmt.Sprintf("By^2=x^3+Ax^2+x\nF: %v\nA: %v\nB: %v\n", e.F, e.A, e.B)
}
func (e *ecMt) IsOnCurve(p Point) bool {
	if _, isZero := p.(*infPoint); isZero {
		return isZero
	}
	P := p.(*ptMt)
	F := e.F
	var t0, t1 GF.Elt
	t0 = F.Add(P.x, e.A)    // x+A
	t0 = F.Mul(t0, P.x)     // (x+A)x
	t0 = F.Add(t0, F.One()) // (x+A)x+1
	t0 = F.Mul(t0, P.x)     // ((x+A)x+1)x
	t1 = F.Sqr(P.y)         // y^2
	t1 = F.Mul(t1, e.B)     // By^2
	return F.AreEqual(t0, t1)
}

// NewCurve returns a point on a Montgomery curve
func (e *ecMt) NewPoint(x, y GF.Elt) Point { return &ptMt{e, &afPoint{x: x, y: y}} }
func (e *ecMt) Identity() Point            { return &infPoint{} }
func (e *ecMt) Add(p, q Point) Point {
	if p.IsIdentity() {
		return q.Copy()
	} else if q.IsIdentity() {
		return p.Copy()
	} else if p.IsEqual(e.Neg(q)) {
		return e.Identity()
	} else if p.IsEqual(q) {
		return e.Double(p)
	} else {
		return e.add(p, q)
	}
}
func (e *ecMt) Neg(p Point) Point {
	if _, isZero := p.(*infPoint); isZero {
		return e.Identity()
	}
	P := p.(*ptMt)
	return &ptMt{e, &afPoint{x: P.x.Copy(), y: e.F.Neg(P.y)}}
}
func (e *ecMt) add(p, q Point) Point {
	P := p.(*ptMt)
	Q := q.(*ptMt)
	F := e.F

	if F.AreEqual(P.x, Q.x) {
		panic("wrong inputs")
	}

	var t0, t1, ll GF.Elt
	t0 = F.Sub(Q.y, P.y) // (y2-y1)
	t1 = F.Sub(Q.x, P.x) // (x2-x1)
	t1 = F.Inv(t1)       // 1/(x2-x1)
	ll = F.Mul(t0, t1)   // l = (y2-y1)/(x2-x1)

	t0 = F.Sqr(ll)      // l^2
	t0 = F.Mul(t0, e.B) // Bl^2
	t0 = F.Sub(t0, e.A) // Bl^2-A
	t0 = F.Sub(t0, P.x) // Bl^2-A-x1
	x := F.Sub(t0, Q.x) // x' = Bl^2-A-x1-x2

	t0 = F.Sub(P.x, x)  // x1-x3
	t0 = F.Mul(t0, ll)  // l(x1-x3)
	y := F.Sub(t0, P.y) // y3 = l(x1-x3)-y1

	return &ptMt{e, &afPoint{x: x, y: y}}
}
func (e *ecMt) Double(p Point) Point {
	if _, ok := p.(*infPoint); ok {
		return e.Identity()
	}
	P := p.(*ptMt)
	if P.IsTwoTorsion() {
		return e.Identity()
	}

	F := e.F
	var t0, t1, ll GF.Elt
	t0 = F.Mul(F.Elt(3), P.x) // 3x
	t1 = F.Mul(F.Elt(2), e.A) // 2A
	t0 = F.Add(t0, t1)        // 3x+2A
	t0 = F.Mul(t0, P.x)       // (3x+2A)x
	t1 = F.Add(t0, F.One())   // (3x+2A)x+1
	t0 = F.Mul(F.Elt(2), e.B) // 2B
	t0 = F.Mul(t0, P.y)       // 2By
	t0 = F.Inv(t0)            // 1/2By
	ll = F.Mul(t1, t0)        // l = (3x^2+2Ax+1)/(2By)

	t0 = F.Sqr(ll)      // l^2
	t0 = F.Mul(t0, e.B) // Bl^2
	t0 = F.Sub(t0, e.A) // Bl^2-A
	t0 = F.Sub(t0, P.x) // Bl^2-A-x
	x := F.Sub(t0, P.x) // x' = Bl^2-A-2x

	t0 = F.Sub(P.x, x)  // x-x'
	t0 = F.Mul(t0, ll)  // l(x-x')
	y := F.Sub(t0, P.y) // y3 = l(x-x')-y1

	return &ptMt{e, &afPoint{x: x, y: y}}
}

// ptMt is an affine point on a Montgomery curve.
type ptMt struct {
	*ecMt
	*afPoint
}

func (p *ptMt) Copy() Point {
	return &ptMt{p.ecMt, &afPoint{x: p.x.Copy(), y: p.y.Copy()}}
}
func (p *ptMt) IsEqual(q Point) bool {
	qq := q.(*ptMt)
	return p.ecMt == qq.ecMt && p.isEqual(p.F, qq.afPoint)
}
func (p *ptMt) IsIdentity() bool   { return false }
func (p *ptMt) IsTwoTorsion() bool { return p.F.IsZero(p.y) }
