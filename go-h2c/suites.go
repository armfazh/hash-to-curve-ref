package h2c

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/bf"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/elligator2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
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

func getCurve(name string) C.EllCurve {
	var E C.EllCurve
	switch name {
	case "curve25519":
	case "curve448":
	case "P256":
	case "P384":
	case "P521":
	default:
		panic("curve not recognized")
	}
	return E
}

func getHash(name string) func() hash.Hash {
	switch name {
	case "SHA256":
		return sha256.New
	case "SHA384":
		return sha512.New384
	case "SHA512":
		return sha512.New
	default:
		panic("not supported")
	}
}

func getMapping(curve C.EllCurve, name string) M.Map {
	var Z GF.Elt       // TBD
	var sgn0 GF.Sgn0ID // TBD
	switch name {
	case "SVDW":
		return svdw.New(curve, Z, sgn0)
	case "SSWU":
		return sswu.New(curve, Z, sgn0)
	case "SSWUAB0":
		return sswuAB0.New(curve, Z, sgn0, nil)
	case "ELL2":
		return elligator2.New(curve, Z, sgn0, nil)
	case "BF":
		return bf.New(curve)
	default:
		panic("mapping not supported")
	}
}

// GetSuite returns an implementation of hash to point
func GetSuite(suiteID string) HashToPoint {
	v := strings.Split(suiteID, "-")
	if len(v) != 5 {
		panic("wrong suiteID")
	}
	E := getCurve(v[0])
	h := getHash(v[1])
	m := getMapping(E, v[2])
	L := uint(128) //TBD
	switch v[3] {
	case "NU":
		return &EncodeToCurve{E, L, h, m}
	case "RO":
		return &HashToCurve{E, L, h, m}
	default:
		panic("wrong format")
	}
}
