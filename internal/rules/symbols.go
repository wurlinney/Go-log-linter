package rules

import (
	"github.com/wurlinney/go-log-linter/internal/core"
	"github.com/wurlinney/go-log-linter/internal/inspectors"
)

type SymbolsRule struct{}

func NewSymbolsRule() *SymbolsRule {
	return &SymbolsRule{}
}

func (r *SymbolsRule) ID() string {
	return "symbols"
}

// Check отклоняет сообщения с неразрешёнными символами.
func (r *SymbolsRule) Check(entry core.LogEntry) []core.Violation {
	for _, ch := range entry.Message {
		if !inspectors.IsAllowedChar(ch) {
			return []core.Violation{
				{
					RuleID:  r.ID(),
					Message: "log message contains disallowed symbols",
					Pos:     entry.Pos,
					End:     entry.End,
				},
			}
		}
	}
	return nil
}
