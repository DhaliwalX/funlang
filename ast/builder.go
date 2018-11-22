package ast

import (
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/lex"
)

type Builder struct {
	ctx context.Context
}

func NewBuilder(ctx context.Context) *Builder {
	return &Builder{ctx}
}

func (b *Builder) NewNilLiteral(pos lex.Position) *NilLiteral {
	return &NilLiteral{pos}
}

func (b *Builder) NewNumericLiteral(pos lex.Position, val string, isFloat bool) *NumericLiteral {
	return &NumericLiteral{pos: pos, val: val, isFloating:isFloat}
}

func (b *Builder) NewStringLiteral(pos lex.Position, val string) *StringLiteral {
	return &StringLiteral{pos: pos, val: val}
}

func (b *Builder) NewBooleanLiteral(pos lex.Position, val bool) *BooleanLiteral {
	return &BooleanLiteral{pos, val}
}

func (b *Builder) NewIdentifier(pos lex.Position, name string) *Identifier {
	return &Identifier{pos, name}
}

func (b *Builder) NewArgumentList(pos lex.Position, args []Expression) *ArgumentList {
	return &ArgumentList{pos, args}
}

func (b *Builder) NewMemberExpression(pos lex.Position,
	token lex.TokenType, member Expression, x Expression) *MemberExpression {
	return &MemberExpression{pos:pos,token:token,member:member,x:x}
}

func (b *Builder) NewPrefixExpression(pos lex.Position, op lex.TokenType,
	x Expression) *PrefixExpression {
	return &PrefixExpression{pos: pos, op:op, x:x}
}

func (b *Builder) NewPostfixExpression(pos lex.Position, op lex.TokenType,
	x Expression) *PostfixExpression {
	return &PostfixExpression{pos:pos, op:op,x:x}
}

func (b *Builder) NewBinaryExpression(pos lex.Position, left, right Expression) *BinaryExpression {
	return &BinaryExpression{pos:pos,left:left, right:right}
}

func (b *Builder) NewAssignExpression(pos lex.Position, left,right Expression) *AssignExpression {
	return &AssignExpression{pos:pos, left:left, right:right}
}

func (b *Builder) NewArrayType(pos lex.Position, size, t Expression) *ArrayType {
	return &ArrayType{pos:pos, size:size, t: t}
}

func (b *Builder) NewField(name, t Expression) *Field {
	return &Field{name:name, t:t}
}

func (b *Builder) NewStructType(pos lex.Position, fields []*Field) *StructType {
	return &StructType{pos:pos, fields:fields}
}

func (b *Builder) NewFuncType(pos lex.Position, params []Expression, ret Expression) *FuncType {
	return &FuncType{pos:pos, params: params, ret:ret}
}

func (b *Builder) NewDeclaration(pos lex.Position, name string, t TypeDeclaration,
	init Expression) *Declaration {
	return &Declaration{pos:pos, name:name, t:t, init:init}
}

func (b *Builder) NewDeclarationList(pos lex.Position, decls []Declaration) *DeclarationList {
	return &DeclarationList{pos:pos, decls:decls}
}

func (b *Builder) NewTypeDeclaration(pos lex.Position, name *Identifier, t Expression) *TypeDeclaration {
	return &TypeDeclaration{pos:pos, name:name, t:t}
}
