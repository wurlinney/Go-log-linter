package rules

import (
	"go/parser"
	"testing"

	"github.com/wurlinney/go-log-linter/internal/core"
)

func TestSensitiveRule_AllowedMessage(t *testing.T) {
	rule := NewSensitiveRule()

	entry := core.LogEntry{
		Message: "user authenticated successfully",
	}

	violations := rule.Check(entry)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %d", len(violations))
	}
}

func TestSensitiveRule_DetectsPasswordConcatenation(t *testing.T) {
	rule := NewSensitiveRule()

	entry := core.LogEntry{
		Message: "user password: " + "secret",
	}

	violations := rule.Check(entry)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}

	if violations[0].RuleID != "sensitive" {
		t.Errorf("RuleID = %q, want %q", violations[0].RuleID, "sensitive")
	}
}

func TestSensitiveRule_KnownKeywords(t *testing.T) {
	rule := NewSensitiveRule()

	entry := core.LogEntry{
		Message: "token: abc",
	}

	violations := rule.Check(entry)
	if len(violations) == 0 {
		t.Fatalf("expected violations for token keyword")
	}
}

func TestSensitiveRule_DetectsSensitiveIdentifierInExpr(t *testing.T) {
	rule := NewSensitiveRule()

	expr, err := parser.ParseExpr(`"password: " + password`)
	if err != nil {
		t.Fatalf("parse expr: %v", err)
	}

	entry := core.LogEntry{
		Message:     "",
		MessageExpr: expr,
	}

	violations := rule.Check(entry)
	if len(violations) == 0 {
		t.Fatalf("expected violations for sensitive identifier in expression")
	}
}

func TestSensitiveRule_CustomKeywords(t *testing.T) {
	rule := NewSensitiveRuleWithKeywords([]string{"mysecret"})

	entry := core.LogEntry{
		Message: "value mysecret=1",
	}

	violations := rule.Check(entry)
	if len(violations) == 0 {
		t.Fatalf("expected violations for custom keyword")
	}
}
