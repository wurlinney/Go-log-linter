package golangci

import "testing"

func TestAnalyzers_Exposed(t *testing.T) {
	if len(Analyzers) == 0 {
		t.Fatalf("expected at least one analyzer to be exposed")
	}
	if Analyzers[0] == nil {
		t.Fatalf("first analyzer must not be nil")
	}
}

