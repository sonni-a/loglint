package analyzer

var (
	sensitiveWords  string
	enableSensitive bool
	enableStyle     bool
	enableEnglish   bool
	enableSpecial   bool
)

func init() {
	Analyzer.Flags.BoolVar(&enableSensitive, "enable-sensitive", true, "enable sensitive data check")
	Analyzer.Flags.StringVar(&sensitiveWords, "sensitive-words", "password,token,api_key,secret", "comma-separated list of sensitive words")
	Analyzer.Flags.BoolVar(&enableStyle, "enable-style", true, "enable lowercase check")
	Analyzer.Flags.BoolVar(&enableEnglish, "enable-english", true, "enable english language check")
	Analyzer.Flags.BoolVar(&enableSpecial, "enable-special", true, "enable special characters check")
}
