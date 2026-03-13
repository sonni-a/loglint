package analyzer_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/sonni-a/loglint/pkg/analyzer"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer.Analyzer, "a")
	analysistest.RunWithSuggestedFixes(t, testdata, analyzer.Analyzer, "b")
}
