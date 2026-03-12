package analyzer

import (
	"go/ast"

	"github.com/wurlinney/go-log-linter/internal/config"
	"github.com/wurlinney/go-log-linter/internal/core"
	"github.com/wurlinney/go-log-linter/internal/engine"
	"github.com/wurlinney/go-log-linter/internal/report"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/ast/inspector"
)

// Engine описывает минимальный контракт движка.
type Engine interface {
	Analyze(call *ast.CallExpr, pass *analysis.Pass) []core.Violation
}

// Analyzer используется с go vet и golangci-lint.
var Analyzer = newDefaultAnalyzer()

func newDefaultAnalyzer() *analysis.Analyzer {
	cfg, err := config.Load("")
	if err != nil {
		cfg = config.Default()
	}
	eng := engine.NewEngineWithConfig(cfg, nil)
	return New(eng)
}

// New создаёт новый analysis.Analyzer поверх переданного движка.
func New(eng Engine) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "loglint",
		Doc:  "checks log messages against consistency and safety rules",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			run(pass, eng)
			return nil, nil
		},
	}
}

// run обходит AST и конвертирует нарушения в diagnostics.
func run(pass *analysis.Pass, eng Engine) {
	if pass == nil || eng == nil {
		return
	}

	ins := inspector.New(pass.Files)

	ins.Preorder([]ast.Node{(*ast.CallExpr)(nil)}, func(n ast.Node) {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return
		}

		violations := eng.Analyze(call, pass)
		for _, v := range violations {
			diag := report.ViolationToDiagnostic(v)
			pass.Report(diag)
		}
	})
}
