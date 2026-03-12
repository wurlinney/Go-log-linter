package extractors

import (
	"testing"

	"github.com/wurlinney/go-log-linter/internal/testutil"
)

func TestSlogExtractor_MatchAndExtract_BasicLiteral(t *testing.T) {
	src := `slog.Info("starting server")`

	call, pass := testutil.ParseFirstCall(t, src)

	ex := NewSlogExtractor()

	if !ex.Match(call, pass) {
		t.Fatalf("expected Match to be true for slog.Info")
	}

	entry, ok := ex.Extract(call, pass)
	if !ok {
		t.Fatalf("expected Extract to succeed")
	}

	if entry.Logger != "slog" {
		t.Errorf("Logger = %q, want %q", entry.Logger, "slog")
	}
	if entry.Level != "Info" {
		t.Errorf("Level = %q, want %q", entry.Level, "Info")
	}
	if entry.Message != "starting server" {
		t.Errorf("Message = %q, want %q", entry.Message, "starting server")
	}
	if entry.Call == nil {
		t.Errorf("Call must not be nil")
	}
}

func TestSlogExtractor_Extract_NonStringFirstArg(t *testing.T) {
	src := `slog.Info(123)`

	call, pass := testutil.ParseFirstCall(t, src)
	ex := NewSlogExtractor()

	entry, ok := ex.Extract(call, pass)
	if !ok || entry == nil {
		t.Fatalf("expected Extract to succeed even for non-string first arg")
	}
	if entry.Message != "" {
		t.Fatalf("expected Message to be empty for non-string first arg")
	}
}
