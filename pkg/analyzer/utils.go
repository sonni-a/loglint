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
	"log":               true,
	"log/slog":          true,
	"go.uber.org/zap":   true,
}

func logMethodBase(name string) string {
	for _, suffix := range []string{"f", "w", "ln", "Context"} {
		name = strings.TrimSuffix(name, suffix)
	}
	return name
}

func isLogMethod(name string) bool {
	if name == "Log" || name == "LogAttrs" {
		return true
	}
	base := logMethodBase(name)
	return logLevels[base] || base == "Print"
}

func isFormatMethod(name string) bool {
	if !strings.HasSuffix(name, "f") {
		return false
	}
	base := logMethodBase(name)
	return logLevels[base] || base == "Print"
}

func selectorObject(pass *analysis.Pass, sel *ast.SelectorExpr) types.Object {
	if obj := pass.TypesInfo.ObjectOf(sel.Sel); obj != nil {
		return obj
	}
	if seln := pass.TypesInfo.Selections[sel]; seln != nil {
		return seln.Obj()
	}
	return nil
}

func isSupportedLogger(obj types.Object) bool {
	if obj == nil {
		return false
	}
	if pkg := obj.Pkg(); pkg != nil && knownLogPackages[pkg.Path()] {
		return true
	}
	sig, ok := obj.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return false
	}
	return strings.Contains(strings.ToLower(namedTypeName(sig.Recv().Type())), "logger")
}

func namedTypeName(t types.Type) string {
	if p, ok := t.(*types.Pointer); ok {
		t = p.Elem()
	}
	named, ok := t.(*types.Named)
	if !ok {
		return ""
	}
	return named.Obj().Name()
}

func msgArgIndex(obj types.Object) int {
	sig, ok := obj.Type().(*types.Signature)
	if !ok {
		return 0
	}
	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		if isStringType(params.At(i).Type()) {
			return i
		}
	}
	return 0
}

func isStringType(t types.Type) bool {
	basic, ok := t.Underlying().(*types.Basic)
	return ok && basic.Kind() == types.String
}
