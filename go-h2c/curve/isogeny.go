package curve

import "fmt"

type Isogeny struct {
	E0, E1 EllCurve
	Map
}

func (i *Isogeny) String() string      { return fmt.Sprintf("Isogeny from E0:%v to E1:%v", i.E0, i.E1) }
func (i *Isogeny) Domain() EllCurve    { return i.E0 }
func (i *Isogeny) Codomain() EllCurve  { return i.E1 }
func (i *Isogeny) Apply(p Point) Point { return i.Map(i.E0, i.E1, p) }

type Map func(e0, e1 EllCurve, p Point) Point
