package weierstrass

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type curve struct {
	F        GF.Field
	A, B     GF.Elt
	Cofactor uint
}

// NewCurve returns a Weierstrass curve
func NewCurve(f GF.Field, a, b GF.Elt, h uint) C.EllCurve {
	return &curve{F: f, A: a, B: b, Cofactor: h}
}

// NewCurve returns a point on a Weierstrass curve
func (e *curve) NewPoint(x, y GF.Elt) C.Point {
	return &point{e: e, x: x, y: y, isIdentity: false}
}

func (e *curve) Identity() C.Point {
	return &point{e: e, x: e.F.Zero(), y: e.F.Zero(), isIdentity: true}
}

func (e *curve) IsOnCurve(p C.Point) bool {
	P := p.(*point)
	F := e.F
	var t0, t1 GF.Elt
	t0 = F.Sqr(P.x)     // x^2
	t0 = F.Add(t0, e.A) // x^2+A
	t0 = F.Mul(t0, P.x) // (x^2+A)x
	t0 = F.Add(t0, e.B) // (x^2+A)x+B
	t1 = F.Sqr(P.y)     // y^2
	return F.AreEqual(t0, t1) || P.isIdentity
}
func (e *curve) Add(p, q C.Point) C.Point {
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

func (e *curve) Neg(p C.Point) C.Point {
	P := p.(*point)
	return &point{
		e:          P.e,
		x:          P.x.Copy(),
		y:          P.e.F.Neg(P.y),
		isIdentity: P.isIdentity,
	}
}

func (e *curve) add(p, q C.Point) C.Point {
	P := p.(*point)
	Q := q.(*point)
	F := e.F

	if F.AreEqual(P.x, Q.x) {
		panic("wrong inputs")
	}

	R := &point{e: e, isIdentity: false}
	var t0, t1, ll GF.Elt
	t0 = F.Sub(Q.y, P.y) // (y2-y1)
	t1 = F.Sub(Q.x, P.x) // (x2-x1)
	t1 = F.Inv(t1)       // 1/(x2-x1)
	ll = F.Mul(t0, t1)   // l = (y2-y1)/(x2-x1)

	t0 = F.Sqr(ll)       // l^2
	t0 = F.Sub(t0, P.x)  // l^2-x1
	R.x = F.Sub(t0, Q.x) // x' = l^2-x1-x2

	t0 = F.Sub(P.x, R.x) // x1-x3
	t0 = F.Mul(t0, ll)   // l(x1-x3)
	R.y = F.Sub(t0, P.y) // y3 = l(x1-x3)-y1

	return R
}

func (e *curve) Double(p C.Point) C.Point {
	P := p.(*point)
	if P.IsTwoTorsion() {
		return e.Identity()
	}

	F := e.F
	R := &point{e: e, isIdentity: false}
	var t0, t1, ll GF.Elt
	t0 = F.Sqr(P.x)          // x^2
	t0 = F.Mul(t0, F.Elt(3)) // 3x^2
	t0 = F.Add(t0, e.A)      // 3x^2+A
	t1 = F.Add(P.y, P.y)     // 2y
	t1 = F.Inv(t1)           // 1/2y
	ll = F.Mul(t0, t1)       // l = (3x^2+2A)/(2y)

	t0 = F.Sqr(ll)       // l^2
	t0 = F.Sub(t0, P.x)  // l^2-x
	R.x = F.Sub(t0, P.x) // x' = l^2-2x

	t0 = F.Sub(P.x, R.x) // x-x'
	t0 = F.Mul(t0, ll)   // l(x-x')
	R.y = F.Sub(t0, P.y) // y3 = l(x-x')-y1

	return R
}

func (e curve) String() string {
	return fmt.Sprintf("y^2=x^3+Ax+B\nF: %v\nA: %v\nB: %v\n", e.F, e.A, e.B)
}

// Point is a projective point on a Montgomery curve.
type point struct {
	e          *curve
	x, y       GF.Elt
	isIdentity bool
}

func (p point) String() string {
	if p.isIdentity {
		return "(inf)"
	}
	return fmt.Sprintf("(%v, %v)", p.x, p.y)
}
func (p *point) IsIdentity() bool   { return p.isIdentity }
func (p *point) IsTwoTorsion() bool { return p.e.F.IsZero(p.y) || p.isIdentity }

func (p *point) IsEqual(q C.Point) bool {
	qq := q.(*point)
	return p.e == qq.e &&
		p.isIdentity == qq.isIdentity &&
		p.e.F.AreEqual(p.x, qq.x) &&
		p.e.F.AreEqual(p.y, qq.y)
}

func (p *point) Copy() C.Point {
	q := &point{}
	q.e = p.e
	q.x = p.x.Copy()
	q.y = p.y.Copy()
	q.isIdentity = p.isIdentity
	return q
}
