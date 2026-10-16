package analyzer

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

type Config struct {
	SensitiveWords  string `json:"sensitive-words"`
	EnableSensitive *bool  `json:"enable-sensitive"`
	EnableStyle     *bool  `json:"enable-style"`
	EnableEnglish   *bool  `json:"enable-english"`
	EnableSpecial   *bool  `json:"enable-special"`
}

func init() {
	register.Plugin("loglint", New)
}

type LogLintPlugin struct {
	config Config
}

func New(settings any) (register.LinterPlugin, error) {
	cfg, err := register.DecodeSettings[Config](settings)
	if err != nil {
		return nil, err
	}

	return &LogLintPlugin{config: cfg}, nil
}

func (p *LogLintPlugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{NewAnalyzer(p.analyzerConfig())}, nil
}

func (p *LogLintPlugin) analyzerConfig() analyzerConfig {
	cfg := defaultConfig()
	if words := parseWords(p.config.SensitiveWords); len(words) > 0 {
		cfg.sensitiveWords = words
	}
	if p.config.EnableSensitive != nil {
		cfg.enableSensitive = *p.config.EnableSensitive
	}
	if p.config.EnableStyle != nil {
		cfg.enableStyle = *p.config.EnableStyle
	}
	if p.config.EnableEnglish != nil {
		cfg.enableEnglish = *p.config.EnableEnglish
	}
	if p.config.EnableSpecial != nil {
		cfg.enableSpecial = *p.config.EnableSpecial
	}
	return cfg
}

func (p *LogLintPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
