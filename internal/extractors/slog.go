package extractors

import (
	"go/ast"
	"go/token"

	"github.com/wurlinney/go-log-linter/internal/core"
	"golang.org/x/tools/go/analysis"
)

// SlogExtractor обрабатывает вызовы slog.
type SlogExtractor struct{}

func NewSlogExtractor() *SlogExtractor {
	return &SlogExtractor{}
}

// Match сообщает, относится ли вызов к slog.
func (e *SlogExtractor) Match(call *ast.CallExpr, _ *analysis.Pass) bool {
	if call == nil || len(call.Args) == 0 {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	if pkgIdent.Name != "slog" {
		return false
	}

	return true
}

// Extract формирует LogEntry для slog вызова.
func (e *SlogExtractor) Extract(call *ast.CallExpr, _ *analysis.Pass) (*core.LogEntry, bool) {
	if call == nil || len(call.Args) == 0 {
		return nil, false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	pkgIdent, ok := sel.X.(*ast.Ident)
	if !ok || pkgIdent.Name != "slog" {
		return nil, false
	}

	var msg string
	if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
		msg = trimQuotes(lit.Value)
	}

	entry := &core.LogEntry{
		Logger:      "slog",
		Level:       sel.Sel.Name,
		Message:     msg,
		MessageExpr: call.Args[0],
		Pos:         call.Pos(),
		End:         call.End(),
		Call:        call,
	}

	return entry, true
}

// trimQuotes убирает внешние кавычки у строкового литерала.
func trimQuotes(v string) string {
	if len(v) >= 2 {
		if (v[0] == '"' && v[len(v)-1] == '"') || (v[0] == '`' && v[len(v)-1] == '`') {
			return v[1 : len(v)-1]
		}
	}
	return v
}
