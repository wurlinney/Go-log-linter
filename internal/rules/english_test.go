package rules

import (
	"testing"

	"github.com/wurlinney/go-log-linter/internal/core"
)

func TestEnglishRule_ValidEnglish(t *testing.T) {
	rule := NewEnglishRule()

	entry := core.LogEntry{
		Message: "starting server",
	}

	violations := rule.Check(entry)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestEnglishRule_InvalidNonEnglish(t *testing.T) {
	rule := NewEnglishRule()

	entry := core.LogEntry{
		Message: "запуск сервера",
	}

	violations := rule.Check(entry)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].RuleID != "english" {
		t.Errorf("RuleID = %q, want %q", violations[0].RuleID, "english")
	}
}
