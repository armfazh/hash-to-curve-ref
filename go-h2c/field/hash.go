package field

import (
	"hash"
	"io"
	"math/big"

	"golang.org/x/crypto/hkdf"
)

// HashToField is a function that hashes a string msg of any length into an
// element of a field fq.
func HashToField(
	H func() hash.Hash,
	msg, dst []byte,
	ctr byte,
	k uint,
	F Field) Elt {

	info := []byte{'H', '2', 'C', ctr, byte(1)}
	msgPrime := hkdf.Extract(H, append(msg, byte(0)), dst)

	m := F.Ext()
	e := make([]*big.Int, m)
	L := (uint(F.BitLen()) + k) / 8
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
