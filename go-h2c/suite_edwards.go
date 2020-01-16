package h2c

import (
	"crypto/sha256"
	"crypto/sha512"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/elligator2"
)

func suitesMCurves() {
	E := C.Curve25519.Get()
	h := sha256.New
	L := uint(48)
	sgn0 := GF.SignLE
	Suites["curve25519-SHA256-ELL2-NU-"] = GetEncodeToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
	Suites["curve25519-SHA256-ELL2-RO-"] = GetHashToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})

	E = C.Edwards25519.Get()
	Suites["edwards25519-SHA256-EDELL2-NU-"] = GetEncodeToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
	// Suites["edwards25519-SHA256-EDELL2-RO-"] = GetHashToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})

	E = C.Curve448.Get()
	h = sha512.New
	L = uint(84)
	sgn0 = GF.SignLE
	Suites["curve448-SHA512-ELL2-NU-"] = GetEncodeToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
	Suites["curve448-SHA512-ELL2-RO-"] = GetHashToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})

	// E = C.Edwards448.Get()
	// Suites["edwards448-SHA512-EDELL2-NU-"] = GetEncodeToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
	// Suites["edwards448-SHA512-EDELL2-RO-"] = GetHashToCurve(&Params{E, L, h, elligator2.New(E, sgn0)})
}
