package parse

import (
	"fmt"
	"funlang/ast"
	"testing"
)

func TestParseExpression(t *testing.T) {
	parser := newParser("10+30")
	parser.advance()
	a := parser.parseExpression()
	if len(parser.errs.list) > 0 {
		t.Error("errors were not expected in this case", parser.errs.Error())
	}
	if _, ok := a.(*ast.BinaryExpression); !ok {
		t.Error("parsing failed for binary expression")
	}
}

func TestParsePrecendence(t *testing.T) {
	parser := newParser("10*30+50")
	parser.advance()
	a := parser.parseExpression()
	if len(parser.errs.list) > 0 {
		t.Error("errs were not expected in this case", parser.errs.Error())
	}

	if _, ok := a.(*ast.BinaryExpression); !ok {
		t.Error("parsing failed for binary expression")
	}

	if fmt.Sprint(a) != "((10 * 30) + 50)" {
		t.Error("parsing failed")
	}
}

func TestParseMemberExpression(t *testing.T) {
	parser := newParser("x(10)")
	parser.advance()
	fmt.Println(parser.parseExpression())
}

func TestParseExpression2(t *testing.T) {
	parser := newParser("(10+32)*a")
	parser.advance()
	a := parser.parseExpression()
	if len(parser.errs.list) > 0 {
		t.Error("errs were not expected in this case", parser.errs.Error())
	}

	if _, ok := a.(*ast.BinaryExpression); !ok {
		t.Error("expected a binary expression")
	}

	t.Log(a)
}

func TestParseAssignExpression(t *testing.T) {
	parser := newParser("a = 10 + 20;")
	parser.advance()
	a := parser.parseExpression()
	if len(parser.errs.list) > 0 {
		t.Error("errs were not expected in this case", parser.errs.Error())
	}

	if _, ok := a.(*ast.AssignExpression); !ok {
		t.Error("expected a assign expression")
	}

	t.Log(a)
}
