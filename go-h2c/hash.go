package h2c

import (
	"hash"

	C "github.com/armfazh/hash-to-curve-ref/go-h2c/curve"
	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
	"github.com/armfazh/hash-to-curve-ref/go-h2c/mapping"
)

// HashToPoint is
type HashToPoint interface {
	IsRandomOracle() bool
	Hash(in, dst []byte) C.Point
}

// Suite is
type Suite struct {
	E     C.EllCurve
	L     uint
	HFunc func() hash.Hash
	mapping.Map
}

type EncodeToCurve struct{ *Suite }

func (s *EncodeToCurve) IsRandomOracle() bool { return false }
func (s *EncodeToCurve) Hash(in, dst []byte) C.Point {
	u := GF.HashToField(in, dst, byte(2), s.HFunc, s.E.Field(), s.L)
	Q := s.MapToCurve(u)
	P := s.E.ClearCofactor(Q)
	return P
}

type HashToCurve struct{ *Suite }

func (s *HashToCurve) IsRandomOracle() bool { return true }
func (s *HashToCurve) Hash(in, dst []byte) C.Point {
	u0 := GF.HashToField(in, dst, byte(0), s.HFunc, s.E.Field(), s.L)
	u1 := GF.HashToField(in, dst, byte(1), s.HFunc, s.E.Field(), s.L)
	Q0 := s.MapToCurve(u0)
	Q1 := s.MapToCurve(u1)
	R := s.E.Add(Q0, Q1)
	P := s.E.ClearCofactor(R)
	return P
}
