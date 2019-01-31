package parse

import (
	"funlang/ast"
	"funlang/context"
	"funlang/lex"
	"testing"
)

func newParser(source string) *Parser {
	stringSource := lex.NewStringSource(source)
	p := &Parser{lex: lex.NewLexer(stringSource), errs: errorList{}, builder: ast.NewBuilder(&context.Context{})}
	return p
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

func TestParseTypeDeclaration(t *testing.T) {
	parser := newParser("type Name struct { val string }")
	parser.advance()
	a := parser.parseTypeDeclaration()
	if len(parser.errs.list) != 0 {
		t.Error(parser.errs.Error())
	}
	if a == nil {
		t.Error("nil type declaration")
	}

	t.Log(a)
}
