package parse

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"testing"
)

func newParser(source string) *Parser {
	stringSource := lex.NewStringSource(source)
	return &Parser{ lex: lex.NewLexer(stringSource), errs: errorList{}, builder: ast.NewBuilder(context.Context{})  }
}

func TestParseType(t *testing.T) {
	parser := newParser(" int")
	parser.advance()
	a := parser.parseType()
	if len(parser.errs.list) != 0 {
		t.Error(parser.errs.Error())
	}
	if _, ok := a.(*ast.Identifier); !ok {
		t.Error("did not parse identifier")
	}
}

func TestParseEmptyStructType(t *testing.T) {
	parser := newParser("struct{}")
	parser.advance()
	a := parser.parseType()
	if len(parser.errs.list) != 0 {
		t.Error(parser.errs.Error())
	}
	if _, ok := a.(*ast.StructType); !ok {
		t.Error("did not parse struct")
	}

	t.Log(a)
}

func TestParseStructType(t *testing.T) {
	parser := newParser("struct{a int b string}")
	parser.advance()
	a := parser.parseType()
	if len(parser.errs.list) != 0 {
		t.Error(parser.errs.Error())
	}
	if _, ok := a.(*ast.StructType); !ok {
		t.Error("did not parse struct")
	}

	t.Log(a)
}
