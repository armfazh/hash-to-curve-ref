package h2c

import (
	"testing"
)

func TestSuite(t *testing.T) {
	msg := []byte("hello")
	dst := []byte("world")
	for _, name := range suiteNames {
		if suite := suites[name]; suite != nil {
			t.Logf("%v\n", name)
			suite.Hash(msg, dst)
		} else {
			t.Logf("Not Supported Yet: %v\n", name)
		}
	}
}
