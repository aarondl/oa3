package openapi3spec

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestYAML(t *testing.T) {
	t.Parallel()

	oa, err := LoadYAML("testdata/openapi3.yaml", false)
	if err != nil {
		t.Fatal(err)
	}

	out, err := json.MarshalIndent(oa, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	golden(t, out, "testdata/out.yaml", *flagUpdateGolden)
}

func golden(t *testing.T, now []byte, filename string, write bool) {
	t.Helper()

	if write {
		if err := ioutil.WriteFile(filename, now, 0o664); err != nil {
			t.Fatalf("failed to write file (%s): %v", filename, err)
		}
	}

	then, err := ioutil.ReadFile(filename)
	if os.IsNotExist(err) {
		t.Fatal("golden file does not exist:", filename)
	} else if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(now, then) {
		return
	}

	config := difflib.UnifiedDiff{
		A:        difflib.SplitLines(string(then)),
		B:        difflib.SplitLines(string(now)),
		FromFile: filename,
		ToFile:   "TestOutput",
		Context:  3,
	}

	text, err := difflib.GetUnifiedDiffString(config)
	if err != nil {
		t.Fatal(err)
	}

	t.Error("\n" + text)
}
