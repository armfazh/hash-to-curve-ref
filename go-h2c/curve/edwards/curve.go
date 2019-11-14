package edwards

import (
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type curve struct {
	F        GF.Field
	A, D     GF.Elt
	Cofactor uint
}

// NewCurve returns a twisted Edwards curve
func NewCurve(f GF.Field, a, d GF.Elt, h uint) C.EllCurve {
	return &curve{F: f, A: a, D: d, Cofactor: h}
}

// NewCurve returns a point on a twisted Edwards curve
func (e *curve) NewPoint(x, y GF.Elt) C.Point {
	return &point{e: e, x: x, y: y}
}

func (e *curve) Identity() C.Point {
	return &point{e: e, x: e.F.Zero(), y: e.F.One()}
}

func (e *curve) IsOnCurve(p C.Point) bool {
	P := p.(*point)
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

func (e *curve) Add(p, q C.Point) C.Point {
	P := p.(*point)
	Q := q.(*point)
	F := e.F

	R := &point{e: e}
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
	R.x = F.Mul(t0, t2)  // (x1y2+x2y1)/(1+Dx1y1x2y2)

	t0 = F.Mul(P.y, Q.y) // y1y2
	t1 = F.Mul(P.x, Q.x) // x1x2
	t1 = F.Mul(t1, e.A)  // Ax1x2
	t0 = F.Sub(t0, t1)   // y1y2-Ax1x2
	R.y = F.Mul(t0, t3)  // (y1y2-Ax1x2)/(1-Dx1y1x2y2)

	return R
}
func (e *curve) Double(p C.Point) C.Point { return e.Add(p, p) }
func (e *curve) Neg(p C.Point) C.Point {
	P := p.(*point)
	return &point{
		e: P.e,
		x: P.e.F.Neg(P.x),
		y: P.y.Copy(),
	}
}

func (e curve) String() string {
	return fmt.Sprintf("Ax^2+y^2=1+Dx^2y^2\nF: %v\nA: %v\nD: %v\n", e.F, e.A, e.D)
}

// Point is a projective point on a Montgomery curve.
type point struct {
	e    *curve
	x, y GF.Elt
}

func (p point) String() string { return fmt.Sprintf("(%v, %v)", p.x, p.y) }
func (p *point) IsIdentity() bool {
	F := p.e.F
	return F.IsZero(p.x) && F.AreEqual(p.y, F.One())
}

func (p *point) IsEqual(q C.Point) bool {
	qq := q.(*point)
	return p.e == qq.e &&
		p.e.F.AreEqual(p.x, qq.x) &&
		p.e.F.AreEqual(p.y, qq.y)
}

func (p *point) Copy() C.Point {
	q := &point{}
	q.e = p.e
	q.x = p.x.Copy()
	q.y = p.y.Copy()
	return q
}
