package curve

import (
	"errors"
	"fmt"
	"math/big"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// MTCurve is a Montgomery curve
type MTCurve struct{ *params }

type M = *MTCurve

func (e *MTCurve) String() string { return "By^2=x^3+Ax^2+x\n" + e.params.String() }

// NewMontgomery returns a Montgomery curve
func NewMontgomery(id CurveID, f GF.Field, a, b GF.Elt, r, h *big.Int) *MTCurve {
	if e := (&MTCurve{&params{
		Id: id, F: f, A: a, B: b, R: r, H: h,
	}}); e.IsValid() {
		return e
	}
	panic(errors.New("can't instantiate a Montgomery curve"))
}

func (e *MTCurve) NewPoint(x, y GF.Elt) (P Point) {
	if P = (&ptMt{e, &afPoint{x: x, y: y}}); e.IsOnCurve(P) {
		return P
	}
	panic(fmt.Errorf("p:%v not on %v", P, e))
}
func (e *MTCurve) IsValid() bool {
	F := e.F
	t0 := F.Sqr(e.A)         // A^2
	t0 = F.Sub(t0, F.Elt(4)) // A^2-4
	t0 = F.Mul(t0, e.B)      // B(A^2-4)
	return !F.IsZero(t0)     // B(A^2-4) != 0
}
func (e *MTCurve) IsOnCurve(p Point) bool {
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
func (e *MTCurve) Identity() Point { return &infPoint{} }
func (e *MTCurve) Add(p, q Point) Point {
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
func (e *MTCurve) Neg(p Point) Point {
	if _, isZero := p.(*infPoint); isZero {
		return e.Identity()
	}
	P := p.(*ptMt)
	return &ptMt{e, &afPoint{x: P.x.Copy(), y: e.F.Neg(P.y)}}
}
func (e *MTCurve) add(p, q Point) Point {
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
func (e *MTCurve) Double(p Point) Point {
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

func (e *MTCurve) ScalarMult(p Point, k *big.Int) Point {
	Q := e.Identity()
	for i := k.BitLen() - 1; i >= 0; i-- {
		Q = e.Double(Q)
		if k.Bit(i) != 0 {
			Q = e.Add(Q, p)
		}
	}
	return Q
}
func (e *MTCurve) ClearCofactor(p Point) Point { return e.ScalarMult(p, e.H) }

// ptMt is an affine point on a Montgomery curve.
type ptMt struct {
	*MTCurve
	*afPoint
}

func (p *ptMt) String() string { return p.afPoint.String() }
func (p *ptMt) Copy() Point    { return &ptMt{p.MTCurve, p.copy()} }
func (p *ptMt) IsEqual(q Point) bool {
	qq := q.(*ptMt)
	return p.MTCurve == qq.MTCurve && p.isEqual(p.F, qq.afPoint)
}
func (p *ptMt) IsIdentity() bool   { return false }
func (p *ptMt) IsTwoTorsion() bool { return p.F.IsZero(p.y) }
