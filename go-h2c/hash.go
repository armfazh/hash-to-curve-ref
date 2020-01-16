package h2c

import (
	"hash"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	M "github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

// HashToPoint is
type HashToPoint interface {
	// IsRandomOracle returns true if the output distribution is indifferentiable from a random oracle.
	IsRandomOracle() bool
	// Hash returns a point on an elliptic curve given an input string and a domin separation tag.
	Hash(in, dst []byte) C.Point
	// GetParams returns the params of the suite
	GetParams() *Params
}

// Params is
type Params struct {
	E       C.EllCurve
	L       uint
	HFunc   func() hash.Hash
	Mapping M.Map
}

func (p *Params) GetParams() *Params { return p }

// GetEncodeToCurve is a non-uniform encoding. This function encodes bit strings
// to points on an elliptic curve group (G). The distribution of the output is
// not uniformly random in G.
func GetEncodeToCurve(p *Params) HashToPoint { return encodeToCurve{p} }

// GetHashToCurve is a random oracle encoding from bit strings to points on an
// elliptic curve group (G). This function is suitable for applications
// requiring a random oracle in G.
func GetHashToCurve(p *Params) HashToPoint { return hashToCurve{p} }

type encodeToCurve struct{ *Params }

func (s encodeToCurve) IsRandomOracle() bool { return false }
func (s encodeToCurve) Hash(in, dst []byte) C.Point {
	u := GF.HashToField(in, dst, byte(2), s.HFunc, s.E.Field(), s.L)
	Q := s.Mapping.MapToCurve(u)
	P := s.E.ClearCofactor(Q)
	return P
}

type hashToCurve struct{ *Params }

func (s hashToCurve) IsRandomOracle() bool { return true }
func (s hashToCurve) Hash(in, dst []byte) C.Point {
	u0 := GF.HashToField(in, dst, byte(0), s.HFunc, s.E.Field(), s.L)
	u1 := GF.HashToField(in, dst, byte(1), s.HFunc, s.E.Field(), s.L)
	Q0 := s.Mapping.MapToCurve(u0)
	Q1 := s.Mapping.MapToCurve(u1)
	R := s.E.Add(Q0, Q1)
	P := s.E.ClearCofactor(R)
	return P
}
