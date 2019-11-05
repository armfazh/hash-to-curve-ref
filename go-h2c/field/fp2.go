package field

import (
	"crypto/rand"
	"io"
	"math/big"
)

type fp2Elt struct {
	a, b *big.Int
}

func (e fp2Elt) String() string {
	return "\na: 0x" + e.a.Text(16) +
		"\nb: 0x" + e.b.Text(16) + " * i"
}

// IsZero returns true is x is zero,
func (e fp2Elt) IsZero() bool { return e.a.Sign() == 0 && e.b.Sign() == 0 }

type fp2 struct{ prime }

func newFp2(mod prime) Field { return fp2{mod} }

func (f fp2) P() *big.Int { return f.p }
func (f fp2) Elt() Elt    { return fp2Elt{big.NewInt(0), big.NewInt(0)} }
func (f fp2) Rand(r io.Reader) Elt {
	a, _ := rand.Int(r, f.p)
	b, _ := rand.Int(r, f.p)
	return fp2Elt{a, b}
}
func (f fp2) String() string { return "GF(" + f.name + ") Irred: i^2+1" }
func (f fp2) Ext() uint      { return uint(2) }
func (f fp2) One() Elt       { return fp2Elt{big.NewInt(1), big.NewInt(0)} }
func (f fp2) Zero() Elt      { return f.Elt() }
func (f fp2) BitLen() int    { return f.p.BitLen() }
func (f fp2) EltFromList(in []*big.Int) Elt {
	if len(in) != 2 {
		panic("wrong length")
	}
	return fp2Elt{in[0], in[1]}
}
func (f fp2) Add(x, y Elt) Elt {
	a := new(big.Int).Add(x.(fp2Elt).a, y.(fp2Elt).a)
	b := new(big.Int).Add(x.(fp2Elt).b, y.(fp2Elt).b)
	a.Mod(a, f.p)
	b.Mod(b, f.p)
	return fp2Elt{a, b}
}
func (f fp2) Sub(x, y Elt) Elt {
	a := new(big.Int).Sub(x.(fp2Elt).a, y.(fp2Elt).a)
	b := new(big.Int).Sub(x.(fp2Elt).b, y.(fp2Elt).b)
	a.Mod(a, f.p)
	b.Mod(b, f.p)
	return fp2Elt{a, b}
}

func (f fp2) Mul(x, y Elt) Elt          { return nil }
func (f fp2) Sqr(x Elt) Elt             { return nil }
func (f fp2) Inv0(x Elt) Elt            { return nil }
func (f fp2) Neg(x Elt) Elt             { return nil }
func (f fp2) Sgn0(x Elt) int            { return 0 }
func (f fp2) CMov(x, y Elt, b bool) Elt { return nil }
func (f fp2) Sqrt(x Elt) Elt            { return nil }
func (f fp2) IsSquare(x Elt) bool       { return false }
