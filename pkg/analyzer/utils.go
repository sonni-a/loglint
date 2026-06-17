package analyzer

import (
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
)

var logLevels = map[string]bool{
	"Debug": true, "Info": true, "Warn": true, "Error": true,
	"DPanic": true, "Panic": true, "Fatal": true,
}

var knownLogPackages = map[string]bool{
	"log":      true,
	"log/slog": true,
}

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
	base := strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(name, "f"), "w"), "Context")
	return logLevels[base] || name == "Log" || name == "LogAttrs"
}

func isSupportedLogger(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	if ident, ok := sel.X.(*ast.Ident); ok {
		if obj, ok := pass.TypesInfo.Uses[ident].(*types.PkgName); ok {
			if knownLogPackages[obj.Imported().Path()] {
				return true
			}
		}
	}

	typ := pass.TypesInfo.TypeOf(sel.X)
	if typ != nil {
		typeStr := strings.ToLower(typ.String())
		return strings.Contains(typeStr, "logger") || strings.Contains(typeStr, "zap")
	}

	return false
}
