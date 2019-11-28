package curve

import (
	"fmt"
	"math/big"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// WCCurve is a Weierstrass curve
type WCCurve struct {
	*params
	e1    W
	Adiv3 GF.Elt
}

type WC = *WCCurve

func (e *WCCurve) String() string { return "y^2=x^3+Ax^2+Bx\n" + e.params.String() }

// NewWeierstrassC returns a Weierstrass curve
func NewWeierstrassC(f GF.Field, a, b GF.Elt, r, h *big.Int) *WCCurve {
	var t0, t1, t2 GF.Elt

	t0 = f.Inv(f.Elt(3)) // 1/3
	t2 = f.Mul(t0, a)    // A/3
	t0 = f.Neg(t2)       // -A/3
	t0 = f.Mul(t0, a)    // -A^2/3
	AA := f.Add(t0, b)   // -A^2/3 + B

	t0 = f.Mul(f.Elt(9), b) // 9B
	t1 = f.Sqr(a)           // A^2
	t1 = f.Add(t1, t1)      // 2A^2
	t1 = f.Sub(t1, t0)      // 2A^2 - 9B
	t1 = f.Mul(t1, a)       // A(2A^2 - 9B)
	t0 = f.Inv(f.Elt(27))   // 1/27
	BB := f.Mul(t0, t1)     // A(2A^2 - 9B)/27

	fmt.Printf("A: %v\n", AA)
	fmt.Printf("B: %v\n", BB)
	fmt.Printf("A/3: %v\n", t2)
	return &WCCurve{
		params: &params{F: f, A: a, B: b, R: r, H: h},
		e1:     NewWeierstrass(f, AA, BB, r, h),
		Adiv3:  t2,
	}
}
func (e *WCCurve) Codomain() EllCurve { return e.e1 }
func (e *WCCurve) Domain() EllCurve   { return e }
func (e *WCCurve) Push(p Point) Point {
	if P, ok := p.(*ptWc); ok {
		xx := e.F.Sub(P.x, e.Adiv3)
		return &ptWe{e.e1, &afPoint{xx, P.y}}
	}
	panic("point is not in the domain curve")
}
func (e *WCCurve) Pull(p Point) Point {
	if P, ok := p.(*ptWe); ok {
		xx := e.F.Add(P.x, e.Adiv3)
		return &ptWc{e, &afPoint{xx, P.y}}
	}
	panic("point is not in the codomain curve")
}

func (e *WCCurve) NewPoint(x, y GF.Elt) (P Point) {
	if P = (&ptWc{e, &afPoint{x: x, y: y}}); e.IsOnCurve(P) {
		return P
	}
	panic(fmt.Errorf("%v not on %v", P, e))
}

func (e *WCCurve) Identity() Point             { return &infPoint{} }
func (e *WCCurve) IsOnCurve(p Point) bool      { return e.e1.IsOnCurve(e.Push(p)) }
func (e *WCCurve) Add(p, q Point) Point        { return e.Pull(e.e1.Add(e.Push(p), e.Push(q))) }
func (e *WCCurve) Double(p Point) Point        { return e.Pull(e.e1.Double(e.Push(p))) }
func (e *WCCurve) Neg(p Point) Point           { return e.Pull(e.e1.Neg(e.Push(p))) }
func (e *WCCurve) ClearCofactor(p Point) Point { return e.Pull(e.e1.ClearCofactor(e.Push(p))) }

// ptWc is an affine point on a WCCurve curve.
type ptWc struct {
	*WCCurve
	*afPoint
}

func (p *ptWc) String() string { return p.afPoint.String() }
func (p *ptWc) Copy() Point    { return &ptWc{p.WCCurve, p.copy()} }
func (p *ptWc) IsEqual(q Point) bool {
	qq := q.(*ptWc)
	return p.WCCurve == qq.WCCurve && p.isEqual(p.F, qq.afPoint)
}
func (p *ptWc) IsIdentity() bool { return false }
