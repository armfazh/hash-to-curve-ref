package h2c

import (
	"crypto/sha256"
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/svdw"
)

func suiteSECP2556K1() {
	iso := getIsogenySECP256K1()
	E := iso.Domain()
	// E := C.GetFromName("secp256k1-3iso")
	F := E.Field()
	h := sha256.New
	Z := F.Elt(-11)
	L := uint(128)
	sgn0 := GF.SignLE
	suites["secp256k1-SHA256-SSWU-NU-"] = EncodeToCurve{E, L, h, sswuAB0.New(E, Z, sgn0, iso)}
	suites["secp256k1-SHA256-SSWU-RO-"] = HashToCurve{E, L, h, sswuAB0.New(E, Z, sgn0, iso)}
	suites["secp256k1-SHA256-SVDW-NU-"] = EncodeToCurve{E, L, h, svdw.New(E, sgn0)}
	suites["secp256k1-SHA256-SVDW-RO-"] = HashToCurve{E, L, h, svdw.New(E, sgn0)}
}

type isosecp256k1 struct {
	E0, E1                 C.W
	xNum, xDen, yNum, yDen []GF.Elt
}

// getIsogenySECP256K1 returns a 3-isoeny
func getIsogenySECP256K1() C.Isogeny {
	e0 := C.GetFromName("secp256k1-3iso")
	e1 := C.GetFromName("secp256k1")
	F := e0.Field()
	return isosecp256k1{
		E0: e0.(C.W),
		E1: e1.(C.W),
		xNum: []GF.Elt{
			F.Elt("0x8e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38daaaaa8c7"),
			F.Elt("0x07d3d4c80bc321d5b9f315cea7fd44c5d595d2fc0bf63b92dfff1044f17c6581"),
			F.Elt("0x534c328d23f234e6e2a413deca25caece4506144037c40314ecbd0b53d9dd262"),
			F.Elt("0x8e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38daaaaa88c")},
		xDen: []GF.Elt{
			F.Elt("0xd35771193d94918a9ca34ccbb7b640dd86cd409542f8487d9fe6b745781eb49b"),
			F.Elt("0xedadc6f64383dc1df7c4b2d51b54225406d36b641f5e41bbc52a56612a8c6d14"),
			F.One()},
		yNum: []GF.Elt{
			F.Elt("0x4bda12f684bda12f684bda12f684bda12f684bda12f684bda12f684b8e38e23c"),
			F.Elt("0xc75e0c32d5cb7c0fa9d0a54b12a0a6d5647ab046d686da6fdffc90fc201d71a3"),
			F.Elt("0x29a6194691f91a73715209ef6512e576722830a201be2018a765e85a9ecee931"),
			F.Elt("0x2f684bda12f684bda12f684bda12f684bda12f684bda12f684bda12f38e38d84")},
		yDen: []GF.Elt{
			F.Elt("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffff93b"),
			F.Elt("0x7a06534bb8bdb49fd5e9e6632722c2989467c1bfc8e8d978dfb425d2685c2573"),
			F.Elt("0x6484aa716545ca2cf3a70c3fa8fe337e0a3d21162f0d6299a7bf8192bfd2a76f"),
			F.One()},
	}
}
func (m isosecp256k1) String() string       { return fmt.Sprintf("Isogeny from %v to\n%v", m.E0, m.E1) }
func (m isosecp256k1) Domain() C.EllCurve   { return m.E0 }
func (m isosecp256k1) Codomain() C.EllCurve { return m.E1 }
func (m isosecp256k1) Push(p C.Point) C.Point {
	F := m.E0.F
	x0, y0 := p.X(), p.Y()

	xNum := evalPoly(F, x0, m.xNum)
	xDen := evalPoly(F, x0, m.xDen)
	yNum := evalPoly(F, x0, m.yNum)
	yDen := evalPoly(F, x0, m.yDen)

	x1 := F.Mul(xNum, F.Inv(xDen))
	y1 := F.Mul(yNum, F.Inv(yDen))
	y1 = F.Mul(y1, y0)

	return m.E1.NewPoint(x1, y1)
}

// evalPoly evaluates a polynomial, given by its coefficients, on x. fx= sum a_ix^i
func evalPoly(f GF.Field, x GF.Elt, a []GF.Elt) GF.Elt {
	fx := f.Zero()
	for i := len(a) - 1; i >= 0; i-- {
		fx = f.Add(f.Mul(fx, x), a[i])
	}
	return fx
}
