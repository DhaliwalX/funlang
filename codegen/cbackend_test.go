package codegen

import (
	"fmt"
	"funlang/context"
	"funlang/parse"
	"funlang/sema"
	"funlang/ssa"
	"testing"
)

func TestGoBackend(t *testing.T) {
	ctx := &context.Context{}
	p := parse.NewParserFromString(ctx, `type int int

func Add(a int, b int) int {
	return a + b;
}

func print(a int);

func max(a int, b int) int {
	if a > b {
		return a;
	}	return b;
}

func main() {
	var b int = 10;
	print(b);

	print(max(10, 20));
}
`)
	a, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	errs := sema.ResolveProgram(a)
	if len(errs) > 0 {
		t.Error(errs)
	}

	program := ssa.Emit(a, ctx)
	fmt.Print(program)

	backend := &GoBackend{}
	backend.Run(program)
	fmt.Print(backend)
}
