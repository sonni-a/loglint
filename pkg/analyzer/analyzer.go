package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"

	"github.com/sonni-a/loglint/pkg/rules"
	"golang.org/x/tools/go/analysis"
)

func run(pass *analysis.Pass, cfg analyzerConfig) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			checkLogCall(pass, call, cfg)
			return true
		})
	}
	return nil, nil
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr, cfg analyzerConfig) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	funcName := sel.Sel.Name
	if !isLogMethod(funcName) || !isSupportedLogger(pass, sel) {
		return
	}

	msgIdx := getMsgArgIndex(funcName)
	if len(call.Args) <= msgIdx {
		return
	}

	lit, ok := call.Args[msgIdx].(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return
	}

	msg, err := strconv.Unquote(lit.Value)
	if err != nil {
		return
	}

	var checks []func(string) (string, string)

	if cfg.enableSensitive {
		words := cfg.sensitiveWords
		checks = append(checks, func(s string) (string, string) {
			return rules.CheckSensitive(s, words)
		})
	}

	if cfg.enableEnglish {
		checks = append(checks, rules.CheckEnglish)
	}

	if cfg.enableSpecial {
		checks = append(checks, rules.CheckSpecialChars)
	}

	if cfg.enableStyle {
		checks = append(checks, rules.CheckLowercase)
	}

	for _, check := range checks {
		errMsg, fixedMsg := check(msg)
		if errMsg != "" {
			diag := analysis.Diagnostic{
				Pos:     lit.Pos(),
				End:     lit.End(),
				Message: "log message " + errMsg,
			}

			if fixedMsg != msg {
				diag.SuggestedFixes = []analysis.SuggestedFix{{
					Message: "apply fix",
					TextEdits: []analysis.TextEdit{{
						Pos:     lit.Pos(),
						End:     lit.End(),
						NewText: []byte(strconv.Quote(fixedMsg)),
					}},
				}}
			}
			pass.Report(diag)
			return
		}
	}
}
