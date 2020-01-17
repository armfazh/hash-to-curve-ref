package h2c

import (
	"crypto/sha256"
	"crypto/sha512"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/elligator2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/svdw"
)

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

func suitesMCurves() {
	E := C.Curve25519.Get()
	h := sha256.New
	L := uint(48)
	sgn0 := GF.SignLE
	Suites["curve25519-SHA256-ELL2-NU-"] = GetEncodeToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
	Suites["curve25519-SHA256-ELL2-RO-"] = GetHashToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})

	E = C.Edwards25519.Get()
	Suites["edwards25519-SHA256-EDELL2-NU-"] = GetEncodeToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
	Suites["edwards25519-SHA256-EDELL2-RO-"] = GetHashToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})

	E = C.Curve448.Get()
	h = sha512.New
	L = uint(84)
	sgn0 = GF.SignLE
	Suites["curve448-SHA512-ELL2-NU-"] = GetEncodeToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
	Suites["curve448-SHA512-ELL2-RO-"] = GetHashToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})

	E = C.Edwards448.Get()
	Suites["edwards448-SHA512-EDELL2-NU-"] = GetEncodeToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
	Suites["edwards448-SHA512-EDELL2-RO-"] = GetHashToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
}

func suiteSECP2556K1() {
	iso := C.GetSECP256K1Isogeny()
	E0 := iso.Domain()
	E1 := iso.Codomain()
	h := sha256.New
	F := E0.Field()
	Z := F.Elt(-11)
	L := uint(48)
	sgn0 := GF.SignLE
	Suites["secp256k1-SHA256-SSWU-NU-"] = GetEncodeToCurve(&Params{E1, L, h, sswuAB0.New(E0, Z, sgn0, iso)})
	Suites["secp256k1-SHA256-SSWU-RO-"] = GetHashToCurve(&Params{E1, L, h, sswuAB0.New(E0, Z, sgn0, iso)})
	Suites["secp256k1-SHA256-SVDW-NU-"] = GetEncodeToCurve(&Params{E1, L, h, svdw.New(E1, sgn0)})
	Suites["secp256k1-SHA256-SVDW-RO-"] = GetHashToCurve(&Params{E1, L, h, svdw.New(E1, sgn0)})
}
