package h2c

import (
	"crypto/sha256"
	"crypto/sha512"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/svdw"
)

// var suiteNames = []string{
// 	"BLS12381G1-SHA256-SSWU-NU-",
// 	"BLS12381G1-SHA256-SSWU-RO-",
// 	"BLS12381G1-SHA256-SVDW-NU-",
// 	"BLS12381G1-SHA256-SVDW-RO-",
// 	"BLS12381G2-SHA256-SSWU-NU-",
// 	"BLS12381G2-SHA256-SSWU-RO-",
// 	"BLS12381G2-SHA256-SVDW-NU-",
// 	"BLS12381G2-SHA256-SVDW-RO-",
// 	"curve25519-SHA256-ELL2-NU-",
// 	"curve25519-SHA256-ELL2-RO-",
// 	"curve448-SHA512-ELL2-NU-",
// 	"curve448-SHA512-ELL2-RO-",
// 	"edwards25519-SHA256-EDELL2-NU-",
// 	"edwards25519-SHA256-EDELL2-RO-",
// 	"edwards448-SHA512-EDELL2-NU-",
// 	"edwards448-SHA512-EDELL2-RO-",
// 	"P256-SHA256-SSWU-NU-",
// 	"P256-SHA256-SSWU-RO-",
// 	"P256-SHA256-SVDW-NU-",
// 	"P256-SHA256-SVDW-RO-",
// 	"P384-SHA512-SSWU-NU-",
// 	"P384-SHA512-SSWU-RO-",
// 	"P384-SHA512-SVDW-NU-",
// 	"P384-SHA512-SVDW-RO-",
// 	"P521-SHA512-SSWU-NU-",
// 	"P521-SHA512-SSWU-RO-",
// 	"P521-SHA512-SVDW-NU-",
// 	"P521-SHA512-SVDW-RO-",
// 	"secp256k1-SHA256-SSWU-NU-",
// 	"secp256k1-SHA256-SSWU-RO-",
// 	"secp256k1-SHA256-SVDW-NU-",
// 	"secp256k1-SHA256-SVDW-RO-",
// }

// Suites is a list of supported hash to curve suites
var Suites map[string]HashToPoint

func init() {
	Suites = make(map[string]HashToPoint)
	suitesWCurves()
	suitesMCurves()
	suiteSECP2556K1()
}

func suitesWCurves() {
	E := C.P256.Get()
	F := E.Field()
	h := sha256.New
	L := uint(48)
	Z := F.Elt(-10)
	sgn0 := GF.SignLE
	Suites["P256-SHA256-SSWU-NU-"] = GetEncodeToCurve(&Params{E, L, h, sswu.New(E, Z, sgn0)})
	Suites["P256-SHA256-SSWU-RO-"] = GetHashToCurve(&Params{E, L, h, sswu.New(E, Z, sgn0)})
	Suites["P256-SHA256-SVDW-NU-"] = GetEncodeToCurve(&Params{E, L, h, svdw.New(E, sgn0)})
	Suites["P256-SHA256-SVDW-RO-"] = GetHashToCurve(&Params{E, L, h, svdw.New(E, sgn0)})

	E = C.P384.Get()
	F = E.Field()
	h = sha512.New
	L = uint(72)
	Z = F.Elt(-12)
	sgn0 = GF.SignLE
	Suites["P384-SHA512-SSWU-NU-"] = GetEncodeToCurve(&Params{E, L, h, sswu.New(E, Z, sgn0)})
	Suites["P384-SHA512-SSWU-RO-"] = GetHashToCurve(&Params{E, L, h, sswu.New(E, Z, sgn0)})
	Suites["P384-SHA512-SVDW-NU-"] = GetEncodeToCurve(&Params{E, L, h, svdw.New(E, sgn0)})
	Suites["P384-SHA512-SVDW-RO-"] = GetHashToCurve(&Params{E, L, h, svdw.New(E, sgn0)})

	E = C.P521.Get()
	F = E.Field()
	h = sha512.New
	L = uint(96)
	Z = F.Elt(-4)
	sgn0 = GF.SignLE
	Suites["P521-SHA512-SSWU-NU-"] = GetEncodeToCurve(&Params{E, L, h, sswu.New(E, Z, sgn0)})
	Suites["P521-SHA512-SSWU-RO-"] = GetHashToCurve(&Params{E, L, h, sswu.New(E, Z, sgn0)})
	Suites["P521-SHA512-SVDW-NU-"] = GetEncodeToCurve(&Params{E, L, h, svdw.New(E, sgn0)})
	Suites["P521-SHA512-SVDW-RO-"] = GetHashToCurve(&Params{E, L, h, svdw.New(E, sgn0)})
}
