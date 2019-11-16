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

func (e fp2Elt) Copy() Elt { r := &fp2Elt{}; r.a.Set(e.a); r.b.Set(e.b); return r }

type fp2 struct {
	p    *big.Int
	name string
}

// NewFp2 is
func NewFp2(name string, p interface{}) Field {
	prime := fromType(p)
	if !prime.ProbablyPrime(4) {
		panic("p is not prime")
	}
	return fp2{p: prime, name: name}
}

func (f fp2) Elt(in interface{}) Elt {
	var a, b *big.Int
	if v, ok := in.([]interface{}); ok && len(v) == 2 {
		a = fromType(v[0])
		b = fromType(v[1])
	} else {
		a = fromType(in)
		b = big.NewInt(0)
	}
	return f.mod(a, b)
}
func (f fp2) P() *big.Int    { return f.p }
func (f fp2) String() string { return "GF(" + f.name + ") Irred: i^2+1" }
func (f fp2) Ext() uint      { return uint(2) }
func (f fp2) Zero() Elt      { return f.Elt(0) }
func (f fp2) One() Elt       { return f.Elt(1) }
func (f fp2) BitLen() int    { return f.p.BitLen() }

func (f fp2) AreEqual(x, y Elt) bool { return f.IsZero(f.Sub(x, y)) }
func (f fp2) IsZero(x Elt) bool {
	e := x.(*fp2Elt)
	return e.a.Mod(e.a, f.p).Sign() == 0 &&
		e.b.Mod(e.b, f.p).Sign() == 0
}

func (f fp2) Rand(r io.Reader) Elt {
	a, _ := rand.Int(r, f.p)
	b, _ := rand.Int(r, f.p)
	return &fp2Elt{a, b}
}

func (f fp2) mod(a, b *big.Int) Elt { return &fp2Elt{a: a.Mod(a, f.p), b: b.Mod(b, f.p)} }
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
func (f fp2) Inv0(x Elt) Elt            { return f.Inv(x) }
func (f fp2) Inv(x Elt) Elt             { return nil }
func (f fp2) Neg(x Elt) Elt             { return nil }
func (f fp2) Sgn0(x Elt) int            { return 0 }
func (f fp2) CMov(x, y Elt, b bool) Elt { return nil }
func (f fp2) Sqrt(x Elt) Elt            { return nil }
func (f fp2) IsSquare(x Elt) bool       { return false }
