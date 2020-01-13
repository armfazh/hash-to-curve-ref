package field_test

import (
	"crypto/sha512"
	"testing"

	GF "github.com/armfazh/hash-to-curve-ref/go-h2c/field"
)

func TestHashToField(t *testing.T) {
	H := sha512.New
	msg := "hello"              // fmt.Printf("msg: %v\n", msg)
	dst := "asdf"               // fmt.Printf("dst: %v\n", dst)
	ctr := byte(0)              // fmt.Printf("ctr: %v\n", ctr)
	F3 := GF.GetFromID(GF.P448) // fmt.Printf("F: %v\n", F3)

	got := GF.HashToField([]byte(msg), []byte(dst), ctr, H, F3, 84) // fmt.Printf("u: %v\n", got)
	want := F3.Elt("0x9ce2583e1380f1e2bbb9a00267c504c17bd4fa5ddb1ce304a99f842163ca774bb1b934813adee2858f15b94a8eb7b7668dfa22870bcc8cbd")
	if !F3.AreEqual(got, want) {
		t.Errorf("want: %v\ngot:%v\n", want, got)
	}
}

func BenchmarkHashToField(b *testing.B) {
	H := sha512.New
	msg := "hello"
	dst := "asdf"
	ctr := byte(0)
	F3 := GF.GetFromID(GF.P448)
	for i := 0; i < b.N; i++ {
		_ = GF.HashToField([]byte(msg), []byte(dst), ctr, H, F3, 84)
	}
}
