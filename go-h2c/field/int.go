package field

import (
	"fmt"
	"math/big"
)

type modulus struct {
	name string
	p    *big.Int
}

func (m modulus) fromType(in interface{}) *big.Int {
	p := fromType(in)
	return p.Mod(p, m.p)
}

func bigFromString(s string) *big.Int {
	p := new(big.Int)
	if _, ok := p.SetString(s, 0); !ok {
		panic("error setting the number")
	}
	return p
}

func fromType(in interface{}) *big.Int {
	n := new(big.Int)
	switch s := in.(type) {
	case *big.Int:
		n.Set(s)
	case big.Int:
		n.Set(&s)
	case string:
		n = bigFromString(s)
	case uint:
		n.SetUint64(uint64(s))
	case uint8:
		n.SetUint64(uint64(s))
	case uint16:
		n.SetUint64(uint64(s))
	case uint32:
		n.SetUint64(uint64(s))
	case uint64:
		n.SetUint64(uint64(s))
	case int:
		n.SetInt64(int64(s))
	case int8:
		n.SetInt64(int64(s))
	case int16:
		n.SetInt64(int64(s))
	case int32:
		n.SetInt64(int64(s))
	case int64:
		n.SetInt64(int64(s))
	default:
		panic(fmt.Errorf("type %T not supported", in))
	}
	return n
}
