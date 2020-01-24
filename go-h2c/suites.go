package h2c

import (
	"crypto"
	_ "crypto/sha256" // To link the sha256 module
	_ "crypto/sha512" // To link the sha256 module
	"fmt"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

// SuiteID is the identifier of supported hash to curve suites.
type SuiteID uint16

const (
	P256_SHA256_SSWU_NU_ SuiteID = iota
	P256_SHA256_SSWU_RO_
	P256_SHA256_SVDW_NU_
	P256_SHA256_SVDW_RO_
	P384_SHA512_SSWU_NU_
	P384_SHA512_SSWU_RO_
	P384_SHA512_SVDW_NU_
	P384_SHA512_SVDW_RO_
	P521_SHA512_SSWU_NU_
	P521_SHA512_SSWU_RO_
	P521_SHA512_SVDW_NU_
	P521_SHA512_SVDW_RO_
	Curve25519_SHA256_ELL2_NU_
	Curve25519_SHA256_ELL2_RO_
	Edwards25519_SHA256_EDELL2_NU_
	Edwards25519_SHA256_EDELL2_RO_
	Curve448_SHA512_ELL2_NU_
	Curve448_SHA512_ELL2_RO_
	Edwards448_SHA512_EDELL2_NU_
	Edwards448_SHA512_EDELL2_RO_
	SECP256k1_SHA256_SSWU_NU_
	SECP256k1_SHA256_SSWU_RO_
	SECP256k1_SHA256_SVDW_NU_
	SECP256k1_SHA256_SVDW_RO_
)

// Get returns a HashToPoint based on the SuiteID, otherwise returns an error
// if the SuiteID is not supported or invalid.
func (id SuiteID) Get() (HashToPoint, error) {
	if s, ok := supportedSuitesID[id]; ok {
		return s.New(), nil
	}
	return nil, fmt.Errorf("Suite: %v not supported", id)
}
func (id SuiteID) register(name string, s *params) {
	s.ID = id
	supportedSuitesNames[name] = id
	supportedSuitesID[id] = *s
}

// GetSuite returns a HashToPoint function from a named suite.
func GetSuite(name string) (HashToPoint, error) {
	if s, ok := supportedSuitesNames[name]; ok {
		return s.Get()
	}
	return nil, fmt.Errorf("Suite: %v not supported", name)
}

var supportedSuitesNames map[string]SuiteID
var supportedSuitesID map[SuiteID]params

func init() {
	supportedSuitesNames = make(map[string]SuiteID)
	supportedSuitesID = make(map[SuiteID]params)

	P256_SHA256_SSWU_NU_.register("P256-SHA256-SSWU-NU-", &params{Ell: C.P256, H: crypto.SHA256, Map: M.SSWU, Sgn0: GF.SignLE, L: 48, RO: false, Z: -10})
	P256_SHA256_SSWU_RO_.register("P256-SHA256-SSWU-RO-", &params{Ell: C.P256, H: crypto.SHA256, Map: M.SSWU, Sgn0: GF.SignLE, L: 48, RO: true, Z: -10})
	P256_SHA256_SVDW_NU_.register("P256-SHA256-SVDW-NU-", &params{Ell: C.P256, H: crypto.SHA256, Map: M.SVDW, Sgn0: GF.SignLE, L: 48, RO: false})
	P256_SHA256_SVDW_RO_.register("P256-SHA256-SVDW-RO-", &params{Ell: C.P256, H: crypto.SHA256, Map: M.SVDW, Sgn0: GF.SignLE, L: 48, RO: true})

	P384_SHA512_SSWU_NU_.register("P384-SHA512-SSWU-NU-", &params{Ell: C.P384, H: crypto.SHA512, Map: M.SSWU, Sgn0: GF.SignLE, L: 72, RO: false, Z: -12})
	P384_SHA512_SSWU_RO_.register("P384-SHA512-SSWU-RO-", &params{Ell: C.P384, H: crypto.SHA512, Map: M.SSWU, Sgn0: GF.SignLE, L: 72, RO: true, Z: -12})
	P384_SHA512_SVDW_NU_.register("P384-SHA512-SVDW-NU-", &params{Ell: C.P384, H: crypto.SHA512, Map: M.SVDW, Sgn0: GF.SignLE, L: 72, RO: false})
	P384_SHA512_SVDW_RO_.register("P384-SHA512-SVDW-RO-", &params{Ell: C.P384, H: crypto.SHA512, Map: M.SVDW, Sgn0: GF.SignLE, L: 72, RO: true})

	P521_SHA512_SSWU_NU_.register("P521-SHA512-SSWU-NU-", &params{Ell: C.P521, H: crypto.SHA512, Map: M.SSWU, Sgn0: GF.SignLE, L: 96, RO: false, Z: -4})
	P521_SHA512_SSWU_RO_.register("P521-SHA512-SSWU-RO-", &params{Ell: C.P521, H: crypto.SHA512, Map: M.SSWU, Sgn0: GF.SignLE, L: 96, RO: true, Z: -4})
	P521_SHA512_SVDW_NU_.register("P521-SHA512-SVDW-NU-", &params{Ell: C.P521, H: crypto.SHA512, Map: M.SVDW, Sgn0: GF.SignLE, L: 96, RO: false})
	P521_SHA512_SVDW_RO_.register("P521-SHA512-SVDW-RO-", &params{Ell: C.P521, H: crypto.SHA512, Map: M.SVDW, Sgn0: GF.SignLE, L: 96, RO: true})

	Curve25519_SHA256_ELL2_NU_.register("curve25519-SHA256-ELL2-NU-", &params{Ell: C.Curve25519, H: crypto.SHA256, Map: M.ELL2, Sgn0: GF.SignLE, L: 48, RO: false})
	Curve25519_SHA256_ELL2_RO_.register("curve25519-SHA256-ELL2-RO-", &params{Ell: C.Curve25519, H: crypto.SHA256, Map: M.ELL2, Sgn0: GF.SignLE, L: 48, RO: true})
	Edwards25519_SHA256_EDELL2_NU_.register("edwards25519-SHA256-EDELL2-NU-", &params{Ell: C.Edwards25519, H: crypto.SHA256, Map: M.EDELL2, Sgn0: GF.SignLE, L: 48, RO: false})
	Edwards25519_SHA256_EDELL2_RO_.register("edwards25519-SHA256-EDELL2-RO-", &params{Ell: C.Edwards25519, H: crypto.SHA256, Map: M.EDELL2, Sgn0: GF.SignLE, L: 48, RO: true})

	Curve448_SHA512_ELL2_NU_.register("curve448-SHA512-ELL2-NU-", &params{Ell: C.Curve448, H: crypto.SHA512, Map: M.ELL2, Sgn0: GF.SignLE, L: 84, RO: false})
	Curve448_SHA512_ELL2_RO_.register("curve448-SHA512-ELL2-RO-", &params{Ell: C.Curve448, H: crypto.SHA512, Map: M.ELL2, Sgn0: GF.SignLE, L: 84, RO: true})
	Edwards448_SHA512_EDELL2_NU_.register("edwards448-SHA512-EDELL2-NU-", &params{Ell: C.Edwards448, H: crypto.SHA512, Map: M.EDELL2, Sgn0: GF.SignLE, L: 84, RO: false})
	Edwards448_SHA512_EDELL2_RO_.register("edwards448-SHA512-EDELL2-RO-", &params{Ell: C.Edwards448, H: crypto.SHA512, Map: M.EDELL2, Sgn0: GF.SignLE, L: 84, RO: true})

	SECP256k1_SHA256_SSWU_NU_.register("secp256k1-SHA256-SSWU-NU-", &params{Ell: C.SECP256K1, H: crypto.SHA256, Map: M.SSWU, Sgn0: GF.SignLE, L: 48, RO: false, Z: -11, Iso: C.GetSECP256K1Isogeny()})
	SECP256k1_SHA256_SSWU_RO_.register("secp256k1-SHA256-SSWU-RO-", &params{Ell: C.SECP256K1, H: crypto.SHA256, Map: M.SSWU, Sgn0: GF.SignLE, L: 48, RO: true, Z: -11, Iso: C.GetSECP256K1Isogeny()})
	SECP256k1_SHA256_SVDW_NU_.register("secp256k1-SHA256-SVDW-NU-", &params{Ell: C.SECP256K1, H: crypto.SHA256, Map: M.SVDW, Sgn0: GF.SignLE, L: 48, RO: false})
	SECP256k1_SHA256_SVDW_RO_.register("secp256k1-SHA256-SVDW-RO-", &params{Ell: C.SECP256K1, H: crypto.SHA256, Map: M.SVDW, Sgn0: GF.SignLE, L: 48, RO: true})
}

type params struct {
	ID   SuiteID
	Ell  C.CurveID
	H    crypto.Hash
	Map  M.ID
	Sgn0 GF.Sgn0ID
	L    uint
	Z    int
	Iso  C.Isogeny
	RO   bool
}

// New returns a HashToPoint that encodes bit strings to points on an elliptic
// curve group.
func (s *params) New() HashToPoint {
	E := s.Ell.Get()
	H := s.H.New
	Z := E.Field().Elt(s.Z)
	m := s.Map.Get(E, Z, s.Sgn0, s.Iso)
	e := &encoding{E, H, s.L, m, s.RO}
	if s.RO {
		return &hashToCurve{e}
	}
	return &encodeToCurve{e}
}
