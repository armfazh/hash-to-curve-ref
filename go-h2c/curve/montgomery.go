package curve

import (
	"fmt"

	"github.com/armfazh/hash-to-curve-ref/h2c/math"
)

// MontCurve is
type MontCurve struct {
	F        math.FiniteField
	A, B     math.Element
	Cofactor uint
}

// Add is
func (e MontCurve) Add(P, Q MontgomeryPoint) MontgomeryPoint {
	return MontgomeryPoint{}
}

// Double is
func (e MontCurve) Double(P MontgomeryPoint) MontgomeryPoint {
	return MontgomeryPoint{}
}

// NewHash is
func (e MontCurve) NewHash() HashToPoint {
	z := e.F.Elt()
	return elligator2{&e, z}
}

// MontgomeryPoint is a projective point on a Montgomery curve.
type MontgomeryPoint struct {
	x, y, z math.Element
}

func (P MontgomeryPoint) String() string {
	return fmt.Sprintf("x: %v\n,y: %v\n,z: %v", P.x, P.y, P.z)
}

// IsZero returns true if P is the identity point
func (P MontgomeryPoint) IsZero() bool { return P.z.IsZero() }

// ToAffine is
func (P MontgomeryPoint) ToAffine() {}

type hashToPoint interface {
	Hash([]byte) MontgomeryPoint
	HashRO([]byte) MontgomeryPoint
}

// elligator2 implements a mapping to point on a Montgomery curve.
type elligator2 struct {
	E *MontCurve
	Z math.Element
}

// MapToCurve returns a point
func (ell elligator2) MapToCurve(u math.Element) MontgomeryPoint {
	F := ell.E.F
	t1 := F.Sqr(u)                  // 1.  t1 = u^2
	t1 = F.Mul(ell.Z, t1)           // 2.  t1 = Z * t1               // Z * u^2
	x1 := F.Add(t1, F.One())        // 3.  x1 = t1 + 1
	x1 = F.Inv0(x1)                 // 4.  x1 = inv0(x1)
	e1 := x1.IsZero()               // 5.  e1 = x1 == 0
	x1 = F.CMov(x1, F.One(), e1)    // 6.  x1 = CMOV(x1, 1, e1)      // if x1 == 0, set x1 = 1
	x1 = F.Mul(F.Neg(ell.E.A), x1)  // 7.  x1 = -A * x1            // x1 = -A / (1 + Z * u^2)
	gx1 := F.Add(x1, ell.E.A)       // 8. gx1 = x1 + A
	gx1 = F.Mul(gx1, x1)            // 9. gx1 = gx1 * x1
	gx1 = F.Add(gx1, ell.E.B)       //10. gx1 = gx1 + B
	gx1 = F.Mul(gx1, x1)            //11. gx1 = gx1 * x1           // gx1 = x1^3 + A * x1^2 + B * x1
	x2 := F.Neg(F.Add(x1, ell.E.A)) //12.  x2 = -x1 - A
	gx2 := F.Mul(t1, gx1)           //13. gx2 = t1 * gx1
	e2 := F.IsSquare(gx1)           //14.  e2 = is_square(gx1)
	x := F.CMov(x2, x1, e2)         //15.   x = CMOV(x2, x1, e2)   // If is_square(gx1), x = x1, else x = x2
	y2 := F.CMov(gx2, gx1, e2)      //16.  y2 = CMOV(gx2, gx1, e2) // If is_square(gx1), y2 = gx1, else y2 = gx2
	y := F.Sqrt(y2)                 //17.   y = sqrt(y2)
	e3 := F.Sgn0(u) == F.Sgn0(y)    //18.  e3 = sgn0(u) == sgn0(y) // fix sign of y
	y = F.CMov(F.Neg(y), y, e3)     //19.   y = CMOV(-y, y, e3)
	return MontgomeryPoint{x: x, y: y, z: F.One()}
}
func (ell elligator2) clearCofactor(P MontgomeryPoint) MontgomeryPoint {
	Q := P
	for i := uint(0); i < ell.E.Cofactor; i += 2 {
		Q = ell.E.Double(Q)
	}
	return Q
}
func (ell elligator2) Hash([]byte) Point {
	var u math.Element
	P := ell.MapToCurve(u)
	return ell.clearCofactor(P)
}

func (ell elligator2) HashRO([]byte) Point {
	var u0, u1 math.Element
	P0 := ell.MapToCurve(u0)
	P1 := ell.MapToCurve(u1)
	P := ell.E.Add(P0, P1)
	return ell.clearCofactor(P)
}
