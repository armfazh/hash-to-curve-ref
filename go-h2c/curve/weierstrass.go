package curve

import (
	"errors"
	"fmt"
	"math/big"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// WECurve is a Weierstrass curve
type WECurve struct{ *params }

type W = *WECurve

func (e *WECurve) String() string { return "y^2=x^3+Ax+B\n" + e.params.String() }

// NewWeierstrass returns a Weierstrass curve
func NewWeierstrass(f GF.Field, a, b GF.Elt, r, h *big.Int) *WECurve {
	if e := (&WECurve{&params{
		F: f, A: a, B: b, R: r, H: h,
	}}); e.IsValid() {
		return e
	}
	panic(errors.New("can't instantiate a Weierstrass curve"))
}

// NewPoint generates
func (e *WECurve) NewPoint(x, y GF.Elt) (P Point) {
	if P = (&ptWe{e, &afPoint{x: x, y: y}}); e.IsOnCurve(P) {
		return P
	}
	panic(fmt.Errorf("p:%v not on curve", P))
}
func (e *WECurve) IsValid() bool {
	F := e.F
	t0 := F.Sqr(e.A)          // A^2
	t0 = F.Mul(t0, e.A)       // A^3
	t0 = F.Add(t0, t0)        // 2A^3
	t0 = F.Add(t0, t0)        // 4A^3
	t1 := F.Sqr(e.B)          // B^3
	t1 = F.Mul(t1, F.Elt(27)) // 27B^2
	t0 = F.Add(t0, t1)        // 4A^3+27B^2
	t0 = F.Add(t0, t0)        // 2(4A^3+27B^2)
	t0 = F.Add(t0, t0)        // 4(4A^3+27B^2)
	t0 = F.Add(t0, t0)        // 8(4A^3+27B^2)
	t0 = F.Add(t0, t0)        // 16(4A^3+27B^2)
	t0 = F.Neg(t0)            // -16(4A^3+27B^2)
	return !F.IsZero(t0)      // -16(4A^3+27B^2) != 0
}
func (e *WECurve) IsOnCurve(p Point) bool {
	if _, isZero := p.(*infPoint); isZero {
		return isZero
	}
	P := p.(*ptWe)
	F := e.F
	t0 := e.EvalRHS(P.x)
	t1 := F.Sqr(P.y) // y^2
	return F.AreEqual(t0, t1)
}
func (e *WECurve) EvalRHS(x GF.Elt) GF.Elt {
	F := e.F
	t0 := F.Sqr(x)        // x^2
	t0 = F.Add(t0, e.A)   // x^2+A
	t0 = F.Mul(t0, x)     // (x^2+A)x
	return F.Add(t0, e.B) // (x^2+A)x+B
}
func (e *WECurve) Identity() Point { return &infPoint{} }
func (e *WECurve) Add(p, q Point) Point {
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
func (e *WECurve) Neg(p Point) Point {
	if _, isZero := p.(*infPoint); isZero {
		return e.Identity()
	}
	P := p.(*ptWe)
	return &ptWe{e, &afPoint{x: P.x.Copy(), y: e.F.Neg(P.y)}}
}
func (e *WECurve) add(p, q Point) Point {
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
func (e *WECurve) Double(p Point) Point {
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
func (e *WECurve) ClearCofactor(p Point) Point {
	return p
}

// ptWe is an affine point on a WECurve curve.
type ptWe struct {
	*WECurve
	*afPoint
}

func (p *ptWe) String() string { return p.afPoint.String() }
func (p *ptWe) Copy() Point    { return &ptWe{p.WECurve, p.copy()} }
func (p *ptWe) IsEqual(q Point) bool {
	qq := q.(*ptWe)
	return p.WECurve == qq.WECurve && p.isEqual(p.F, qq.afPoint)
}
func (p *ptWe) IsIdentity() bool   { return false }
func (p *ptWe) IsTwoTorsion() bool { return p.F.IsZero(p.y) }
