//+build go1.10

package hexenc_test

import (
	"go/doc"
	"go/parser"
	"go/token"
	"os"
	"runtime"
	"strings"
	"testing"
)

// TestDoc compares our doc with the reference in package encoding/hex.
func TestDoc(t *testing.T) {
	filter := func(info os.FileInfo) bool {
		return !strings.HasSuffix(info.Name(), "_test.go")
	}

	pkgs, err := parser.ParseDir(token.NewFileSet(), ".", filter, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	methodsDoc := make(map[string]string)

	p := doc.New(pkgs["hexenc"], ".", 0)

	for _, t := range p.Types {
		if t.Name != "Encoding" {
			continue
		}
		for _, meth := range t.Methods {
			methodsDoc[meth.Name] = meth.Doc
		}
	}

	pkgs = nil

	pkgs, err = parser.ParseDir(token.NewFileSet(), runtime.GOROOT()+"/src/encoding/hex", filter, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	p = doc.New(pkgs["hex"], "encoding/hex", 0)
	for _, f := range p.Funcs {
		if d, exists := methodsDoc[f.Name]; exists {
			t.Logf(
				"%s:\n%s---------------------------------------------------------------\n%s",
				f.Name, d, f.Doc)
		}
	}
}
