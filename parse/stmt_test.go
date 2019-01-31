package parse

import (
	"funlang/ast"
	"testing"
)

func TestParseExpressionStmt(t *testing.T) {
	parser := newParser("a = 10;")
	parser.advance()
	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil")
	}
	if _, ok := stmt.(*ast.ExpressionStmt); !ok {
		t.Error("!ok", stmt)
	}

	t.Log(stmt)
}

func TestParseEmptyExpression(t *testing.T) {
	parser := newParser(";")
	parser.advance()
	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil")
	}

	if _, ok := stmt.(*ast.ExpressionStmt); !ok {
		t.Error("expected an expression statement", stmt)
	}

	t.Log(stmt)
}

func TestParseReturnStatement(t *testing.T) {
	parser := newParser("return expr;")
	parser.advance()
	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statment was nil")
	}

	if _, ok := stmt.(*ast.ReturnStatement); !ok {
		t.Error("expected a return statement", stmt)
	}
	t.Log(stmt)
}

func TestParseIfStatement(t *testing.T) {
	parser := newParser("if expr x = 10;")
	parser.advance()
	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil", parser.errs.Error())
	}

	if _, ok := stmt.(*ast.IfElseStatement); !ok {
		t.Error("expected an ifelsestatement", stmt)
	}

	t.Log(stmt)
}

func TestParseIfElseStatement(t *testing.T) {
	parser := newParser("if x dothis(); else dothat();")
	parser.advance()
	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil", parser.errs.Error())
	}

	if _, ok := stmt.(*ast.IfElseStatement); !ok {
		t.Error("expected an iufelsestatement", stmt)
	}

	t.Log(stmt)
}

func TestParseBlockStatement(t *testing.T) {
	parser := newParser("{ x = 10; y = &x; *x = 20; }")
	parser.advance()

	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil", parser.errs.Error())
	}

	if b, ok := stmt.(*ast.BlockStatement); ok {
		if b.Len() != 3 {
			t.Error("expected stmts to be 3 but got ", b.Len())
		}
	} else {
		t.Error("expected a block statement", stmt, parser.errs.Error())
	}

	t.Log(stmt)
}

func TestParseForStatement(t *testing.T) {
	parser := newParser("for x = 0; x < 10 { y = 20; x = x + 1; }")
	parser.advance()
	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil", parser.errs.Error())
	}

	if _, ok := stmt.(*ast.ForStatement); !ok {
		t.Error("expected a for statement")
	}

	t.Log(stmt)
}

func TestParseInfiniteForStatement(t *testing.T) {
	parser := newParser("for { x = 10; y > 10; }")
	parser.advance()
	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil", parser.errs.Error())
	}

	if _, ok := stmt.(*ast.ForStatement); !ok {
		t.Error("expected a for statement")
	}

	t.Log(stmt)
}

func TestParseForStatementWithConditionOnly(t *testing.T) {
	parser := newParser("for x > 10 { x = x + 1; }")
	parser.advance()

	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil", parser.errs.Error())
	}

	if _, ok := stmt.(*ast.ForStatement); !ok {
		t.Error("expected a for statement")
	}

	t.Log(stmt)
}

func TestParseFunction(t *testing.T) {
	parser := newParser("func print(arg int) int { return arg; }")
	parser.advance()

	stmt := parser.parseStatement()
	if stmt == nil {
		t.Error("statement was nil", parser.errs.Error())
	}

	if _, ok := stmt.(*ast.FunctionStatement); !ok {
		t.Error("expected a function")
	}

	t.Log(stmt)
}
