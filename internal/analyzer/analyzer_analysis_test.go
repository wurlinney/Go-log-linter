package analyzer

import (
	"os"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func testdataDir(t *testing.T) string {
	t.Helper()
	td := analysistest.TestData()
	if _, err := os.Stat(td); os.IsNotExist(err) {
		t.Skip("analysistest testdata directory not present in this environment")
	}
	return td
}

func TestAnalyzer_LowercaseRule(t *testing.T) {
	t.Parallel()

	testdata := testdataDir(t)
	analysistest.Run(t, testdata, Analyzer, "lowercase")
}

func TestAnalyzer_EnglishRule(t *testing.T) {
	t.Parallel()

	testdata := testdataDir(t)
	analysistest.Run(t, testdata, Analyzer, "english")
}

func TestAnalyzer_SymbolsRule(t *testing.T) {
	t.Parallel()

	testdata := testdataDir(t)
	analysistest.Run(t, testdata, Analyzer, "symbols")
}

func TestAnalyzer_SensitiveRule(t *testing.T) {
	t.Parallel()

	testdata := testdataDir(t)
	analysistest.Run(t, testdata, Analyzer, "sensitive")
}

func TestAnalyzer_Mixed(t *testing.T) {
	t.Parallel()

	testdata := testdataDir(t)
	analysistest.Run(t, testdata, Analyzer, "mixed")
}
