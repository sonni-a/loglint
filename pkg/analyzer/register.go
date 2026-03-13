package analyzer

import (
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

type Config struct {
	SensitiveWords  string `json:"sensitive-words"`
	EnableSensitive bool   `json:"enable-sensitive"`
	EnableStyle     bool   `json:"enable-style"`
	EnableEnglish   bool   `json:"enable-english"`
	EnableSpecial   bool   `json:"enable-special"`
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
	sensitiveWords = p.config.SensitiveWords
	enableSensitive = p.config.EnableSensitive
	enableStyle = p.config.EnableStyle
	enableEnglish = p.config.EnableEnglish
	enableSpecial = p.config.EnableSpecial

	return []*analysis.Analyzer{Analyzer}, nil
}

func (p *LogLintPlugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
