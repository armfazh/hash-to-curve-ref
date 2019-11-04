package h2c

import (
	"golang.org/x/crypto/hkdf"
	"hash"
	"io"

	"github.com/armfazh/hash-to-curve-ref/h2c/math"
)

// HashToField is a function that hashes a string msg of any length into an
// element of a field fq.
func HashToField(
	H func() hash.Hash,
	msg, dst []byte,
	ctr, k uint,
	fq math.Fq) math.EltFq {

	const prefix = "H2C"
	L := (fq.P.Size() + k) / 8
	t := make([]byte, L)
	u := make([]math.Elt, fq.M)
	mp := hkdf.Extract(H, msg, dst)
	ctrOctet := []byte{byte(ctr)}
	for i := uint(1); i <= fq.M; i++ {
		info := append(append([]byte(prefix), ctrOctet...), []byte{byte(i)}...)
		rd := hkdf.Expand(H, mp, info)
		if _, err := io.ReadFull(rd, t); err != nil {
			panic("error on hdkf")
		}
		u[i-1] = fq.P.EltFromBytes(t)
	}
	return u
}
