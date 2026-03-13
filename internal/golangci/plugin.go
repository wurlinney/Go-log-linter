package golangci

import (
	"github.com/wurlinney/go-log-linter/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

// Analyzers экспортируется для использования golangci-lint.
var Analyzers = []*analysis.Analyzer{
	analyzer.Analyzer,
}

