package extractors

import (
	"go/ast"

	"github.com/wurlinney/go-log-linter/internal/core"
	"golang.org/x/tools/go/analysis"
)

// Extractor отвечает за распознавание лог вызовов и извлечение LogEntry из ast.
type Extractor interface {
	Match(call *ast.CallExpr, pass *analysis.Pass) bool
	Extract(call *ast.CallExpr, pass *analysis.Pass) (*core.LogEntry, bool)
}
