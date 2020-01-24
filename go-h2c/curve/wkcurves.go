package curve

import (
	"math/big"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type CurveID int

const (
	Custom CurveID = iota
	P256
	P384
	P521
	Curve25519
	Curve448
	Edwards25519
	Edwards448
	SECP256K1
	SECP256K1_3ISO
)

// GetFromID is
func (id CurveID) Get() EllCurve {
	switch id {
	case P256:
		f := GF.P256.Get()
		return NewWeierstrass(id, f,
			f.Elt("-3"),
			f.Elt("0x5ac635d8aa3a93e7b3ebbd55769886bc651d06b0cc53b0f63bce3c3e27d2604b"),
			GF.FromType("0xffffffff00000000ffffffffffffffffbce6faada7179e84f3b9cac2fc632551"),
			big.NewInt(1))
	case P384:
		f := GF.P384.Get()
		return NewWeierstrass(id, f,
			f.Elt("-3"),
			f.Elt("0xb3312fa7e23ee7e4988e056be3f82d19181d9c6efe8141120314088f5013875ac656398d8a2ed19d2a85c8edd3ec2aef"),
			GF.FromType("0xffffffffffffffffffffffffffffffffffffffffffffffffc7634d81f4372ddf581a0db248b0a77aecec196accc52973"),
			big.NewInt(1))
	case P521:
		f := GF.P521.Get()
		return NewWeierstrass(id, f,
			f.Elt("-3"),
			f.Elt("0x051953eb9618e1c9a1f929a21a0b68540eea2da725b99b315f3b8b489918ef109e156193951ec7e937b1652c0bd3bb1bf073573df883d2c34f1ef451fd46b503f00"),
			GF.FromType("0x7ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffd15b6c64746fc85f736b8af5e7ec53f04fbd8c4569a8f1f4540ea2435f5180d6b"),
			big.NewInt(1))
	case SECP256K1:
		f := GF.P256K1.Get()
		return NewWeierstrass(id, f,
			f.Zero(),
			f.Elt("7"),
			GF.FromType("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"),
			big.NewInt(1))
	case SECP256K1_3ISO:
		f := GF.P256K1.Get()
		return NewWeierstrass(id, f,
			f.Elt("0x3f8731abdd661adca08a5558f0f5d272e953d363cb6f0e5d405447c01a444533"),
			f.Elt("1771"),
			GF.FromType("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"),
			big.NewInt(1))
	case Curve25519:
		f := GF.P25519.Get()
		return NewMontgomery(id, f,
			f.Elt("486662"),
			f.One(),
			GF.FromType("0x1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed"),
			big.NewInt(8))
	case Edwards25519:
		f := GF.P25519.Get()
		return NewEdwards(id, f,
			f.Elt("-1"),
			f.Elt("0x52036cee2b6ffe738cc740797779e89800700a4d4141d8ab75eb4dca135978a3"),
			GF.FromType("0x1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed"),
			big.NewInt(8))
	case Curve448:
		f := GF.P448.Get()
		return NewMontgomery(id, f,
			f.Elt("156326"),
			f.One(),
			GF.FromType("0x3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3"),
			big.NewInt(4))
	case Edwards448:
		f := GF.P448.Get()
		return NewEdwards(id, f,
			f.One(),
			f.Elt("-39081"),
			GF.FromType("0x3fffffffffffffffffffffffffffffffffffffffffffffffffffffff7cca23e9c44edb49aed63690216cc2728dc58f552378c292ab5844f3"),
			big.NewInt(4))
	default:
		panic("curve not supported")
	}
}
