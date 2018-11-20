package ast

import (
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"bitbucket.org/dhaliwalprince/funlang/types"
)

type stmtNode interface {
	Node
	stmt()
}

type BlockStatement struct {
	pos lex.Position
	end lex.Position
	stmts []stmtNode
}

type ForStatement struct {
	pos lex.Position
	init exprNode
	condition exprNode
	body stmtNode
}

type ExpressionStmt struct {
	expr exprNode
}

type FunctionProtoType struct {
	pos lex.Position
	end lex.Position
	name string
	args DeclarationList
	t types.Type
}

type FunctionStatement struct {
	proto FunctionProtoType
	// for linking c functions
	isExtern bool
	body BlockStatement
}

type IfElseStatement struct {
	pos lex.Position
	condition exprNode
	body stmtNode
	elseNode stmtNode
}

type ReturnStatement struct {
	pos lex.Position
	expr exprNode
}

func (*BlockStatement) stmt() {}
func (*ForStatement) stmt() {}
func (*ExpressionStmt) stmt() {}
func (*FunctionProtoType) stmt() {}
func (*FunctionStatement) stmt() {}
func (*IfElseStatement) stmt() {}
func (*ReturnStatement) stmt() {}

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
