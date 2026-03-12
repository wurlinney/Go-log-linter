package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/wurlinney/go-log-linter/internal/core"
	"golang.org/x/tools/go/analysis"
)

type fakeEngine struct {
	calls int
}

func (f *fakeEngine) Analyze(*ast.CallExpr, *analysis.Pass) []core.Violation {
	f.calls++
	return []core.Violation{
		{
			RuleID:  "test",
			Message: "violation",
		},
	}
}

func TestRun_ReportsDiagnosticsForEachCallExpr(t *testing.T) {
	src := `package p

func f() {
	foo()
	bar()
}
`

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}

	var diags []analysis.Diagnostic

	pass := &analysis.Pass{
		Fset:  fset,
		Files: []*ast.File{file},
		Report: func(d analysis.Diagnostic) {
			diags = append(diags, d)
		},
	}

	eng := &fakeEngine{}

	run(pass, eng)

	if eng.calls != 2 {
		t.Fatalf("expected engine.Analyze to be called 2 times, got %d", eng.calls)
	}

	if len(diags) != 2 {
		t.Fatalf("expected 2 diagnostics, got %d", len(diags))
	}

	for _, d := range diags {
		if d.Message != "violation" {
			t.Errorf("unexpected diagnostic message: %q", d.Message)
		}
	}
}
