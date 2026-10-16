package analyzer_test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"testing"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"github.com/sonni-a/loglint/pkg/analyzer"
)

func TestPluginSettings(t *testing.T) {
	const src = `package p
import "log/slog"
func f() {
	slog.Info("Starting server!!!")
	slog.Info("user password: 123")
}
`
	styleAndSpecial := []string{
		"log message contains special characters",
		"log message must start with lowercase letter",
	}
	allDefaults := append(append([]string{}, styleAndSpecial...), "log message may contain sensitive data")

	tests := []struct {
		name     string
		settings any
		want     []string
		wantErr  bool
	}{
		{name: "nil", settings: nil, want: allDefaults},
		{name: "empty", settings: map[string]any{}, want: allDefaults},
		{name: "words only", settings: map[string]any{"sensitive-words": "credit_card"}, want: styleAndSpecial},
		{
			name:     "style off",
			settings: map[string]any{"enable-style": false},
			want: []string{
				"log message contains special characters",
				"log message may contain sensitive data",
			},
		},
		{
			name: "all off",
			settings: map[string]any{
				"enable-sensitive": false,
				"enable-style":     false,
				"enable-english":   false,
				"enable-special":   false,
			},
			want: nil,
		},
		{name: "unknown field", settings: map[string]any{"nope": true}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := analyzer.New(tt.settings)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			if p.GetLoadMode() != register.LoadModeTypesInfo {
				t.Fatalf("load mode %q", p.GetLoadMode())
			}
			got := pluginDiagnostics(t, p, src)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("diagnostics:\n got %q\nwant %q", got, tt.want)
			}
		})
	}
}

func pluginDiagnostics(t *testing.T, p register.LinterPlugin, src string) []string {
	t.Helper()
	analyzers, err := p.BuildAnalyzers()
	if err != nil {
		t.Fatalf("BuildAnalyzers: %v", err)
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	info := &types.Info{
		Types:      map[ast.Expr]types.TypeAndValue{},
		Uses:       map[*ast.Ident]types.Object{},
		Selections: map[*ast.SelectorExpr]*types.Selection{},
	}
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("p", fset, []*ast.File{file}, info)
	if err != nil {
		t.Fatal(err)
	}
	var msgs []string
	pass := &analysis.Pass{
		Fset:      fset,
		Files:     []*ast.File{file},
		Pkg:       pkg,
		TypesInfo: info,
		Report: func(d analysis.Diagnostic) {
			msgs = append(msgs, d.Message)
		},
	}
	if _, err := analyzers[0].Run(pass); err != nil {
		t.Fatal(err)
	}
	return msgs
}
