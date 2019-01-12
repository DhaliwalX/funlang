package main

import (
	"bitbucket.org/dhaliwalprince/funlang/codegen"
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/parse"
	"bitbucket.org/dhaliwalprince/funlang/sema"
	"bitbucket.org/dhaliwalprince/funlang/ssa"
	"fmt"
	"os"
)

func main() {
	ctx := &context.Context{}
	if len(os.Args) < 2 {
		fmt.Println("usage: fun filename")
		return
	}

	filename := os.Args[1]
	p := parse.NewParserFromFile(ctx, os.Args[1])
	a, err := p.Parse()
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("== resolving variables")
	errs := sema.ResolveProgram(a)
	if len(errs) > 0 {
		fmt.Print(errs)
	}
	fmt.Println(a)

	program := ssa.Emit(a, ctx)
	fmt.Print(program)

	backend := &codegen.GoBackend{}
	backend.Run(program)
	o, err := os.OpenFile(fmt.Sprintf("%s.c", filename), os.O_CREATE|os.O_WRONLY, 0555)
	if err != nil {
		panic(err)
	}
	defer o.Close()
	_, err = o.Write([]byte(backend.String()))

	if err != nil {
		panic(err)
	}
}