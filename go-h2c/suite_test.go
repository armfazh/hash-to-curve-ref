package h2c

import (
	"testing"
)

func TestSuite(t *testing.T) {
	msg := []byte("hello")
	dst := []byte("world")
	for id, suite := range suites {
		t.Logf("%v\n", id)
		suite.Hash(msg, dst)
	}
}
