package math

import "math/big"

type fquadratic struct{ ff }

type fquadraticElement struct {
	a, b *big.Int
}

func (e fquadraticElement) String() string {
	return "a: " + e.a.Text(16) + "\n" +
		"b: " + e.b.Text(16) + " * i"
}

// IsZero returns true is x is zero,
func (e fquadraticElement) IsZero() bool { return e.a.Sign() == 0 && e.b.Sign() == 0 }
func (f fquadratic) String() string {
	return "GF(" + f.ff.String() + ")\n" +
		"Irred: i^2+1"
}
func (f fquadratic) IsPrime() bool { return false }
func (f fquadratic) Elt() Element  { return fquadraticElement{big.NewInt(0), big.NewInt(0)} }
func (f fquadratic) One() Element  { return fquadraticElement{big.NewInt(1), big.NewInt(0)} }
func (f fquadratic) Zero() Element { return f.Elt() }
func (f fquadratic) Add(x, y Element) Element {
	a := new(big.Int).Add(x.(fquadraticElement).a, y.(fquadraticElement).a)
	b := new(big.Int).Add(x.(fquadraticElement).b, y.(fquadraticElement).b)
	a.Mod(a, f.p)
	b.Mod(b, f.p)
	return fquadraticElement{a, b}
}
func (f fquadratic) Sub(x, y Element) Element {
	a := new(big.Int).Sub(x.(fquadraticElement).a, y.(fquadraticElement).a)
	b := new(big.Int).Sub(x.(fquadraticElement).b, y.(fquadraticElement).b)
	a.Mod(a, f.p)
	b.Mod(b, f.p)
	return fquadraticElement{a, b}
}

func (f fquadratic) Mul(x, y Element) Element          { return nil }
func (f fquadratic) Sqr(x Element) Element             { return nil }
func (f fquadratic) Sqrt(x Element) Element            { return nil }
func (f fquadratic) Inv0(x Element) Element            { return nil }
func (f fquadratic) Neg(x Element) Element             { return nil }
func (f fquadratic) Sgn0(x Element) int                { return 0 }
func (f fquadratic) CMov(x, y Element, b bool) Element { return nil }
func (f fquadratic) IsSquare(x Element) bool           { return false }
