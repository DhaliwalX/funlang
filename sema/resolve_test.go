package sema

import (
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/parse"
	"fmt"
	"testing"
)

func TestResolve(t *testing.T) {
	ctx := &context.Context{}
	parser := parse.NewParserFromString(ctx, `
var a int = 10;
var a = 10;
var c = "string";
var x Struct;
type Struct struct {
	hello int
}
`)
	ast, err := parser.Parse()
	if err != nil {
		fmt.Print(err)
		return
	}

	r := resolver{}
	r.openScope()
	errs := ResolveProgram(ast)
	if len(errs) != 0 {
		t.Error(fmt.Sprintf("len(errs) > 0:\n%s", errs))
	}
}
