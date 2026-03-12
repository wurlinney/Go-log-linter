package rules

import (
	"testing"

	"github.com/wurlinney/go-log-linter/internal/core"
)

func TestSymbolsRule_ValidMessage(t *testing.T) {
	rule := NewSymbolsRule()

	entry := core.LogEntry{
		Message: "server started",
	}

	violations := rule.Check(entry)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestSymbolsRule_InvalidExclamationMarks(t *testing.T) {
	rule := NewSymbolsRule()

	entry := core.LogEntry{
		Message: "server started!!!",
	}

	violations := rule.Check(entry)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].RuleID != "symbols" {
		t.Errorf("RuleID = %q, want %q", violations[0].RuleID, "symbols")
	}
}

func TestSymbolsRule_InvalidEmoji(t *testing.T) {
	rule := NewSymbolsRule()

	entry := core.LogEntry{
		Message: "server started 🚀",
	}

	violations := rule.Check(entry)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}
