package core

import (
	"go/token"

	"golang.org/x/tools/go/analysis"
)

// Violation представляет найденное нарушение правила.
type Violation struct {
	RuleID  string
	Message string
	Pos     token.Pos
	End     token.Pos
	Fix     *analysis.SuggestedFix
}
