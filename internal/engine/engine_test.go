package engine

import (
	"go/ast"
	"testing"

	"github.com/wurlinney/go-log-linter/internal/core"
	"github.com/wurlinney/go-log-linter/internal/extractors"
	"golang.org/x/tools/go/analysis"
)

type fakeExtractor struct {
	shouldMatch bool
	entry       *core.LogEntry
}

func (f *fakeExtractor) Match(*ast.CallExpr, *analysis.Pass) bool {
	return f.shouldMatch
}

func (f *fakeExtractor) Extract(*ast.CallExpr, *analysis.Pass) (*core.LogEntry, bool) {
	if !f.shouldMatch || f.entry == nil {
		return nil, false
	}
	return f.entry, true
}

type fakeRule struct {
	id         string
	violations []core.Violation
}

func (r *fakeRule) ID() string { return r.id }

func (r *fakeRule) Check(core.LogEntry) []core.Violation {
	return r.violations
}

func TestEngine_Analyze_ExtractorAndRulesApplied(t *testing.T) {
	entry := &core.LogEntry{
		Logger:  "slog",
		Level:   "Info",
		Message: "starting server",
	}

	ex := &fakeExtractor{
		shouldMatch: true,
		entry:       entry,
	}

	v := core.Violation{
		RuleID:  "test-rule",
		Message: "violation",
	}

	rule := &fakeRule{
		id:         "test-rule",
		violations: []core.Violation{v},
	}

	eng := NewEngine([]extractors.Extractor{ex}, []core.Rule{rule})

	pass := &analysis.Pass{}
	call := &ast.CallExpr{}

	violations := eng.Analyze(call, pass)

	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].RuleID != "test-rule" {
		t.Errorf("RuleID = %q, want %q", violations[0].RuleID, "test-rule")
	}
}

func TestEngine_Analyze_NoExtractorMatch(t *testing.T) {
	ex := &fakeExtractor{
		shouldMatch: false,
	}

	eng := NewEngine([]extractors.Extractor{ex}, nil)

	pass := &analysis.Pass{}
	var call *ast.CallExpr

	violations := eng.Analyze(call, pass)

	if len(violations) != 0 {
		t.Fatalf("expected 0 violations when no extractor matches, got %d", len(violations))
	}
}
