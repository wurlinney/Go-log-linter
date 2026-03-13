package extractors

import (
	"go/ast"
	"go/token"

	"github.com/wurlinney/go-log-linter/internal/core"
	"golang.org/x/tools/go/analysis"
)

// ZapExtractor обрабатывает вызовы zap логгера.
type ZapExtractor struct{}

func NewZapExtractor() *ZapExtractor {
	return &ZapExtractor{}
}

func (e *ZapExtractor) Match(call *ast.CallExpr, _ *analysis.Pass) bool {
	if call == nil || len(call.Args) == 0 {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Поддерживаем два типа вызовов:
	// 1) zap.L().Info("msg")
	// 2) logger.Info("msg") / log.Info("msg") - когда zap-логгер хранится в переменной
	if _, ok := sel.X.(*ast.CallExpr); ok {
		// Любой вызов вида something().Info(...) интерпретируем как zap.L().Info(...)
		return true
	}
	if ident, ok := sel.X.(*ast.Ident); ok {
		// Ограничиваемся типичными именами переменных-логгеров,
		// чтобы не ловить все pkg.Func(...) как zap-вызовы.
		return ident.Name == "logger" || ident.Name == "log"
	}

	return false
}

func (e *ZapExtractor) Extract(call *ast.CallExpr, _ *analysis.Pass) (*core.LogEntry, bool) {
	if call == nil || len(call.Args) == 0 {
		return nil, false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, false
	}

	var msg string
	if lit, ok := call.Args[0].(*ast.BasicLit); ok && lit.Kind == token.STRING {
		msg = trimQuotes(lit.Value)
	}

	entry := &core.LogEntry{
		Logger:      "zap",
		Level:       sel.Sel.Name,
		Message:     msg,
		MessageExpr: call.Args[0],
		Pos:         call.Pos(),
		End:         call.End(),
		Call:        call,
	}

	return entry, true
}
