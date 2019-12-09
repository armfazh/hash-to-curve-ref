package suite

import (
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"

	"github.com/armfazh/hash-to-curve-ref/go-h2c"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/bf"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/elligator2"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswu"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/sswuAB0"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping/svdw"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/toy"
)

var ToySuites map[string]h2c.HashToPoint

func init() {
	initSuites()
}

func initSuites() {
	ToySuites = make(map[string]h2c.HashToPoint)
	// TODO
	// RegisterToySuite("W0-SHA256-SSWU-NU-")
	// RegisterToySuite("W0-SHA256-SSWU-RO-")
}

func RegisterToySuite(suiteID string) {
	v := strings.Split(suiteID, "-")
	if len(v) != 5 {
		panic("wrong suiteID")
	}
	curveID := v[0]
	ecc := toy.ToyCurves[curveID]

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
	var Z GF.Elt       // TBD
	var sgn0 GF.Sgn0ID // TBD
	var l uint = 16    // TBD
	var mm mapping.Map
	switch v[2] {
	case "SVDW":
		mm = svdw.New(ecc.E, Z, sgn0)
	case "SSWU":
		mm = sswu.New(ecc.E, Z, sgn0)
	case "SSWUAB0":
		mm = sswuAB0.New(ecc.E, Z, sgn0, nil)
	case "ELL2":
		mm = elligator2.New(ecc.E, Z, sgn0, nil)
	case "BF":
		mm = bf.New(ecc.E)
	}
	s := &h2c.Suite{E: ecc.E, L: l, HFunc: h, Map: mm}
	var h2p h2c.HashToPoint

	switch v[3] {
	case "NU":
		h2p = &h2c.EncodeToCurve{Suite: s}
	case "RO":
		h2p = &h2c.HashToCurve{Suite: s}
	default:
		panic("wrong suiteID")
	}
	ToySuites[suiteID] = h2p
}
