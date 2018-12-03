package ast

import (
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"fmt"
	"strings"
)

type Statement interface {
	Node
	stmt()
}

type BlockStatement struct {
	pos lex.Position
	end lex.Position
	stmts []Statement
}

type ForStatement struct {
	pos       lex.Position
	init      Expression
	condition Expression
	body      Statement
}

type ExpressionStmt struct {
	expr Expression
}

type FunctionProtoType struct {
	pos lex.Position
	end lex.Position
	name string
	args []DeclNode
	ret Expression
}

type FunctionStatement struct {
	proto *FunctionProtoType
	// for linking c functions
	isExtern bool
	body *BlockStatement
}

type IfElseStatement struct {
	pos       lex.Position
	condition Expression
	body      Statement
	elseNode  Statement
}

type ReturnStatement struct {
	pos  lex.Position
	expr Expression
}

func (*BlockStatement) stmt() {}
func (*ForStatement) stmt() {}
func (*ExpressionStmt) stmt() {}
func (*FunctionProtoType) stmt() {}
func (*FunctionStatement) stmt() {}
func (*IfElseStatement) stmt() {}
func (*ReturnStatement) stmt() {}

func (b *BlockStatement) Accept(v Visitor) {
	v.VisitBlockStatement(b)
}

func (f *ForStatement) Accept(v Visitor) {
	v.VisitForStatement(f)
}

func (e *ExpressionStmt) Accept(v Visitor) {
	v.VisitExpressionStmt(e)
}

func (f *FunctionProtoType) Accept(v Visitor) {
	v.VisitFunctionProtoType(f)
}

func (f *FunctionStatement) Accept(v Visitor) {
	v.VisitFunctionStatement(f)
}

func (i *IfElseStatement) Accept(v Visitor) {
	v.VisitIfElseStatement(i)
}

func (r *ReturnStatement) Accept(v Visitor) {
	v.VisitReturnStatement(r)
}

func (b *BlockStatement) Beg() lex.Position { return b.pos }
func (b *BlockStatement) End() lex.Position { return b.end }

func (f *ForStatement) Beg() lex.Position { return f.pos }
func (f *ForStatement) End() lex.Position { return f.body.End() }

func (e *ExpressionStmt) Beg() lex.Position { return e.expr.Beg() }
func (e *ExpressionStmt) End() lex.Position { return e.expr.End() }

func (f *FunctionProtoType) Beg() lex.Position { return f.pos }
func (f *FunctionProtoType) End() lex.Position { return f.end }

func (f *FunctionStatement) Beg() lex.Position { return f.proto.Beg() }
func (f *FunctionStatement) End() lex.Position { return f.body.End() }

func (i *IfElseStatement) Beg() lex.Position { return i.pos }
func (i *IfElseStatement) End() lex.Position { return i.body.End() }

func (r *ReturnStatement) Beg() lex.Position { return r.pos }
func (r *ReturnStatement) End() lex.Position { return r.expr.End() }

func (b *BlockStatement) String() string {
	builder := strings.Builder{}
	builder.WriteString("{")
	for _, stmt := range b.stmts {
		builder.WriteString(fmt.Sprint(stmt))
	}
	builder.WriteString("}")
	return builder.String()
}

func (b *ExpressionStmt) String() string {
	return fmt.Sprintf("%s;", b.expr)
}

func (f *ForStatement) String() string {
	return fmt.Sprintf("for %s; %s %s", f.init, f.condition, f.body)
}

func (i *IfElseStatement) String() string {
	return fmt.Sprintf("if %s %s else %s", i.condition, i.body, i.elseNode)
}

func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;", r.expr)
}

func (f *FunctionStatement) String() string {
	return fmt.Sprintf("%s %s", f.proto, f.body)
}

func (f *FunctionProtoType) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("func %s(", f.name))
	for _, arg := range f.args {
		builder.WriteString(fmt.Sprint(arg))
		builder.WriteString(",")
	}

	builder.WriteString(fmt.Sprintf(") %s", f.ret))
	return builder.String()
}

func (b *BlockStatement) Len() int {
	return len(b.stmts)
}
