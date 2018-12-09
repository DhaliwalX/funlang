package sema

import (
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/parse"
	"fmt"
	"testing"
)

func TestResolve(t *testing.T) {
	ctx := &context.Context{}
	parser := parse.NewParserFromString(ctx, `var x Struct;
type int int
type float float
type string string

type Struct struct {
	hello string
}

func Hello(a string) {
	x.hello = a;
}

func some(a int, b float) string {
	return a+b;
}
`)
	ast, err := parser.Parse()
	if err != nil {
		fmt.Print(err)
		return
	}

	t.Log(ast)
	errs := ResolveProgram(ast)
	if len(errs) != 0 {
		t.Error(fmt.Sprintf("len(errs) > 0:\n%s", errs))
	}
}
