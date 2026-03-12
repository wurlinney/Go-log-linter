package inspectors

import (
	"go/ast"
	"go/token"
	"strings"
)

var defaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"secret",
	"token",
	"api_key",
	"apikey",
	"jwt",
	"credential",
	"cookie",
}

// ContainsSensitiveKeywordWithKeywords использует дефолтный список и дополнительные слова.
func ContainsSensitiveKeywordWithKeywords(s string, custom []string) bool {
	ls := strings.ToLower(s)
	keywords := append([]string{}, defaultSensitiveKeywords...)
	keywords = append(keywords, custom...)

	for _, k := range keywords {
		if strings.Contains(ls, k) {
			return true
		}
	}
	return false
}

// ExprContainsSensitiveIdentifierWithKeywords учитывает дополнительные слова.
func ExprContainsSensitiveIdentifierWithKeywords(expr ast.Expr, custom []string) bool {
	switch n := expr.(type) {
	case *ast.Ident:
		return isSensitiveNameWithKeywords(n.Name, custom)
	case *ast.BinaryExpr:
		if n.Op == token.ADD {
			return ExprContainsSensitiveIdentifierWithKeywords(n.X, custom) || ExprContainsSensitiveIdentifierWithKeywords(n.Y, custom)
		}
	case *ast.CallExpr:
		for _, a := range n.Args {
			if ExprContainsSensitiveIdentifierWithKeywords(a, custom) {
				return true
			}
		}
	case *ast.SelectorExpr:
		if ExprContainsSensitiveIdentifierWithKeywords(n.Sel, custom) {
			return true
		}
		return ExprContainsSensitiveIdentifierWithKeywords(n.X, custom)
	}
	return false
}

func isSensitiveNameWithKeywords(name string, custom []string) bool {
	lower := strings.ToLower(name)
	keywords := append([]string{}, defaultSensitiveKeywords...)
	keywords = append(keywords, custom...)

	for _, k := range keywords {
		if strings.Contains(lower, k) {
			return true
		}
	}
	for _, extra := range []string{"key", "api", "secret"} {
		if strings.Contains(lower, extra) {
			return true
		}
	}
	return false
}
