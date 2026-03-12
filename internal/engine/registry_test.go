package engine

import (
	"testing"

	"github.com/wurlinney/go-log-linter/internal/core"
)

// dummyRule нужен только чтобы убедиться, что NewDefaultEngine принимает правила.
type dummyRule struct{}

func (d *dummyRule) ID() string { return "dummy" }

func (d *dummyRule) Check(core.LogEntry) []core.Violation { return nil }

func TestNewDefaultEngine_RegistersExtractorsAndRules(t *testing.T) {
	rule := &dummyRule{}

	eng := NewDefaultEngine([]core.Rule{rule})

	if eng == nil {
		t.Fatalf("expected non-nil engine")
	}
	// Косвенно проверяем, что движок может быть использован, вызвав Analyze с nil-аргументами.
	vios := eng.Analyze(nil, nil)
	if len(vios) != 0 {
		t.Fatalf("expected 0 violations for nil input, got %d", len(vios))
	}
}
