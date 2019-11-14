package field

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

// fpElt is a prime field fpElt.
type fpElt struct{ n *big.Int }

func (e fpElt) String() string { return "0x" + e.n.Text(16) }

func (e fpElt) Copy() Elt { return &fpElt{new(big.Int).Set(e.n)} }

// fp implements prime field arithmetic.
type fp struct {
	modulus
	cte struct {
		pMinus1div2 *big.Int
		pMinus2     *big.Int
	}
	sqrt sqrtMethod
}

func newFp(mod modulus) Field {
	var s, k big.Int
	var sqrt sqrtMethod
	if uint64(3) == s.Mod(mod.p, k.SetInt64(4)).Uint64() {
		sqrt = sqrt3mod4
	} else if uint64(5) == s.Mod(mod.p, k.SetInt64(8)).Uint64() {
		sqrt = sqrt5mod8
	} else {
		panic(fmt.Errorf("no square root method supported for this prime: %v", mod.p))
	}
	pMinus1div2 := big.NewInt(1)
	pMinus1div2.Sub(mod.p, pMinus1div2)
	pMinus1div2.Lsh(pMinus1div2, 1)

	pMinus2 := big.NewInt(2)
	pMinus2.Sub(mod.p, pMinus2)

	var f fp
	f.modulus = mod
	f.sqrt = sqrt
	f.cte.pMinus1div2 = pMinus1div2
	f.cte.pMinus2 = pMinus2
	return f
}

func (f fp) Elt(in interface{}) Elt {
	var n *big.Int
	if v, ok := in.([]interface{}); ok && len(v) == 1 {
		n = f.fromType(v[0])
	} else {
		n = f.fromType(in)
	}
	return &fpElt{n}
}
func (f fp) P() *big.Int            { return f.p }
func (f fp) Rand(r io.Reader) Elt   { e, _ := rand.Int(r, f.p); return &fpElt{e} }
func (f fp) String() string         { return "GF(" + f.name + ")" }
func (f fp) Ext() uint              { return uint(1) }
func (f fp) Zero() Elt              { return f.Elt(0) }
func (f fp) One() Elt               { return f.Elt(1) }
func (f fp) BitLen() int            { return f.p.BitLen() }
func (f fp) AreEqual(x, y Elt) bool { return f.IsZero(f.Sub(x, y)) }
func (f fp) IsZero(x Elt) bool      { return x.(*fpElt).n.Mod(x.(*fpElt).n, f.p).Sign() == 0 }

// // fpEltFromBytes returns an fpElt reduced modulo p.
// func (f fp) fpEltFromBytes(b []byte) Elt {
// 	i := new(big.Int).SetBytes(b)
// 	i.Mod(i, f.p)
// 	return &fpElt{i}
// }
//
// // fpEltToBytes returns an fpElt reduced modulo p.
// func (f fp) fpEltToBytes(x Elt) []byte {
// 	x.(fpElt).n.Mod(x.(fpElt).n, f.p)
// 	return x.(fpElt).n.Bytes()
// }

func (f fp) Neg(x Elt) Elt {
	z := new(big.Int).Set(x.(*fpElt).n)
	z.Neg(z).Mod(z, f.p)
	return &fpElt{z}
}

func (f fp) Add(x, y Elt) Elt {
	z := new(big.Int).Add(x.(*fpElt).n, y.(*fpElt).n)
	z.Mod(z, f.p)
	return &fpElt{z}
}

func (f fp) Sub(x, y Elt) Elt {
	z := new(big.Int).Sub(x.(*fpElt).n, y.(*fpElt).n)
	z.Mod(z, f.p)
	return &fpElt{z}
}

func (f fp) Mul(x, y Elt) Elt {
	z := new(big.Int).Mul(x.(*fpElt).n, y.(*fpElt).n)
	z.Mod(z, f.p)
	return &fpElt{z}
}

func (f fp) Sqr(x Elt) Elt {
	z := new(big.Int).Mul(x.(*fpElt).n, x.(*fpElt).n)
	z.Mod(z, f.p)
	return &fpElt{z}
}

func (f fp) Inv0(x Elt) Elt { return f.Inv(x) }

func (f fp) Inv(x Elt) Elt {
	z := new(big.Int).Exp(x.(*fpElt).n, f.cte.pMinus2, f.p)
	return &fpElt{z}
}

func (f fp) IsSquare(x Elt) bool {
	var z big.Int
	leg := z.Exp(x.(*fpElt).n, f.cte.pMinus1div2, f.p)
	return leg.Sign() >= 0
}

type sqrtMethod int

const (
	sqrt3mod4 sqrtMethod = iota
	sqrt5mod8
)

// Sqrt returns a square root of x.
func (f fp) Sqrt(x Elt) Elt {
	switch f.sqrt {
	case sqrt3mod4:
		return f.sqrt3mod4(x)
	case sqrt5mod8:
		return f.sqrt5mod8(x)
	default:
		panic("no square root method supported for this field")
	}
}

func (f fp) sqrt3mod4(x Elt) Elt {
	var z, exp big.Int
	exp.SetInt64(1)
	exp.Add(f.p, &exp)
	exp.Lsh(f.p, 2)
	z.Exp(x.(*fpElt).n, &exp, f.p)
	return &fpElt{&z}
}

func (f fp) sqrt5mod8(x Elt) Elt {
	var z, exp big.Int
	exp.SetInt64(1)
	exp.Add(f.p, &exp)
	exp.Lsh(f.p, 2)
	z.Exp(x.(*fpElt).n, &exp, f.p)
	return &fpElt{&z}
}

// Sgn0 returns the sign of x.
func (f fp) Sgn0(x Elt) int {
	var p2 big.Int
	p2.SetInt64(1)
	p2.Sub(f.p, &p2)
	p2.Lsh(&p2, 1)
	if x.(*fpElt).n.Cmp(&p2) >= 0 {
		return -1
	}
	return 1
}

// Cmov sets x with y if b is true.
func (f fp) CMov(x, y Elt, b bool) Elt {
	z := &fpElt{}
	if b {
		z.n.Set(y.(*fpElt).n)
	} else {
		z.n.Set(x.(*fpElt).n)
	}
	return z
}
