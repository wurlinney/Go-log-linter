package rules

import (
	"go/ast"
	"go/token"
	"strconv"
	"unicode"

	"github.com/wurlinney/go-log-linter/internal/core"
	"github.com/wurlinney/go-log-linter/internal/inspectors"
	"golang.org/x/tools/go/analysis"
)

type LowercaseRule struct{}

func NewLowercaseRule() *LowercaseRule {
	return &LowercaseRule{}
}

func (r *LowercaseRule) ID() string {
	return "lowercase"
}

func (r *LowercaseRule) Check(entry core.LogEntry) []core.Violation {
	if inspectors.FirstLetterLower(entry.Message) {
		return nil
	}

	v := core.Violation{
		RuleID:  r.ID(),
		Message: "log message should start with lowercase letter",
		Pos:     entry.Pos,
		End:     entry.End,
	}

	// Попробуем построить SuggestedFix только если сообщение — строковый литерал.
	if lit, ok := entry.MessageExpr.(*ast.BasicLit); ok && lit.Kind == token.STRING {
		if fix := buildLowercaseFix(lit); fix != nil {
			v.Fix = fix
		}
	}

	return []core.Violation{v}
}

// buildLowercaseFix строит SuggestedFix для первой буквы.
func buildLowercaseFix(lit *ast.BasicLit) *analysis.SuggestedFix {
	unquoted, err := strconv.Unquote(lit.Value)
	if err != nil || unquoted == "" {
		return nil
	}

	runes := []rune(unquoted)
	applied := false

	for i, r := range runes {
		if unicode.IsLetter(r) {
			runes[i] = unicode.ToLower(r)
			applied = true
			break
		}
	}

	if !applied {
		return nil
	}

	newLiteral := strconv.Quote(string(runes))

	return &analysis.SuggestedFix{
		Message: "convert first letter to lowercase",
		TextEdits: []analysis.TextEdit{
			{
				Pos:     lit.Pos(),
				End:     lit.End(),
				NewText: []byte(newLiteral),
			},
		},
	}
}
