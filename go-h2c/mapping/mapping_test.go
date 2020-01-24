package mapping_test

import (
	"testing"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/internal/toy"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

func TestBF(t *testing.T) {
	var curves = []string{"W2"}
	for _, id := range curves {
		E := toy.ToyCurves[id].E
		F := E.Field()
		n := F.Order().Int64()
		m := mapping.NewBF(E)
		for i := int64(0); i < n; i++ {
			u := F.Elt(i)
			P := m.Map(u)
			if !E.IsOnCurve(P) {
				t.Fatalf("u: %v got P: %v\n", u, P)
			}
		}
	}
}

func TestEll2(t *testing.T) {
	var curves = []string{"M0", "M1", "E0", "W3"}
	for _, id := range curves {
		E := toy.ToyCurves[id].E
		F := E.Field()
		n := F.Order().Int64()
		for _, m := range []mapping.MapToCurve{
			mapping.NewElligator2(E, GF.SignLE),
			mapping.NewElligator2(E, GF.SignBE),
		} {
			for i := int64(0); i < n; i++ {
				u := F.Elt(i)
				P := m.Map(u)
				if !E.IsOnCurve(P) {
					t.Fatalf("u: %v got P: %v\n", u, P)
				}
			}
		}
	}
}

func TestSVDW(t *testing.T) {
	var curves = []string{"W0"}
	for _, c := range curves {
		E := toy.ToyCurves[c].E
		F := E.Field()
		n := F.Order().Int64()
		for _, m := range []mapping.MapToCurve{
			mapping.NewSVDW(E, GF.SignLE),
			mapping.NewSVDW(E, GF.SignBE),
		} {
			for i := int64(0); i < n; i++ {
				u := F.Elt(i)
				P := m.Map(u)
				if !E.IsOnCurve(P) {
					t.Fatalf("%vu: %v\nP: %v not on curve.", m, u, P)
				}
			}
		}
	}
}

func TestSSWU(t *testing.T) {
	var curves = []struct {
		Name string
		Z    uint
	}{
		{"W0", 3},
		{"W1iso", 3},
	}
	for _, c := range curves {
		E := toy.ToyCurves[c.Name].E
		F := E.Field()
		n := F.Order().Int64()
		iso := doubleIso{E}
		Z := F.Elt(c.Z)
		for _, m := range []mapping.MapToCurve{
			mapping.NewSSWU(E, Z, GF.SignLE, nil),
			mapping.NewSSWU(E, Z, GF.SignBE, nil),
			mapping.NewSSWU(E, Z, GF.SignLE, iso),
			mapping.NewSSWU(E, Z, GF.SignBE, iso),
		} {
			for i := int64(0); i < n; i++ {
				u := F.Elt(i)
				P := m.Map(u)
				if !E.IsOnCurve(P) {
					t.Fatalf("u: %v got P: %v\n", u, P)
				}
			}
		}
	}
}

type doubleIso struct{ E C.EllCurve }

func (d doubleIso) Domain() C.EllCurve     { return d.E }
func (d doubleIso) Codomain() C.EllCurve   { return d.E }
func (d doubleIso) Push(p C.Point) C.Point { return d.E.Double(p) }
