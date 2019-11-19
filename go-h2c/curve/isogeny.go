package curve

import "fmt"

type Isogeny interface {
	Domain() EllCurve
	Codomain() EllCurve
	Degree() uint
	Apply(Point) Point
}

type iso struct {
	E0, E1 EllCurve
	deg    uint
	isoMap
}

func (i iso) String() string      { return fmt.Sprintf("Isogeny from E0:%v to E1:%v", i.E0, i.E1) }
func (i iso) Domain() EllCurve    { return i.E0 }
func (i iso) Codomain() EllCurve  { return i.E1 }
func (i iso) Degree() uint        { return i.deg }
func (i iso) Apply(p Point) Point { return i.isoMap(p) }

type isoMap func(p Point) Point

// NewIsogeny creates a map m from E0 -> E1.
func NewIsogeny(e0, e1 EllCurve, degree uint, m isoMap) Isogeny {
	return iso{E0: e0, E1: e1, isoMap: m, deg: degree}
}
