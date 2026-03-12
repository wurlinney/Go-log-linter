package testutil

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"golang.org/x/tools/go/analysis"
)

// ParseFirstCall разбирает исходник и возвращает первый CallExpr и Pass.
func ParseFirstCall(t *testing.T, src string) (*ast.CallExpr, *analysis.Pass) {
	t.Helper()

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", "package p\nfunc f() {"+src+"}", 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	var call *ast.CallExpr
	ast.Inspect(file, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok && call == nil {
			call = c
			return false
		}
		return true
	})
	if call == nil {
		t.Fatal("no CallExpr found")
	}

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
	}

	return call, pass
}
