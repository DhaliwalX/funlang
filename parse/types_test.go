package parse

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"testing"
)

func newParser(source string) *Parser {
	source := lex.FileSource{source:source}
	return &Parser{ lex: lex.NewLexer(), errs: errorList{}, builder: ast.NewBuilder(context.Context{})  }
}

func TestParseType(t *testing.T) {

}
