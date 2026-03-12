package rules

import (
	"github.com/wurlinney/go-log-linter/internal/core"
	"github.com/wurlinney/go-log-linter/internal/inspectors"
)

type SensitiveRule struct {
	customKeywords []string
}

func NewSensitiveRule() *SensitiveRule {
	return &SensitiveRule{}
}

// NewSensitiveRuleWithKeywords создаёт правило с дополнительными словами.
func NewSensitiveRuleWithKeywords(keywords []string) *SensitiveRule {
	return &SensitiveRule{
		customKeywords: keywords,
	}
}

func (r *SensitiveRule) ID() string {
	return "sensitive"
}

// Check ищет потенциально чувствительные данные в сообщении.
func (r *SensitiveRule) Check(entry core.LogEntry) []core.Violation {
	if inspectors.ContainsSensitiveKeywordWithKeywords(entry.Message, r.customKeywords) {
		return []core.Violation{
			{
				RuleID:  r.ID(),
				Message: "log message may contain sensitive data",
				Pos:     entry.Pos,
				End:     entry.End,
			},
		}
	}

	if entry.MessageExpr != nil && inspectors.ExprContainsSensitiveIdentifierWithKeywords(entry.MessageExpr, r.customKeywords) {
		return []core.Violation{
			{
				RuleID:  r.ID(),
				Message: "log message may contain sensitive data",
				Pos:     entry.Pos,
				End:     entry.End,
			},
		}
	}

	return nil
}
