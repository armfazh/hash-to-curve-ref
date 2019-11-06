package field

import (
	"crypto/sha512"
	"fmt"
	"testing"
)

func TestHashToField(t *testing.T) {
	H := sha512.New
	msg := "hello"
	dst := "asdf"
	ctr := byte(0)
	F3 := NewFromID(P448)

	fmt.Printf("msg: %v\n", msg)
	fmt.Printf("dst: %v\n", dst)
	fmt.Printf("F: %v\n", F3)
	a3 := HashToField([]byte(msg), []byte(dst), ctr, H, F3, 84)
	fmt.Printf("u: %v\n", a3)
	want := "0x9ce2583e1380f1e2bbb9a00267c504c17bd4fa5ddb1ce304a99f842163ca774bb1b934813adee2858f15b94a8eb7b7668dfa22870bcc8cbd"
	got := fmt.Sprintf("%v", a3)
	if want != got {
		t.Errorf("want: %v\ngot:%v\n", want, got)
	}
}

func BenchmarkHashToField(b *testing.B) {
	H := sha512.New
	msg := "hello"
	dst := "asdf"
	ctr := byte(0)
	F3 := NewFromID(P448)
	for i := 0; i < b.N; i++ {
		_ = HashToField([]byte(msg), []byte(dst), ctr, H, F3, 84)
	}
}
