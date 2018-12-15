package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/parse"
	"bitbucket.org/dhaliwalprince/funlang/sema"
	"bitbucket.org/dhaliwalprince/funlang/types"
	"fmt"
	"testing"
)

func newTestFunction() *Function {
	return &Function{current:&BasicBlock{}}
}

func TestEmitDecl(t *testing.T) {
	f := newTestFunction()
	ctx := &context.Context{}
	tr := transformer{
		function: f,
		factory: types.NewFactory(ctx),
		types: make(map[string]types.Type),
	}

	p := parse.NewParserFromString(ctx, "var a int = 10 + 20;")
	a, err := p.Parse()
	if err != nil {
		t.Error(err)
	}

	tr.Visit(a.Decls()[0])
	fmt.Print(tr.function.current)
}

func TestEmit(t *testing.T) {
	ctx := &context.Context{}
	p := parse.NewParserFromString(ctx, `type int int

type string string

type Person struct {
	name string
	age int
}

func add(a int, b int) int {
	return a + b;
}

func something(a int) int {
	var b = 20;
	b = b + a;
	return b - 40;
}

func Name(person Person) string {
	return person.name;
}

func AddAge(person Person, age int) int {
	person.age = person.age + age;
	return person.age;
}

func IfElse(person Person) int {
	var b int;
	if person.age > 100 {
		b = 20;
    } else {
		b = 10;
	}
	return b;
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

	program := Emit(a, ctx)
	t.Log(program)
}
