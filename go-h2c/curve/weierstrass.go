package curve

import (
	"fmt"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type ecWe struct{ *Params }

// Curve is a Weierstrass curve
type WECurve = *ecWe

func (e *ecWe) String() string {
	return fmt.Sprintf("y^2=x^3+Ax+B\nF: %v\nA: %v\nD: %v\n", e.F, e.A, e.B)
}
func (e *ecWe) IsOnCurve(p Point) bool {
	if _, isZero := p.(*infPoint); isZero {
		return isZero
	}
	P := p.(*ptWe)
	F := e.F
	var t0, t1 GF.Elt
	t0 = F.Sqr(P.x)     // x^2
	t0 = F.Add(t0, e.A) // x^2+A
	t0 = F.Mul(t0, P.x) // (x^2+A)x
	t0 = F.Add(t0, e.B) // (x^2+A)x+B
	t1 = F.Sqr(P.y)     // y^2
	return F.AreEqual(t0, t1)
}

// NewWeierstrass returns a Weierstrass curve
func NewWeierstrass(ecParams *Params) EllCurve { return &ecWe{ecParams} }

func (e *ecWe) NewPoint(x, y GF.Elt) Point { return &ptWe{e, &afPoint{x: x, y: y}} }
func (e *ecWe) Identity() Point            { return &infPoint{} }
func (e *ecWe) Add(p, q Point) Point {
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
func (e *ecWe) Neg(p Point) Point {
	if _, isZero := p.(*infPoint); isZero {
		return e.Identity()
	}
	P := p.(*ptWe)
	return &ptWe{e, &afPoint{x: P.x.Copy(), y: e.F.Neg(P.y)}}
}
func (e *ecWe) add(p, q Point) Point {
	P := p.(*ptWe)
	Q := q.(*ptWe)
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
	t0 = F.Sub(t0, P.x) // l^2-x1
	x := F.Sub(t0, Q.x) // x' = l^2-x1-x2

	t0 = F.Sub(P.x, x)  // x1-x3
	t0 = F.Mul(t0, ll)  // l(x1-x3)
	y := F.Sub(t0, P.y) // y3 = l(x1-x3)-y1

	return &ptWe{e, &afPoint{x: x, y: y}}
}
func (e *ecWe) Double(p Point) Point {
	if _, ok := p.(*infPoint); ok {
		return e.Identity()
	}
	P := p.(*ptWe)
	if P.IsTwoTorsion() {
		return e.Identity()
	}

	F := e.F
	var t0, t1, ll GF.Elt
	t0 = F.Sqr(P.x)          // x^2
	t0 = F.Mul(t0, F.Elt(3)) // 3x^2
	t0 = F.Add(t0, e.A)      // 3x^2+A
	t1 = F.Add(P.y, P.y)     // 2y
	t1 = F.Inv(t1)           // 1/2y
	ll = F.Mul(t0, t1)       // l = (3x^2+2A)/(2y)

	t0 = F.Sqr(ll)      // l^2
	t0 = F.Sub(t0, P.x) // l^2-x
	x := F.Sub(t0, P.x) // x' = l^2-2x

	t0 = F.Sub(P.x, x)  // x-x'
	t0 = F.Mul(t0, ll)  // l(x-x')
	y := F.Sub(t0, P.y) // y3 = l(x-x')-y1

	return &ptWe{e, &afPoint{x: x, y: y}}
}

// ptWe is an affine point on a Weierstrass curve.
type ptWe struct {
	*ecWe
	*afPoint
}

func (p *ptWe) Copy() Point {
	return &ptWe{p.ecWe, &afPoint{x: p.x.Copy(), y: p.y.Copy()}}
}
func (p *ptWe) IsEqual(q Point) bool {
	qq := q.(*ptWe)
	return p.ecWe == qq.ecWe && p.isEqual(p.F, qq.afPoint)
}
func (p *ptWe) IsIdentity() bool   { return false }
func (p *ptWe) IsTwoTorsion() bool { return p.F.IsZero(p.y) }
