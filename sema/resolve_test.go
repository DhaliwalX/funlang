package sema

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/parse"
	"fmt"
	"testing"
)


type IdentPrinter struct {
}

func (printer *IdentPrinter) Visit(node ast.Node) ast.Visitor {
	if i, ok := node.(*ast.Identifier); ok {
		fmt.Printf("%s -> %s\n", i, fmt.Sprint(i.Object))
		return nil
	}

	return printer
}

func TestResolve(t *testing.T) {
	ctx := &context.Context{}
	parser := parse.NewParserFromString(ctx, `var x Struct;
type int int
type float float
type string string

type Struct struct {
	hello string
}

func Hello(a string) int {
	x.hello = a;
}

func some(a int, b float) string {
	var x = 10;
	
	return a+b-x;
}
`)
	a, err := parser.Parse()
	if err != nil {
		fmt.Print(err)
		return
	}

	t.Log(a)
	errs := ResolveProgram(a)
	if len(errs) != 0 {
		t.Error(fmt.Sprintf("len(errs) > 0:\n%s", errs))
	}

	printer := &IdentPrinter{}
	ast.Walk(printer, a)
}
