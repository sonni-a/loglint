package analyzer

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
)

func getMsgArgIndex(name string) int {
	switch {
	case name == "Log" || name == "LogAttrs":
		return 2
	case strings.HasSuffix(name, "Context"):
		return 1
	default:
		return 0
	}
}

func isLogMethod(name string) bool {
	levels := map[string]bool{
		"Debug": true, "Info": true, "Warn": true, "Error": true,
		"DPanic": true, "Panic": true, "Fatal": true,
	}
	base := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(name, "f"), "w"), "Context")
	return levels[base] || name == "Log" || name == "LogAttrs"
}

func isSupportedLogger(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	if ident, ok := sel.X.(*ast.Ident); ok {
		name := strings.ToLower(ident.Name)
		if name == "zaplogger" || name == "logger" || name == "log" || name == "slog" {
			return true
		}
	}
	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ != nil {
		typeStr := strings.ToLower(typ.String())
		return strings.Contains(typeStr, "logger") || strings.Contains(typeStr, "zap")
	}
	return false
}
