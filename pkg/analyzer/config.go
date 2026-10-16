package analyzer

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

type analyzerConfig struct {
	sensitiveWords  []string
	enableSensitive bool
	enableStyle     bool
	enableEnglish   bool
	enableSpecial   bool
}

func defaultConfig() analyzerConfig {
	return analyzerConfig{
		sensitiveWords:  []string{"password", "token", "api_key", "secret"},
		enableSensitive: true,
		enableStyle:     true,
		enableEnglish:   true,
		enableSpecial:   true,
	}
}

func parseWords(s string) []string {
	words := strings.Split(s, ",")
	var result []string
	for _, w := range words {
		if w = strings.TrimSpace(w); w != "" {
			result = append(result, w)
		}
	}
	return result
}

func NewAnalyzer(cfg analyzerConfig) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "loglint",
		Doc:  "checks log messages for style and security rules",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return run(pass, cfg)
		},
	}
}

// Analyzer is the CLI entry point with flag-based configuration.
var Analyzer = newCLIAnalyzer()

func newCLIAnalyzer() *analysis.Analyzer {
	cfg := defaultConfig()
	words := strings.Join(cfg.sensitiveWords, ",")

	a := NewAnalyzer(cfg)

	a.Flags.BoolVar(&cfg.enableSensitive, "enable-sensitive", cfg.enableSensitive, "enable sensitive data check")
	a.Flags.StringVar(&words, "sensitive-words", words, "comma-separated list of sensitive words")
	a.Flags.BoolVar(&cfg.enableStyle, "enable-style", cfg.enableStyle, "enable lowercase check")
	a.Flags.BoolVar(&cfg.enableEnglish, "enable-english", cfg.enableEnglish, "enable english language check")
	a.Flags.BoolVar(&cfg.enableSpecial, "enable-special", cfg.enableSpecial, "enable special characters check")

	a.Run = func(pass *analysis.Pass) (interface{}, error) {
		local := cfg
		local.sensitiveWords = parseWords(words)
		return run(pass, local)
	}

	return a
}
