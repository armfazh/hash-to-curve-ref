package curve

import (
	"fmt"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type mt2wec struct {
	E0   *MTCurve
	E1   *WCCurve
	invB GF.Elt
}

func (e *MTCurve) ToWeierstrassC() RationalMap {
	F := e.Field()
	invB := F.Inv(e.params.B)
	a := F.Mul(invB, e.params.A)
	b := F.Sqr(invB)
	return &mt2wec{E0: e, E1: NewWeierstrassC(Custom, F, a, b, e.params.R, e.params.H), invB: invB}
}

func (r *mt2wec) Domain() EllCurve   { return r.E0 }
func (r *mt2wec) Codomain() EllCurve { return r.E1 }
func (r *mt2wec) Push(p Point) Point {
	if p.IsIdentity() {
		return r.E1.Identity()
	}
	F := r.E0.Field()

	P := p.(*ptMt)
	x := F.Mul(P.x, r.invB) // s = x/B
	y := F.Mul(P.y, r.invB) // t = y/B
	return r.E1.NewPoint(x, y)
}
func (r *mt2wec) Pull(p Point) Point {
	if p.IsIdentity() {
		return r.E0.Identity()
	}
	F := r.E0.Field()
	P := p.(*ptWc)
	x := F.Mul(P.x, r.E0.B) // x = s*B
	y := F.Mul(P.y, r.E0.B) // y = t*B
	return r.E0.NewPoint(x, y)
}

type te2wec struct {
	E0       *TECurve
	E1       *WCCurve
	invSqrtD GF.Elt // 4/(a-d)
}

func (e *TECurve) ToWeierstrassC() RationalMap {
	F := e.Field()
	half := F.Inv(F.Elt(2))             // 1/2
	t0 := F.Add(e.params.A, e.params.D) // a+d
	a := F.Mul(t0, half)                // A = (a+d)/2

	t0 = F.Sub(e.params.A, e.params.D) // a-d
	t0 = F.Mul(t0, half)               // (a-d)/2
	t0 = F.Mul(t0, half)               // (a-d)/4
	invSqrtD := F.Inv(t0)              // 4/(a-d)
	b := F.Sqr(t0)                     // B = (a-d)^2/16
	return &te2wec{E0: e, E1: NewWeierstrassC(Custom, F, a, b, e.params.R, e.params.H), invSqrtD: invSqrtD}
}

func (r *te2wec) Domain() EllCurve   { return r.E0 }
func (r *te2wec) Codomain() EllCurve { return r.E1 }
func (r *te2wec) Push(p Point) Point {
	if p.IsIdentity() {
		return r.E1.Identity()
	}
	F := r.E0.Field()
	P := p.(*ptTe)
	t0 := F.Add(F.One(), P.y)  // 1+y
	t1 := F.Sub(F.One(), P.y)  // 1-y
	t1 = F.Mul(t1, r.invSqrtD) // invSqrtD*(1-y)
	t1 = F.Inv(t1)             // 1/(invSqrtD*(1-y))
	x := F.Mul(t0, t1)         // x = (1+y)/(invSqrtD*(1-y))
	t0 = F.Inv(P.y)            // 1/y
	y := F.Mul(x, t0)          // y = x/y
	return r.E1.NewPoint(x, y)
}
func (r *te2wec) Pull(p Point) Point {
	if p.IsIdentity() {
		return r.E0.Identity()
	}
	P := p.(*ptWc)
	F := r.E0.Field()
	if P.IsTwoTorsion() {
		return r.E0.NewPoint(F.Zero(), F.Elt(-1))
	}
	t0 := F.Inv(P.y)            // 1/y
	x := F.Mul(P.x, t0)         // X = x/y
	t0 = F.Mul(r.invSqrtD, P.x) // invSqrtD*x
	t1 := F.Add(t0, F.One())    // invSqrtD*x+1
	t2 := F.Sub(t0, F.One())    // invSqrtD*x-1
	t1 = F.Inv(t1)              // 1/(invSqrtD*x+1)
	y := F.Mul(t1, t2)          // Y = (invSqrtD*x-1)/(invSqrtD*x+1)
	return r.E0.NewPoint(x, y)
}

type wc2we struct {
	E0    *WCCurve
	E1    *WECurve
	Adiv3 GF.Elt
}

func (e *WCCurve) ToWeierstrass() RationalMap {
	F := e.Field()
	var t0, t1 GF.Elt
	t0 = F.Inv(F.Elt(3))    // 1/3
	Adiv3 := F.Mul(t0, e.A) // A/3
	t0 = F.Neg(Adiv3)       // -A/3
	t0 = F.Mul(t0, e.A)     // -A^2/3
	A := F.Add(t0, e.B)     // -A^2/3 + B

	t0 = F.Mul(F.Elt(9), e.B) // 9B
	t1 = F.Sqr(e.A)           // A^2
	t1 = F.Add(t1, t1)        // 2A^2
	t1 = F.Sub(t1, t0)        // 2A^2 - 9B
	t1 = F.Mul(t1, e.A)       // A(2A^2 - 9B)
	t0 = F.Inv(F.Elt(27))     // 1/27
	B := F.Mul(t0, t1)        // A(2A^2 - 9B)/27
	return &wc2we{E0: e, E1: NewWeierstrass(Custom, F, A, B, e.params.R, e.params.H), Adiv3: Adiv3}
}
func (r *wc2we) Domain() EllCurve   { return r.E0 }
func (r *wc2we) Codomain() EllCurve { return r.E1 }
func (r *wc2we) Push(p Point) Point {
	if p.IsIdentity() {
		return r.E1.Identity()
	}
	P := p.(*ptWc)
	F := r.E0.Field()
	xx := F.Add(P.x, r.Adiv3)
	return r.E1.NewPoint(xx, P.y)
}
func (r *wc2we) Pull(p Point) Point {
	if p.IsIdentity() {
		return r.E0.Identity()
	}
	P := p.(*ptWe)
	F := r.E0.Field()
	xx := F.Sub(P.x, r.Adiv3)
	return r.E0.NewPoint(xx, P.y)
}

type te2mt25519 struct {
	E0       T
	E1       M
	invSqrtD GF.Elt // sqrt(-486664) such that sgn0(sqrt_neg_486664) == 1
}

// FromTe2Mt25519 returns the birational map between Edwards25519 and Curve25519 curves.
func FromTe2Mt25519() RationalMap {
	e0 := Edwards25519.Get()
	e1 := Curve25519.Get()
	F := e0.Field()
	return te2mt25519{
		E0:       e0.(T),
		E1:       e1.(M),
		invSqrtD: F.Elt("6853475219497561581579357271197624642482790079785650197046958215289687604742"),
	}
}
func (m te2mt25519) String() string     { return fmt.Sprintf("Rational Map from %v to\n%v", m.E0, m.E1) }
func (m te2mt25519) Domain() EllCurve   { return m.E0 }
func (m te2mt25519) Codomain() EllCurve { return m.E1 }
func (m te2mt25519) Push(p Point) Point {
	if p.IsIdentity() {
		return m.E1.Identity()
	}
	F := m.E0.Field()
	x0, y0, one := p.X(), p.Y(), F.One()
	t0 := F.Add(one, y0)       // 1+y
	t1 := F.Sub(one, y0)       // 1-y
	t1 = F.Inv(t1)             // 1/(1-y)
	xx := F.Mul(t0, t1)        // xx = (1+y)/(1-y)
	t0 = F.Inv(x0)             // 1/x
	t0 = F.Mul(t0, m.invSqrtD) // invSqrtD/x
	yy := F.Mul(t0, xx)        // yy = invSqrtD*xx/x
	return m.E1.NewPoint(xx, yy)
}
func (m te2mt25519) Pull(p Point) Point {
	if p.IsIdentity() {
		return m.E0.Identity()
	}
	F := m.E0.Field()
	if p.IsTwoTorsion() {
		return m.E0.NewPoint(F.Zero(), F.Elt(-1))
	}
	x0, y0, one := p.X(), p.Y(), F.One()
	t0 := F.Inv(y0)             // 1/y
	t0 = F.Mul(t0, x0)          // x/y
	xx := F.Mul(m.invSqrtD, t0) // xx = invSqrtD*x/y
	t0 = F.Add(x0, one)         // x+1
	t1 := F.Sub(x0, one)        // x-1
	t0 = F.Inv(t0)              // 1/(x+1)
	yy := F.Mul(t0, t1)         // yy = (x-1)/(x+1)
	return m.E0.NewPoint(xx, yy)
}

type te2mt4iso448 struct {
	E0 T
	E1 M
}

// FromTe2Mt4ISO448 returns the four-degree isogeny between Edwards448 and Curve448 curves.
func FromTe2Mt4ISO448() RationalMap       { return te2mt4iso448{Edwards448.Get().(T), Curve448.Get().(M)} }
func (m te2mt4iso448) String() string     { return fmt.Sprintf("4-Isogeny from %v to\n%v", m.E0, m.E1) }
func (m te2mt4iso448) Domain() EllCurve   { return m.E0 }
func (m te2mt4iso448) Codomain() EllCurve { return m.E1 }
func (m te2mt4iso448) Push(p Point) Point {
	if p.IsIdentity() {
		return m.E1.Identity()
	}
	F := m.E0.Field()
	x, y := p.X(), p.Y()
	t0 := F.Inv(x)           // 1/x
	t1 := F.Sqr(t0)          // 1/x^2
	t2 := F.Mul(t0, t1)      // 1/x^3
	t0 = F.Sqr(y)            // y^2
	xx := F.Mul(t0, t1)      // xx = y^2/x^2
	t0 = F.Sub(F.Elt(2), t0) // 2-y^2
	t1 = F.Sqr(x)            // x^2
	t0 = F.Sub(t0, t1)       // 2-y^2-x^2
	t0 = F.Mul(t0, y)        // (2-y^2-x^2)*y
	yy := F.Mul(t0, t2)      // (2-y^2-x^2)*y/x^3
	return m.E1.NewPoint(xx, yy)
}
func (m te2mt4iso448) Pull(p Point) Point {
	if p.IsIdentity() {
		return m.E0.Identity()
	}
	F := m.E0.Field()
	x, y, one := p.X(), p.Y(), F.One()

	t0 := F.Sqr(x)       // x^2
	t1 := F.Add(t0, one) // x^2+1
	t0 = F.Sub(t0, one)  // x^2-1
	t2 := F.Sqr(y)       // y^2
	t2 = F.Add(t2, t2)   // 2y^2
	t3 := F.Add(x, x)    // 2x

	t4 := F.Mul(t0, y)    // y(x^2-1)
	t4 = F.Add(t4, t4)    // 2y(x^2-1)
	xNum := F.Add(t4, t4) // xNum = 4y(x^2-1)

	t5 := F.Sqr(t0)       // x^4-2x^2+1
	t4 = F.Add(t5, t2)    // x^4-2x^2+1+2y^2
	xDen := F.Add(t4, t2) // xDen = x^4-2x^2+1+4y^2

	t5 = F.Mul(t5, x)     // x^5-2x^3+x
	t4 = F.Mul(t2, t3)    // 4xy^2
	yNum := F.Sub(t4, t5) // yNum = -(x^5-2x^3+x-4xy^2)

	t4 = F.Mul(t1, t2)    // 2x^2y^2+2y^2
	yDen := F.Sub(t5, t4) // yDen = x^5-2x^3+x-2x^2y^2-2y^2

	xx := F.Mul(xNum, F.Inv(xDen))
	yy := F.Mul(yNum, F.Inv(yDen))
	return m.E0.NewPoint(xx, yy)
}
