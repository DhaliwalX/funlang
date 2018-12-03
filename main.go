package main

import (
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/parse"
	"fmt"
)

func main() {
	ctx := &context.Context{}
	parser := parse.NewParserFromString(ctx, `var a int = 10; var b = 10; var c = "string";`)
	ast, err := parser.Parse()
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Print(ast)
}
