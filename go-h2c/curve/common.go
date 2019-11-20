package curve

import (
	"fmt"
	"math/big"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type params struct {
	F       GF.Field
	A, B, D GF.Elt
	R       *big.Int
	H       *big.Int
}

func (e *params) String() string {
	return fmt.Sprintf("F: %v\nA: %v\nB: %v\n", e.F, e.A, e.B)
}
func (e *params) Field() GF.Field    { return e.F }
func (e *params) Order() *big.Int    { return e.R }
func (e *params) Cofactor() *big.Int { return e.H }

// afPoint is an affine point.
type afPoint struct{ x, y GF.Elt }

func (p afPoint) String() string { return fmt.Sprintf("(%v, %v)", p.x, p.y) }
func (p *afPoint) X() GF.Elt     { return p.x }
func (p *afPoint) Y() GF.Elt     { return p.y }
func (p *afPoint) copy() *afPoint {
	q := &afPoint{}
	q.x = p.x.Copy()
	q.y = p.y.Copy()
	return q
}
func (p *afPoint) isEqual(f GF.Field, q *afPoint) bool {
	return f.AreEqual(p.x, q.x) && f.AreEqual(p.y, q.y)
}

// infPoint is the point at infinity.
type infPoint struct{}

func (p infPoint) String() string        { return "(inf)" }
func (p *infPoint) X() GF.Elt            { return nil }
func (p *infPoint) Y() GF.Elt            { return nil }
func (p *infPoint) Copy() Point          { return &infPoint{} }
func (p *infPoint) IsEqual(q Point) bool { _, t := q.(*infPoint); return t }
func (p *infPoint) IsIdentity() bool     { return true }
