package h2c

import (
	"crypto/sha256"
	"crypto/sha512"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/elligator2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/svdw"
)

var suiteNames = []string{
	"BLS12381G1-SHA256-SSWU-NU-",
	"BLS12381G1-SHA256-SSWU-RO-",
	"BLS12381G1-SHA256-SVDW-NU-",
	"BLS12381G1-SHA256-SVDW-RO-",
	"BLS12381G2-SHA256-SSWU-NU-",
	"BLS12381G2-SHA256-SSWU-RO-",
	"BLS12381G2-SHA256-SVDW-NU-",
	"BLS12381G2-SHA256-SVDW-RO-",
	"curve25519-SHA256-ELL2-NU-",
	"curve25519-SHA256-ELL2-RO-",
	"curve448-SHA512-ELL2-NU-",
	"curve448-SHA512-ELL2-RO-",
	"edwards25519-SHA256-EDELL2-NU-",
	"edwards25519-SHA256-EDELL2-RO-",
	"edwards448-SHA512-EDELL2-NU-",
	"edwards448-SHA512-EDELL2-RO-",
	"P256-SHA256-SSWU-NU-",
	"P256-SHA256-SSWU-RO-",
	"P256-SHA256-SVDW-NU-",
	"P256-SHA256-SVDW-RO-",
	"P384-SHA512-SSWU-NU-",
	"P384-SHA512-SSWU-RO-",
	"P384-SHA512-SVDW-NU-",
	"P384-SHA512-SVDW-RO-",
	"P521-SHA512-SSWU-NU-",
	"P521-SHA512-SSWU-RO-",
	"P521-SHA512-SVDW-NU-",
	"P521-SHA512-SVDW-RO-",
	"secp256k1-SHA256-SSWU-NU-",
	"secp256k1-SHA256-SSWU-RO-",
	"secp256k1-SHA256-SVDW-NU-",
	"secp256k1-SHA256-SVDW-RO-",
}

var suites map[string]HashToPoint

func init() {
	suites = make(map[string]HashToPoint)
	suitesWCurves()
	suitesMCurves()
	suiteSECP2556K1()
}

func suitesMCurves() {
	E := C.GetFromName("curve25519")
	h := sha256.New
	L := uint(128)
	sgn0 := GF.SignLE
	suites["curve25519-SHA256-ELL2-NU-"] = EncodeToCurve{E, L, h, elligator2.New(E, sgn0, nil)}
	suites["curve25519-SHA256-ELL2-RO-"] = HashToCurve{E, L, h, elligator2.New(E, sgn0, nil)}

	E = C.GetFromName("edwards25519")
	suites["edwards25519-SHA256-EDELL2-NU-"] = EncodeToCurve{E, L, h, elligator2.New(E, sgn0, nil)}
	suites["edwards25519-SHA256-EDELL2-RO-"] = HashToCurve{E, L, h, elligator2.New(E, sgn0, nil)}

	E = C.GetFromName("curve448")
	h = sha512.New
	L = uint(224)
	sgn0 = GF.SignLE
	suites["curve448-SHA512-ELL2-NU-"] = EncodeToCurve{E, L, h, elligator2.New(E, sgn0, nil)}
	suites["curve448-SHA512-ELL2-RO-"] = HashToCurve{E, L, h, elligator2.New(E, sgn0, nil)}

	E = C.GetFromName("edwards448")
	suites["edwards448-SHA512-EDELL2-NU-"] = EncodeToCurve{E, L, h, elligator2.New(E, sgn0, nil)}
	suites["edwards448-SHA512-EDELL2-RO-"] = HashToCurve{E, L, h, elligator2.New(E, sgn0, nil)}
}

func suitesWCurves() {
	E := C.GetFromName("P256")
	F := E.Field()
	h := sha256.New
	L := uint(128)
	Z := F.Elt(-10)
	sgn0 := GF.SignLE
	suites["P256-SHA256-SSWU-NU-"] = EncodeToCurve{E, L, h, sswu.New(E, Z, sgn0)}
	suites["P256-SHA256-SSWU-RO-"] = HashToCurve{E, L, h, sswu.New(E, Z, sgn0)}
	suites["P256-SHA256-SVDW-NU-"] = EncodeToCurve{E, L, h, svdw.New(E, sgn0)}
	suites["P256-SHA256-SVDW-RO-"] = HashToCurve{E, L, h, svdw.New(E, sgn0)}

	E = C.GetFromName("P384")
	F = E.Field()
	h = sha512.New384
	L = uint(192)
	Z = F.Elt(-12)
	sgn0 = GF.SignLE
	suites["P384-SHA384-SSWU-NU-"] = EncodeToCurve{E, L, h, sswu.New(E, Z, sgn0)}
	suites["P384-SHA384-SSWU-RO-"] = HashToCurve{E, L, h, sswu.New(E, Z, sgn0)}
	suites["P384-SHA384-SVDW-NU-"] = EncodeToCurve{E, L, h, svdw.New(E, sgn0)}
	suites["P384-SHA384-SVDW-RO-"] = HashToCurve{E, L, h, svdw.New(E, sgn0)}

	E = C.GetFromName("P521")
	F = E.Field()
	h = sha512.New
	L = uint(256)
	Z = F.Elt(-4)
	sgn0 = GF.SignLE
	suites["P521-SHA512-SSWU-NU-"] = EncodeToCurve{E, L, h, sswu.New(E, Z, sgn0)}
	suites["P521-SHA512-SSWU-RO-"] = HashToCurve{E, L, h, sswu.New(E, Z, sgn0)}
	suites["P521-SHA512-SVDW-NU-"] = EncodeToCurve{E, L, h, svdw.New(E, sgn0)}
	suites["P521-SHA512-SVDW-RO-"] = HashToCurve{E, L, h, svdw.New(E, sgn0)}
}
