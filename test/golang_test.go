package test

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
	"os"
	"testing"
)

func TestGolangSSA(t *testing.T) {
	fs := token.NewFileSet()
	a, err := parser.ParseFile(fs, "testfile.go", nil, 0)
	if err != nil {
		panic(err)
	}
	s := types.NewPackage("hello", "")
	hello, _, err := ssautil.BuildPackage(&types.Config{
		Importer: importer.Default(),
	}, fs, s, []*ast.File{a},  ssa.PrintFunctions|ssa.PrintPackages|ssa.LogSource)

	hello.WriteTo(os.Stdout)

	hello.Build()
	hello.WriteTo(os.Stdout)

	hello.Func("Pointer").WriteTo(os.Stdout)
}
