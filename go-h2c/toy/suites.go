package toy

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"

	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/ell2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/svdw"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/suite"
)

var ToySuites map[string]suite.HashToPoint

func init() {
	initCurves()
	initSuites()
}

func initSuites() {
	ToySuites = make(map[string]suite.HashToPoint)
	RegisterToySuite("W0-SHA256-SSWU-NU-")
	RegisterToySuite("W0-SHA256-SSWU-RO-")
}

func RegisterToySuite(suiteID string) {
	v := strings.Split(suiteID, "-")
	if len(v) != 5 {
		panic("wrong suiteID")
	}
	curveID := v[0]
	ecc := ToyCurves[curveID]

	var h func() hash.Hash
	switch v[1] {
	case "SHA256":
		h = sha256.New
	case "SHA384":
		h = sha512.New384
	case "SHA512":
		h = sha512.New
	default:
		panic("not supported")
	}
	var mm mapping.Map
	switch v[2] {
	case "SVDW":
		mm = svdw.New(ecc.E, ecc.Z, ecc.Sgn0)
	case "SSWU":
		mm = sswu.New(ecc.E, ecc.Z, ecc.Sgn0)
	case "SSWUAB0":
		mm = sswuAB0.New(ecc.E, ecc.Z, ecc.Sgn0, ecc.Iso)
	case "ELL2":
		mm = ell2.New(ecc.E, ecc.Z, ecc.Sgn0)
	}
	l := uint(16)
	s := &suite.Suite{E: ecc.E, L: l, HFunc: h, Map: mm}
	var h2p suite.HashToPoint

	switch v[3] {
	case "NU":
		h2p = &suite.EncodeToCurve{Suite: s}
	case "RO":
		h2p = &suite.HashToCurve{Suite: s}
	default:
		panic("wrong suiteID")
	}
	ToySuites[suiteID] = h2p
}
