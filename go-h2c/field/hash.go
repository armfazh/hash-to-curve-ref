package field

import (
	"hash"
	"io"
	"math/big"

	"golang.org/x/crypto/hkdf"
)

// HashToField is a function that hashes a string msg of any length into an
// element of a field fq.
//
// Parameters:
// - msg is the message to hash.
// - DST, a domain separation tag (see discussion above).
// - ctr is 0, 1, or 2.
// - H, a cryptographic hash function.
// - F, a finite field of characteristic p and order q = p^m.
// - L = ceil((ceil(log2(p)) + k) / 8), where k is the security parameter of
// the cryptosystem (e.g., k = 128).
// - HKDF as defined in RFC-5869 and instantiated with H.
func HashToField(
	msg, dst []byte,
	ctr byte,
	H func() hash.Hash,
	F Field,
	L uint) Elt {

	info := []byte{'H', '2', 'C', ctr, byte(1)}
	msgPrime := hkdf.Extract(H, append(msg, byte(0)), dst)

	m := F.Ext()
	e := make([]*big.Int, m)
	t := make([]byte, L)

	for i := uint(1); i <= m; i++ {
		info[4] = byte(i)
		rd := hkdf.Expand(H, msgPrime, info)
		if _, err := io.ReadFull(rd, t); err != nil {
			panic("error on hdkf")
		}
		e[i-1] = new(big.Int).SetBytes(t)
		e[i-1].Mod(e[i-1], F.P())
	}
	return F.EltFromList(e)
}
