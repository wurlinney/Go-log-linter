package engine

import (
	"go/ast"

	"github.com/wurlinney/go-log-linter/internal/core"
	"github.com/wurlinney/go-log-linter/internal/extractors"
	"golang.org/x/tools/go/analysis"
)

// Engine запускает экстракторы и правила.
type Engine struct {
	extractors []extractors.Extractor
	rules      []core.Rule
}

func NewEngine(extrs []extractors.Extractor, rules []core.Rule) *Engine {
	return &Engine{
		extractors: extrs,
		rules:      rules,
	}
}

// Analyze извлекает LogEntry и применяет правила.
func (e *Engine) Analyze(call *ast.CallExpr, pass *analysis.Pass) []core.Violation {
	if call == nil {
		return nil
	}

	var entry *core.LogEntry
	for _, ex := range e.extractors {
		if ex.Match(call, pass) {
			if le, ok := ex.Extract(call, pass); ok && le != nil {
				entry = le
				break
			}
		}
	}

	if entry == nil {
		return nil
	}

	var result []core.Violation
	for _, r := range e.rules {
		vios := r.Check(*entry)
		if len(vios) > 0 {
			result = append(result, vios...)
		}
	}

	return result
}
