package curve

import (
	"fmt"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

type mt2wec struct {
	E0   *MTCurve
	E1   *WCCurve
	invB GF.Elt
}

func (e *MTCurve) ToWeierstrassC() RationalMap {
	F := e.Field()
	invB := F.Inv(e.params.B)
	a := F.Mul(invB, e.params.A)
	b := F.Sqr(invB)
	return &mt2wec{E0: e, E1: NewWeierstrassC(Custom, F, a, b, e.params.R, e.params.H), invB: invB}
}

func (r *mt2wec) Domain() EllCurve   { return r.E0 }
func (r *mt2wec) Codomain() EllCurve { return r.E1 }
func (r *mt2wec) Push(p Point) Point {
	if p.IsIdentity() {
		return r.E1.Identity()
	}
	F := r.E0.Field()

	P := p.(*ptMt)
	x := F.Mul(P.x, r.invB) // s = x/B
	y := F.Mul(P.y, r.invB) // t = y/B
	return r.E1.NewPoint(x, y)
}
func (r *mt2wec) Pull(p Point) Point {
	if p.IsIdentity() {
		return r.E0.Identity()
	}
	F := r.E0.Field()
	P := p.(*ptWc)
	x := F.Mul(P.x, r.E0.B) // x = s*B
	y := F.Mul(P.y, r.E0.B) // y = t*B
	return r.E0.NewPoint(x, y)
}

type te2wec struct {
	E0       *TECurve
	E1       *WCCurve
	invSqrtD GF.Elt // 4/(a-d)
}

func (e *TECurve) ToWeierstrassC() RationalMap {
	F := e.Field()
	half := F.Inv(F.Elt(2))             // 1/2
	t0 := F.Add(e.params.A, e.params.D) // a+d
	a := F.Mul(t0, half)                // A = (a+d)/2

	t0 = F.Sub(e.params.A, e.params.D) // a-d
	t0 = F.Mul(t0, half)               // (a-d)/2
	t0 = F.Mul(t0, half)               // (a-d)/4
	invSqrtD := F.Inv(t0)              // 4/(a-d)
	b := F.Sqr(t0)                     // B = (a-d)^2/16
	return &te2wec{E0: e, E1: NewWeierstrassC(Custom, F, a, b, e.params.R, e.params.H), invSqrtD: invSqrtD}
}

func (r *te2wec) Domain() EllCurve   { return r.E0 }
func (r *te2wec) Codomain() EllCurve { return r.E1 }
func (r *te2wec) Push(p Point) Point {
	if p.IsIdentity() {
		return r.E1.Identity()
	}
	F := r.E0.Field()
	P := p.(*ptTe)
	t0 := F.Add(F.One(), P.y)  // 1+y
	t1 := F.Sub(F.One(), P.y)  // 1-y
	t1 = F.Mul(t1, r.invSqrtD) // invSqrtD*(1-y)
	t1 = F.Inv(t1)             // 1/(invSqrtD*(1-y))
	x := F.Mul(t0, t1)         // x = (1+y)/(invSqrtD*(1-y))
	t0 = F.Inv(P.y)            // 1/y
	y := F.Mul(x, t0)          // y = x/y
	return r.E1.NewPoint(x, y)
}
func (r *te2wec) Pull(p Point) Point {
	if p.IsIdentity() {
		return r.E0.Identity()
	}
	P := p.(*ptWc)
	F := r.E0.Field()
	if P.IsTwoTorsion() {
		return r.E0.NewPoint(F.Zero(), F.Elt(-1))
	}
	t0 := F.Inv(P.y)            // 1/y
	x := F.Mul(P.x, t0)         // X = x/y
	t0 = F.Mul(r.invSqrtD, P.x) // invSqrtD*x
	t1 := F.Add(t0, F.One())    // invSqrtD*x+1
	t2 := F.Sub(t0, F.One())    // invSqrtD*x-1
	t1 = F.Inv(t1)              // 1/(invSqrtD*x+1)
	y := F.Mul(t1, t2)          // Y = (invSqrtD*x-1)/(invSqrtD*x+1)
	return r.E0.NewPoint(x, y)
}

type wc2we struct {
	E0    *WCCurve
	E1    *WECurve
	Adiv3 GF.Elt
}

func (e *WCCurve) ToWeierstrass() RationalMap {
	F := e.Field()
	var t0, t1 GF.Elt
	t0 = F.Inv(F.Elt(3))    // 1/3
	Adiv3 := F.Mul(t0, e.A) // A/3
	t0 = F.Neg(Adiv3)       // -A/3
	t0 = F.Mul(t0, e.A)     // -A^2/3
	A := F.Add(t0, e.B)     // -A^2/3 + B

	t0 = F.Mul(F.Elt(9), e.B) // 9B
	t1 = F.Sqr(e.A)           // A^2
	t1 = F.Add(t1, t1)        // 2A^2
	t1 = F.Sub(t1, t0)        // 2A^2 - 9B
	t1 = F.Mul(t1, e.A)       // A(2A^2 - 9B)
	t0 = F.Inv(F.Elt(27))     // 1/27
	B := F.Mul(t0, t1)        // A(2A^2 - 9B)/27
	return &wc2we{E0: e, E1: NewWeierstrass(Custom, F, A, B, e.params.R, e.params.H), Adiv3: Adiv3}
}
func (r *wc2we) Domain() EllCurve   { return r.E0 }
func (r *wc2we) Codomain() EllCurve { return r.E1 }
func (r *wc2we) Push(p Point) Point {
	if p.IsIdentity() {
		return r.E1.Identity()
	}
	P := p.(*ptWc)
	F := r.E0.Field()
	xx := F.Add(P.x, r.Adiv3)
	return r.E1.NewPoint(xx, P.y)
}
func (r *wc2we) Pull(p Point) Point {
	if p.IsIdentity() {
		return r.E0.Identity()
	}
	P := p.(*ptWe)
	F := r.E0.Field()
	xx := F.Sub(P.x, r.Adiv3)
	return r.E0.NewPoint(xx, P.y)
}

type te2mt25519 struct {
	E0       T
	E1       M
	invSqrtD GF.Elt // sqrt(-486664) such that sgn0(sqrt_neg_486664) == 1
}

// FromTe2Mt25519 returns the birational map between Edwards25519 and Curve25519 curves.
func FromTe2Mt25519() RationalMap {
	e0 := Edwards25519.Get()
	e1 := Curve25519.Get()
	F := e0.Field()
	return te2mt25519{
		E0:       e0.(T),
		E1:       e1.(M),
		invSqrtD: F.Elt("6853475219497561581579357271197624642482790079785650197046958215289687604742"),
	}
}
func (m te2mt25519) String() string     { return fmt.Sprintf("Rational Map from %v to\n%v", m.E0, m.E1) }
func (m te2mt25519) Domain() EllCurve   { return m.E0 }
func (m te2mt25519) Codomain() EllCurve { return m.E1 }
func (m te2mt25519) Push(p Point) Point {
	if p.IsIdentity() {
		return m.E1.Identity()
	}
	F := m.E0.Field()
	x0, y0, one := p.X(), p.Y(), F.One()
	t0 := F.Add(one, y0)       // 1+y
	t1 := F.Sub(one, y0)       // 1-y
	t1 = F.Inv(t1)             // 1/(1-y)
	xx := F.Mul(t0, t1)        // xx = (1+y)/(1-y)
	t0 = F.Inv(x0)             // 1/x
	t0 = F.Mul(t0, m.invSqrtD) // invSqrtD/x
	yy := F.Mul(t0, xx)        // yy = invSqrtD*xx/x
	return m.E1.NewPoint(xx, yy)
}
func (m te2mt25519) Pull(p Point) Point {
	if p.IsIdentity() {
		return m.E0.Identity()
	}
	F := m.E0.Field()
	if p.IsTwoTorsion() {
		return m.E0.NewPoint(F.Zero(), F.Elt(-1))
	}
	x0, y0, one := p.X(), p.Y(), F.One()
	t0 := F.Inv(y0)             // 1/y
	t0 = F.Mul(t0, x0)          // x/y
	xx := F.Mul(m.invSqrtD, t0) // xx = invSqrtD*x/y
	t0 = F.Add(x0, one)         // x+1
	t1 := F.Sub(x0, one)        // x-1
	t0 = F.Inv(t0)              // 1/(x+1)
	yy := F.Mul(t0, t1)         // yy = (x-1)/(x+1)
	return m.E0.NewPoint(xx, yy)
}

type te2mt4iso448 struct {
	E0 T
	E1 M
}

// FromTe2Mt4ISO448 returns the four-degree isogeny between Edwards448 and Curve448 curves.
func FromTe2Mt4ISO448() RationalMap       { return te2mt4iso448{Edwards448.Get().(T), Curve448.Get().(M)} }
func (m te2mt4iso448) String() string     { return fmt.Sprintf("4-Isogeny from %v to\n%v", m.E0, m.E1) }
func (m te2mt4iso448) Domain() EllCurve   { return m.E0 }
func (m te2mt4iso448) Codomain() EllCurve { return m.E1 }
func (m te2mt4iso448) Push(p Point) Point {
	if p.IsIdentity() {
		return m.E1.Identity()
	}
	F := m.E0.Field()
	x, y := p.X(), p.Y()
	t0 := F.Inv(x)           // 1/x
	t1 := F.Sqr(t0)          // 1/x^2
	t2 := F.Mul(t0, t1)      // 1/x^3
	t0 = F.Sqr(y)            // y^2
	xx := F.Mul(t0, t1)      // xx = y^2/x^2
	t0 = F.Sub(F.Elt(2), t0) // 2-y^2
	t1 = F.Sqr(x)            // x^2
	t0 = F.Sub(t0, t1)       // 2-y^2-x^2
	t0 = F.Mul(t0, y)        // (2-y^2-x^2)*y
	yy := F.Mul(t0, t2)      // (2-y^2-x^2)*y/x^3
	return m.E1.NewPoint(xx, yy)
}
func (m te2mt4iso448) Pull(p Point) Point {
	if p.IsIdentity() {
		return m.E0.Identity()
	}
	F := m.E0.Field()
	x, y, one := p.X(), p.Y(), F.One()

	t0 := F.Sqr(x)       // x^2
	t1 := F.Add(t0, one) // x^2+1
	t0 = F.Sub(t0, one)  // x^2-1
	t2 := F.Sqr(y)       // y^2
	t2 = F.Add(t2, t2)   // 2y^2
	t3 := F.Add(x, x)    // 2x

	t4 := F.Mul(t0, y)    // y(x^2-1)
	t4 = F.Add(t4, t4)    // 2y(x^2-1)
	xNum := F.Add(t4, t4) // xNum = 4y(x^2-1)

	t5 := F.Sqr(t0)       // x^4-2x^2+1
	t4 = F.Add(t5, t2)    // x^4-2x^2+1+2y^2
	xDen := F.Add(t4, t2) // xDen = x^4-2x^2+1+4y^2

	t5 = F.Mul(t5, x)     // x^5-2x^3+x
	t4 = F.Mul(t2, t3)    // 4xy^2
	yNum := F.Sub(t4, t5) // yNum = -(x^5-2x^3+x-4xy^2)

	t4 = F.Mul(t1, t2)    // 2x^2y^2+2y^2
	yDen := F.Sub(t5, t4) // yDen = x^5-2x^3+x-2x^2y^2-2y^2

	xx := F.Mul(xNum, F.Inv(xDen))
	yy := F.Mul(yNum, F.Inv(yDen))
	return m.E0.NewPoint(xx, yy)
}

type isosecp256k1 struct {
	E0, E1                 W
	xNum, xDen, yNum, yDen []GF.Elt
}

// GetSECP256K1Isogeny returns a 3-degree isogeny from SECP256K1_3ISO to the SECP256K1 elliptic curve.
func GetSECP256K1Isogeny() Isogeny {
	e0 := SECP256K1_3ISO.Get()
	e1 := SECP256K1.Get()
	F := e0.Field()
	return isosecp256k1{
		E0: e0.(W),
		E1: e1.(W),
		xNum: []GF.Elt{
			F.Elt("0x8e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38daaaaa8c7"),
			F.Elt("0x07d3d4c80bc321d5b9f315cea7fd44c5d595d2fc0bf63b92dfff1044f17c6581"),
			F.Elt("0x534c328d23f234e6e2a413deca25caece4506144037c40314ecbd0b53d9dd262"),
			F.Elt("0x8e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38e38daaaaa88c")},
		xDen: []GF.Elt{
			F.Elt("0xd35771193d94918a9ca34ccbb7b640dd86cd409542f8487d9fe6b745781eb49b"),
			F.Elt("0xedadc6f64383dc1df7c4b2d51b54225406d36b641f5e41bbc52a56612a8c6d14"),
			F.One(),
			F.Zero()},
		yNum: []GF.Elt{
			F.Elt("0x4bda12f684bda12f684bda12f684bda12f684bda12f684bda12f684b8e38e23c"),
			F.Elt("0xc75e0c32d5cb7c0fa9d0a54b12a0a6d5647ab046d686da6fdffc90fc201d71a3"),
			F.Elt("0x29a6194691f91a73715209ef6512e576722830a201be2018a765e85a9ecee931"),
			F.Elt("0x2f684bda12f684bda12f684bda12f684bda12f684bda12f684bda12f38e38d84")},
		yDen: []GF.Elt{
			F.Elt("0xfffffffffffffffffffffffffffffffffffffffffffffffffffffffefffff93b"),
			F.Elt("0x7a06534bb8bdb49fd5e9e6632722c2989467c1bfc8e8d978dfb425d2685c2573"),
			F.Elt("0x6484aa716545ca2cf3a70c3fa8fe337e0a3d21162f0d6299a7bf8192bfd2a76f"),
			F.One()},
	}
}
func (m isosecp256k1) String() string     { return fmt.Sprintf("3-Isogeny from %v to\n%v", m.E0, m.E1) }
func (m isosecp256k1) Domain() EllCurve   { return m.E0 }
func (m isosecp256k1) Codomain() EllCurve { return m.E1 }
func (m isosecp256k1) Push(p Point) Point {
	F := m.E0.F
	x, y := p.X(), p.Y()
	xNum, xDen, yNum, yDen := F.Zero(), F.Zero(), F.Zero(), F.Zero()
	for i := 3; i >= 0; i-- {
		xNum = F.Add(F.Mul(xNum, x), m.xNum[i])
		xDen = F.Add(F.Mul(xDen, x), m.xDen[i])
		yNum = F.Add(F.Mul(yNum, x), m.yNum[i])
		yDen = F.Add(F.Mul(yDen, x), m.yDen[i])
	}
	xx := F.Mul(xNum, F.Inv(xDen))
	yy := F.Mul(yNum, F.Inv(yDen))
	yy = F.Mul(yy, y)
	return m.E1.NewPoint(xx, yy)
}

type isobls12381G1 struct {
	E0, E1                 W
	xNum, xDen, yNum, yDen []GF.Elt
}

// GetSECP256K1Isogeny returns an 11-degree isogeny from BLS12381G1_11ISO to the BLS12381G1 elliptic curve.
func GetBLS12381G1Isogeny() Isogeny {
	e0 := BLS12381G1_11ISO.Get()
	e1 := BLS12381G1.Get()
	F := e0.Field()
	return isobls12381G1{
		E0: e0.(W),
		E1: e1.(W),
		xNum: []GF.Elt{
			F.Elt("0x11a05f2b1e833340b809101dd99815856b303e88a2d7005ff2627b56cdb4e2c85610c2d5f2e62d6eaeac1662734649b7"),
			F.Elt("0x17294ed3e943ab2f0588bab22147a81c7c17e75b2f6a8417f565e33c70d1e86b4838f2a6f318c356e834eef1b3cb83bb"),
			F.Elt("0xd54005db97678ec1d1048c5d10a9a1bce032473295983e56878e501ec68e25c958c3e3d2a09729fe0179f9dac9edcb0"),
			F.Elt("0x1778e7166fcc6db74e0609d307e55412d7f5e4656a8dbf25f1b33289f1b330835336e25ce3107193c5b388641d9b6861"),
			F.Elt("0xe99726a3199f4436642b4b3e4118e5499db995a1257fb3f086eeb65982fac18985a286f301e77c451154ce9ac8895d9"),
			F.Elt("0x1630c3250d7313ff01d1201bf7a74ab5db3cb17dd952799b9ed3ab9097e68f90a0870d2dcae73d19cd13c1c66f652983"),
			F.Elt("0xd6ed6553fe44d296a3726c38ae652bfb11586264f0f8ce19008e218f9c86b2a8da25128c1052ecaddd7f225a139ed84"),
			F.Elt("0x17b81e7701abdbe2e8743884d1117e53356de5ab275b4db1a682c62ef0f2753339b7c8f8c8f475af9ccb5618e3f0c88e"),
			F.Elt("0x80d3cf1f9a78fc47b90b33563be990dc43b756ce79f5574a2c596c928c5d1de4fa295f296b74e956d71986a8497e317"),
			F.Elt("0x169b1f8e1bcfa7c42e0c37515d138f22dd2ecb803a0c5c99676314baf4bb1b7fa3190b2edc0327797f241067be390c9e"),
			F.Elt("0x10321da079ce07e272d8ec09d2565b0dfa7dccdde6787f96d50af36003b14866f69b771f8c285decca67df3f1605fb7b"),
			F.Elt("0x6e08c248e260e70bd1e962381edee3d31d79d7e22c837bc23c0bf1bc24c6b68c24b1b80b64d391fa9c8ba2e8ba2d229")},
		xDen: []GF.Elt{
			F.Elt("0x8ca8d548cff19ae18b2e62f4bd3fa6f01d5ef4ba35b48ba9c9588617fc8ac62b558d681be343df8993cf9fa40d21b1c"),
			F.Elt("0x12561a5deb559c4348b4711298e536367041e8ca0cf0800c0126c2588c48bf5713daa8846cb026e9e5c8276ec82b3bff"),
			F.Elt("0xb2962fe57a3225e8137e629bff2991f6f89416f5a718cd1fca64e00b11aceacd6a3d0967c94fedcfcc239ba5cb83e19"),
			F.Elt("0x3425581a58ae2fec83aafef7c40eb545b08243f16b1655154cca8abc28d6fd04976d5243eecf5c4130de8938dc62cd8"),
			F.Elt("0x13a8e162022914a80a6f1d5f43e7a07dffdfc759a12062bb8d6b44e833b306da9bd29ba81f35781d539d395b3532a21e"),
			F.Elt("0xe7355f8e4e667b955390f7f0506c6e9395735e9ce9cad4d0a43bcef24b8982f7400d24bc4228f11c02df9a29f6304a5"),
			F.Elt("0x772caacf16936190f3e0c63e0596721570f5799af53a1894e2e073062aede9cea73b3538f0de06cec2574496ee84a3a"),
			F.Elt("0x14a7ac2a9d64a8b230b3f5b074cf01996e7f63c21bca68a81996e1cdf9822c580fa5b9489d11e2d311f7d99bbdcc5a5e"),
			F.Elt("0xa10ecf6ada54f825e920b3dafc7a3cce07f8d1d7161366b74100da67f39883503826692abba43704776ec3a79a1d641"),
			F.Elt("0x95fc13ab9e92ad4476d6e3eb3a56680f682b4ee96f7d03776df533978f31c1593174e4b4b7865002d6384d168ecdd0a"),
			F.One(),
			F.Zero()},
		yNum: []GF.Elt{
			F.Elt("0x90d97c81ba24ee0259d1f094980dcfa11ad138e48a869522b52af6c956543d3cd0c7aee9b3ba3c2be9845719707bb33"),
			F.Elt("0x134996a104ee5811d51036d776fb46831223e96c254f383d0f906343eb67ad34d6c56711962fa8bfe097e75a2e41c696"),
			F.Elt("0xcc786baa966e66f4a384c86a3b49942552e2d658a31ce2c344be4b91400da7d26d521628b00523b8dfe240c72de1f6"),
			F.Elt("0x1f86376e8981c217898751ad8746757d42aa7b90eeb791c09e4a3ec03251cf9de405aba9ec61deca6355c77b0e5f4cb"),
			F.Elt("0x8cc03fdefe0ff135caf4fe2a21529c4195536fbe3ce50b879833fd221351adc2ee7f8dc099040a841b6daecf2e8fedb"),
			F.Elt("0x16603fca40634b6a2211e11db8f0a6a074a7d0d4afadb7bd76505c3d3ad5544e203f6326c95a807299b23ab13633a5f0"),
			F.Elt("0x4ab0b9bcfac1bbcb2c977d027796b3ce75bb8ca2be184cb5231413c4d634f3747a87ac2460f415ec961f8855fe9d6f2"),
			F.Elt("0x987c8d5333ab86fde9926bd2ca6c674170a05bfe3bdd81ffd038da6c26c842642f64550fedfe935a15e4ca31870fb29"),
			F.Elt("0x9fc4018bd96684be88c9e221e4da1bb8f3abd16679dc26c1e8b6e6a1f20cabe69d65201c78607a360370e577bdba587"),
			F.Elt("0xe1bba7a1186bdb5223abde7ada14a23c42a0ca7915af6fe06985e7ed1e4d43b9b3f7055dd4eba6f2bafaaebca731c30"),
			F.Elt("0x19713e47937cd1be0dfd0b8f1d43fb93cd2fcbcb6caf493fd1183e416389e61031bf3a5cce3fbafce813711ad011c132"),
			F.Elt("0x18b46a908f36f6deb918c143fed2edcc523559b8aaf0c2462e6bfe7f911f643249d9cdf41b44d606ce07c8a4d0074d8e"),
			F.Elt("0xb182cac101b9399d155096004f53f447aa7b12a3426b08ec02710e807b4633f06c851c1919211f20d4c04f00b971ef8"),
			F.Elt("0x245a394ad1eca9b72fc00ae7be315dc757b3b080d4c158013e6632d3c40659cc6cf90ad1c232a6442d9d3f5db980133"),
			F.Elt("0x5c129645e44cf1102a159f748c4a3fc5e673d81d7e86568d9ab0f5d396a7ce46ba1049b6579afb7866b1e715475224b"),
			F.Elt("0x15e6be4e990f03ce4ea50b3b42df2eb5cb181d8f84965a3957add4fa95af01b2b665027efec01c7704b456be69c8b604")},
		yDen: []GF.Elt{
			F.Elt("0x16112c4c3a9c98b252181140fad0eae9601a6de578980be6eec3232b5be72e7a07f3688ef60c206d01479253b03663c1"),
			F.Elt("0x1962d75c2381201e1a0cbd6c43c348b885c84ff731c4d59ca4a10356f453e01f78a4260763529e3532f6102c2e49a03d"),
			F.Elt("0x58df3306640da276faaae7d6e8eb15778c4855551ae7f310c35a5dd279cd2eca6757cd636f96f891e2538b53dbf67f2"),
			F.Elt("0x16b7d288798e5395f20d23bf89edb4d1d115c5dbddbcd30e123da489e726af41727364f2c28297ada8d26d98445f5416"),
			F.Elt("0xbe0e079545f43e4b00cc912f8228ddcc6d19c9f0f69bbb0542eda0fc9dec916a20b15dc0fd2ededda39142311a5001d"),
			F.Elt("0x8d9e5297186db2d9fb266eaac783182b70152c65550d881c5ecd87b6f0f5a6449f38db9dfa9cce202c6477faaf9b7ac"),
			F.Elt("0x166007c08a99db2fc3ba8734ace9824b5eecfdfa8d0cf8ef5dd365bc400a0051d5fa9c01a58b1fb93d1a1399126a775c"),
			F.Elt("0x16a3ef08be3ea7ea03bcddfabba6ff6ee5a4375efa1f4fd7feb34fd206357132b920f5b00801dee460ee415a15812ed9"),
			F.Elt("0x1866c8ed336c61231a1be54fd1d74cc4f9fb0ce4c6af5920abc5750c4bf39b4852cfe2f7bb9248836b233d9d55535d4a"),
			F.Elt("0x167a55cda70a6e1cea820597d94a84903216f763e13d87bb5308592e7ea7d4fbc7385ea3d529b35e346ef48bb8913f55"),
			F.Elt("0x4d2f259eea405bd48f010a01ad2911d9c6dd039bb61a6290e591b36e636a5c871a5c29f4f83060400f8b49cba8f6aa8"),
			F.Elt("0xaccbb67481d033ff5852c1e48c50c477f94ff8aefce42d28c0f9a88cea7913516f968986f7ebbea9684b529e2561092"),
			F.Elt("0xad6b9514c767fe3c3613144b45f1496543346d98adf02267d5ceef9a00d9b8693000763e3b90ac11e99b138573345cc"),
			F.Elt("0x2660400eb2e4f3b628bdd0d53cd76f2bf565b94e72927c1cb748df27942480e420517bd8714cc80d1fadc1326ed06f7"),
			F.Elt("0xe0fa1d816ddc03e6b24255e0d7819c171c40f65e273b853324efcd6356caa205ca2f570f13497804415473a1d634b8f"),
			F.One()},
	}
}
func (m isobls12381G1) String() string     { return fmt.Sprintf("11-Isogeny from %v to\n%v", m.E0, m.E1) }
func (m isobls12381G1) Domain() EllCurve   { return m.E0 }
func (m isobls12381G1) Codomain() EllCurve { return m.E1 }
func (m isobls12381G1) Push(p Point) Point {
	F := m.E0.F
	x, y := p.X(), p.Y()
	xNum, xDen := F.Zero(), F.Zero()
	for i := len(m.xNum) - 1; i >= 0; i-- {
		xNum = F.Add(F.Mul(xNum, x), m.xNum[i])
		xDen = F.Add(F.Mul(xDen, x), m.xDen[i])
	}
	yNum, yDen := F.Zero(), F.Zero()
	for i := len(m.yNum) - 1; i >= 0; i-- {
		yNum = F.Add(F.Mul(yNum, x), m.yNum[i])
		yDen = F.Add(F.Mul(yDen, x), m.yDen[i])
	}
	xx := F.Mul(xNum, F.Inv(xDen))
	yy := F.Mul(yNum, F.Inv(yDen))
	yy = F.Mul(yy, y)
	return m.E1.NewPoint(xx, yy)
}
