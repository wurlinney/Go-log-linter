package report

import (
	"github.com/wurlinney/go-log-linter/internal/core"
	"golang.org/x/tools/go/analysis"
)

// ViolationToDiagnostic создаёт Diagnostic из Violation.
func ViolationToDiagnostic(v core.Violation) analysis.Diagnostic {
	diag := analysis.Diagnostic{
		Pos:     v.Pos,
		End:     v.End,
		Message: v.Message,
	}

	if v.Fix != nil {
		diag.SuggestedFixes = []analysis.SuggestedFix{*v.Fix}
	}

	return diag
}
