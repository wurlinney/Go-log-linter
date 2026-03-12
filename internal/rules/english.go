package rules

import (
	"unicode"

	"github.com/wurlinney/go-log-linter/internal/core"
)

type EnglishRule struct{}

func NewEnglishRule() *EnglishRule {
	return &EnglishRule{}
}

func (r *EnglishRule) ID() string {
	return "english"
}

// Check проверяет, что сообщение содержит только ASCII символы.
func (r *EnglishRule) Check(entry core.LogEntry) []core.Violation {
	if isEnglish(entry.Message) {
		return nil
	}

	return []core.Violation{
		{
			RuleID:  r.ID(),
			Message: "log message should contain only english characters",
			Pos:     entry.Pos,
			End:     entry.End,
		},
	}
}

func isEnglish(s string) bool {
	for _, r := range s {
		if r > unicode.MaxASCII {
			return false
		}
	}
	return true
}
