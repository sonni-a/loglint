package analyzer

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/sonni-a/loglint/pkg/rules"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks log messages for style and security rules",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			checkLogCall(pass, call)
			return true
		})
	}
	return nil, nil
}

func checkLogCall(pass *analysis.Pass, call *ast.CallExpr) {
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

	if enableSensitive {
		words := strings.Split(sensitiveWords, ",")
		var userWords []string
		for _, w := range words {
			if s := strings.TrimSpace(w); s != "" {
				userWords = append(userWords, s)
			}
		}

		checks = append(checks, func(s string) (string, string) {
			return rules.CheckSensitive(s, userWords)
		})
	}

	if enableEnglish {
		checks = append(checks, rules.CheckEnglish)
	}

	if enableSpecial {
		checks = append(checks, rules.CheckSpecialChars)
	}

	if enableStyle {
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
