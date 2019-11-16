package curve

import (
	"fmt"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

// afPoint is an affine point.
type afPoint struct{ x, y GF.Elt }

func (p afPoint) String() string { return fmt.Sprintf("(%v, %v)", p.x, p.y) }

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
func (p *infPoint) Copy() Point          { return &infPoint{} }
func (p *infPoint) IsEqual(q Point) bool { _, t := q.(*infPoint); return t }
func (p *infPoint) IsIdentity() bool     { return true }
