package rules

import (
	"go/parser"
	"testing"

	"github.com/wurlinney/go-log-linter/internal/core"
)

func TestLowercaseRule_ValidMessage(t *testing.T) {
	rule := NewLowercaseRule()

	entry := core.LogEntry{
		Message: "starting server",
	}

	violations := rule.Check(entry)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestLowercaseRule_InvalidMessageStartsWithUpper(t *testing.T) {
	rule := NewLowercaseRule()

	entry := core.LogEntry{
		Message: "Starting server",
	}

	violations := rule.Check(entry)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].RuleID != "lowercase" {
		t.Errorf("RuleID = %q, want %q", violations[0].RuleID, "lowercase")
	}
}

func TestLowercaseRule_SuggestedFixForLiteral(t *testing.T) {
	rule := NewLowercaseRule()

	expr, err := parser.ParseExpr(`"Starting server"`)
	if err != nil {
		t.Fatalf("parse expr: %v", err)
	}

	entry := core.LogEntry{
		Message:     "Starting server",
		MessageExpr: expr,
	}

	violations := rule.Check(entry)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	v := violations[0]
	if v.Fix == nil {
		t.Fatalf("expected SuggestedFix to be set")
	}
	if len(v.Fix.TextEdits) != 1 {
		t.Fatalf("expected 1 TextEdit, got %d", len(v.Fix.TextEdits))
	}
	if string(v.Fix.TextEdits[0].NewText) != `"starting server"` {
		t.Fatalf("unexpected fix text: %s", string(v.Fix.TextEdits[0].NewText))
	}
}
