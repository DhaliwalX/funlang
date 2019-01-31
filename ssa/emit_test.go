package ssa

import (
	"fmt"
	"testing"

	"funlang/context"
	"funlang/parse"
	"funlang/sema"
)

func newTestFunction() *Function {
	return &Function{current: &BasicBlock{}}
}

func TestEmitDecl(t *testing.T) {
	// f := newTestFunction()
	// ctx := &context.Context{}
	// tr := transformer{
	// 	function: f,
	// 	factory: types.NewFactory(ctx),
	// 	types: make(map[string]types.Type),
	// }

	// p := parse.NewParserFromString(ctx, "var a int = 10 + 20;")
	// a, err := p.Parse()
	// if err != nil {
	// 	t.Error(err)
	// }

	// tr.Visit(a.Decls()[0])
	// fmt.Print(tr.function.current)
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

func TestFor(person Person) int {
	var a int = 0;
	for a = 10; a > 0 {
		person.age = person.age + a;
		a = a - 1;
	}
	
	return a;
}

//func TestIndex(arr []int, i *int) int {
//	return *i;
//}
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

func TestEmitCall(t *testing.T) {
	ctx := &context.Context{}
	p := parse.NewParserFromString(ctx, `type int int
func TestCall() int {
	return TestCall();
}

type Person struct { age int }

func TestPerson(person Person) Person {
	var age = TestCall();
	var personAge = TestPerson();
	return TestPerson();
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
	fmt.Print(program)
}

func TestLogical(t *testing.T) {
	ctx := &context.Context{}

	p := parse.NewParserFromString(ctx, `type int int
type bool int
func TestLogicalOperation(a bool, b bool) bool {
	return a && b;
}

func TestLogicalOperation2(a bool, b bool) bool {
	return a || b;
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
	fmt.Print(program)
}

func TestEmit2(t *testing.T) {
	ctx := &context.Context{}
	p := parse.NewParserFromString(ctx, `
	type int int
	//
	//func Sort(a []int, size int) int {
	//	var i = 0;
	//	var min = 0;
	//	for i < size {
	//		min = FindMinimum(a, i, size);
	//		swap(a, i, min);
	//		i = i + 1;
	//	}
	//}

	// find index of the minimum element in this array
	func FindMinimum(a []int, s int, size int) int {
		var i int = s;
		var min = 0;
		for i < size {
			if a[min] > a[i] {
				min = i;
			}
			i = i + 1;
		}

		return min;
	}

	func Swap(a []int, x int, y int) int {
		var t = a[x];
		a[x] = a[y];
		a[y] = t;
	}

	`)
	a, err := p.Parse()
	if err != nil {
		t.Error(err)
	}
	fmt.Print(a.String())

	errs := sema.ResolveProgram(a)

	if len(errs) > 0 {
		t.Error(errs)
	}

	program := Emit(a, ctx)
	fmt.Print(program)
}
