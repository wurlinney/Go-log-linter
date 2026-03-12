package extractors

import (
	"testing"

	"github.com/wurlinney/go-log-linter/internal/testutil"
)

func TestZapExtractor_L_GlobalLogger(t *testing.T) {
	src := `zap.L().Info("starting server")`

	call, pass := testutil.ParseFirstCall(t, src)
	ex := NewZapExtractor()

	if !ex.Match(call, pass) {
		t.Fatalf("expected Match to be true for zap.L().Info")
	}

	entry, ok := ex.Extract(call, pass)
	if !ok {
		t.Fatalf("expected Extract to succeed")
	}

	if entry.Logger != "zap" {
		t.Errorf("Logger = %q, want %q", entry.Logger, "zap")
	}
	if entry.Level != "Info" {
		t.Errorf("Level = %q, want %q", entry.Level, "Info")
	}
	if entry.Message != "starting server" {
		t.Errorf("Message = %q, want %q", entry.Message, "starting server")
	}
}

func TestZapExtractor_LoggerVar(t *testing.T) {
	src := `logger.Error("something failed")`

	call, pass := testutil.ParseFirstCall(t, src)
	ex := NewZapExtractor()

	if !ex.Match(call, pass) {
		t.Fatalf("expected Match to be true for logger.Error")
	}

	entry, ok := ex.Extract(call, pass)
	if !ok {
		t.Fatalf("expected Extract to succeed")
	}

	if entry.Logger != "zap" {
		t.Errorf("Logger = %q, want %q", entry.Logger, "zap")
	}
	if entry.Level != "Error" {
		t.Errorf("Level = %q, want %q", entry.Level, "Error")
	}
	if entry.Message != "something failed" {
		t.Errorf("Message = %q, want %q", entry.Message, "something failed")
	}
}

func TestZapExtractor_NonStringFirstArg(t *testing.T) {
	src := `zap.L().Info(123)`

	call, pass := testutil.ParseFirstCall(t, src)
	ex := NewZapExtractor()

	entry, ok := ex.Extract(call, pass)
	if !ok || entry == nil {
		t.Fatalf("expected Extract to succeed even for non-string first arg")
	}
	if entry.Message != "" {
		t.Fatalf("expected Message to be empty for non-string first arg")
	}
}
