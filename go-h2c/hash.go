// +build ignore

package h2c

import (
	// "crypto/hmac"
	// "golang.org/x/crypto/hkdf"
	"hash"
	"io"
	"math/big"
	// "github.com/cloudflare/circl/math"
)

// Hash2Point implements a hash to point function
type Hash2Point interface {
	io.Writer
	Reset()
	Sum() (x, y *math.FqElt)
	IsRandomOracle() bool
}

type encoding struct {
	hkdf     hash.Hash
	mapping  func(u *math.EltFq) (x, y *math.EltFq)
	clearCof func(x, y *math.EltFq)
}

func (e encoding) Reset()                            { e.hkdf.Reset() }
func (e encoding) Write(p []byte) (n int, err error) { return e.hkdf.Write(p) }

type encodeToCurve struct {
	encoding
}

func (e encodeToCurve) IsRandomOracle() bool { return false }
func (e encodeToCurve) Sum() (x, y *big.Int) {
	m := e.hkdf.Sum(nil)
	u := e.hashTobase(m, ctr)
	x, y = e.mapping(u)
	e.clearCof(x, y)
	return x, y
}

type hashToCurve struct {
	encoding
}

func (h hashToCurve) IsRandomOracle() bool { return true }
func (h hashToCurve) Sum() (x, y *big.Int) {
	m := h.hkdf.Sum(nil)
	u := h.hashTobase(m, ctr)
	x, y = h.mapping(u)
	h.clearCof(x, y)
	return x, y
}
