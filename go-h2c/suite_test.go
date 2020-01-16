package h2c_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	h2c "github.com/armfazh/hash-to-curve-ref/go-h2c"
)

// func TestSuite(t *testing.T) {
// 	msg := []byte("hello")
// 	dst := []byte("world")
// 	for name, suite := range h2c.Suites {
// 		t.Logf("%v\n", name)
// 		suite.Hash(msg, dst)
// 	}
// }

type vectorSuite struct {
	SuiteID   string `json:"ciphersuite"`
	CurveName string `json:"curve"`
	DST       string `json:"dst"`
	Field     struct {
		M string `json:"m"`
		P string `json:"p"`
	} `json:"field"`
	Hash string `json:"hash"`
	Map  struct {
		Name string `json:"name"`
		Sgn0 string `json:"sgn0"`
	} `json:"map"`
	RandomOracle bool `json:"randomOracle"`
	Vectors      []struct {
		P struct {
			X string `json:"x"`
			Y string `json:"y"`
		} `json:"P"`
		Msg string `json:"msg"`
	} `json:"vectors"`
}

func (v vectorSuite) testVector(t *testing.T) {
	t.Helper()
	hashToCurve, ok := h2c.Suites[v.SuiteID]
	if !ok {
		t.Skipf("suite not supported yet: %v\n", v.SuiteID)
	}
	params := hashToCurve.GetParams()
	E := params.E
	F := E.Field()
	for i := range v.Vectors {
		got := hashToCurve.Hash([]byte(v.Vectors[i].Msg), []byte(v.DST))
		want := E.NewPoint(
			F.Elt(v.Vectors[i].P.X),
			F.Elt(v.Vectors[i].P.Y),
		)
		if !got.IsEqual(want) {
			t.Fatalf("suite: %v\ngot:  %v\nwant: %v", v.SuiteID, got, want)
		}
	}
}

func TestVectors(t *testing.T) {
	if errFolder := filepath.Walk("testdata",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			jsonFile, errFile := os.Open(path)
			defer jsonFile.Close()
			if errFile != nil {
				return errFile
			}

			byteValue, errRead := ioutil.ReadAll(jsonFile)
			if errRead != nil {
				return errRead
			}
			var v vectorSuite
			errJSON := json.Unmarshal(byteValue, &v)
			if errJSON != nil {
				return errJSON
			}
			t.Run(v.SuiteID, v.testVector)
			return nil
		}); errFolder != nil {
		t.Fatalf("error on reading testdata folder: %v", errFolder)
	}
}
