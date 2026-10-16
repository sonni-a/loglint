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
	obj := selectorObject(pass, sel)
	if !isLogMethod(funcName) || !isSupportedLogger(obj) {
		return
	}

	msgIdx := msgArgIndex(obj)
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

	type logCheck struct {
		apply func(string) (string, string)
		fix   bool
	}

	var checks []logCheck

	if cfg.enableSensitive {
		words := cfg.sensitiveWords
		checks = append(checks, logCheck{
			apply: func(s string) (string, string) {
				return rules.CheckSensitive(s, words)
			},
			fix: true,
		})
	}

	if cfg.enableEnglish {
		checks = append(checks, logCheck{apply: rules.CheckEnglish, fix: true})
	}

	if cfg.enableSpecial {
		checks = append(checks, logCheck{
			apply: rules.CheckSpecialChars,
			fix:   !isFormatMethod(funcName),
		})
	}

	if cfg.enableStyle {
		checks = append(checks, logCheck{apply: rules.CheckLowercase, fix: true})
	}

	var violations []string

	for _, check := range checks {
		errMsg, _ := check.apply(msg)
		if errMsg != "" {
			violations = append(violations, errMsg)
		}
	}

	combined := msg
	const maxFixPasses = 5
	for range maxFixPasses {
		next := combined
		for _, check := range checks {
			if !check.fix {
				continue
			}
			errMsg, fixedMsg := check.apply(next)
			if errMsg != "" && fixedMsg != "" && fixedMsg != next {
				next = fixedMsg
			}
		}
		if next == combined {
			break
		}
		combined = next
	}

	for i, errMsg := range violations {
		diag := analysis.Diagnostic{
			Pos:     lit.Pos(),
			End:     lit.End(),
			Message: "log message " + errMsg,
		}

		if i == len(violations)-1 && combined != msg {
			diag.SuggestedFixes = []analysis.SuggestedFix{{
				Message: "apply fix",
				TextEdits: []analysis.TextEdit{{
					Pos:     lit.Pos(),
					End:     lit.End(),
					NewText: []byte(strconv.Quote(combined)),
				}},
			}}
		}

		pass.Report(diag)
	}
}
