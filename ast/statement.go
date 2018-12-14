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

func (b *BlockStatement) Statements() []Statement {
	return b.stmts
}

type ForStatement struct {
	pos       lex.Position
	init      Expression
	condition Expression
	body      Statement
}

func (f *ForStatement) Init() Expression {
	return f.init
}

func (f *ForStatement) Condition() Expression {
	return f.condition
}

func (f *ForStatement) Body() Statement {
	return f.body
}

type ExpressionStmt struct {
	expr Expression
}

func (e *ExpressionStmt) Expr() Expression {
	return e.expr
}

type FunctionProtoType struct {
	pos lex.Position
	end lex.Position
	name string
	args []DeclNode
	ret Expression
}

func (proto *FunctionProtoType) Name() string {
	return proto.name
}

func (proto *FunctionProtoType) Params() []DeclNode {
	return proto.args
}

func (proto *FunctionProtoType) Return() Expression {
	return proto.ret
}

type FunctionStatement struct {
	proto *FunctionProtoType
	// for linking c functions
	isExtern bool
	body *BlockStatement
}

func (f *FunctionStatement) Proto() *FunctionProtoType {
	return f.proto
}

func (f *FunctionStatement) Body() *BlockStatement {
	return f.body
}

type IfElseStatement struct {
	pos       lex.Position
	condition Expression
	body      Statement
	elseNode  Statement
}

func (i *IfElseStatement) Condition() Expression {
	return i.condition
}

func (i *IfElseStatement) Body() Statement {
	return i.body
}

func (i *IfElseStatement) ElseNode() Statement {
	return i.elseNode
}

type ReturnStatement struct {
	pos  lex.Position
	expr Expression
}

func (r *ReturnStatement) Expression() Expression {
	return r.expr
}

type DeclarationStatement struct {
	decl DeclNode
}

func (*BlockStatement) stmt() {}
func (*ForStatement) stmt() {}
func (*ExpressionStmt) stmt() {}
func (*FunctionProtoType) stmt() {}
func (*FunctionStatement) stmt() {}
func (*IfElseStatement) stmt() {}
func (*ReturnStatement) stmt() {}
func (*DeclarationStatement) stmt() {}

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

func (r *DeclarationStatement) Beg() lex.Position { return r.decl.Beg() }
func (r *DeclarationStatement) End() lex.Position { return r.decl.End() }

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

func (d *DeclarationStatement) String() string {
	return fmt.Sprintf("%s%s", d.decl, ";")
}

func (d *DeclarationStatement) Decl() DeclNode {
	return d.decl
}

func (b *BlockStatement) Len() int {
	return len(b.stmts)
}
