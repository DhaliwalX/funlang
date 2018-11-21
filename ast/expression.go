package ast

import "bitbucket.org/dhaliwalprince/funlang/lex"

type exprNode interface {
    Node
    expr()
}

type NilLiteral struct {
    pos lex.Position // since nil has three characters only
}

type NumericLiteral struct {
    pos lex.Position
    val string
    isFloating bool
}

type StringLiteral struct {
    pos lex.Position
    val string
}

type ArrayLiteral struct {
    vals []exprNode
}

type StructProp struct {
    pos lex.Position
    name string
    val exprNode
}

type StructLiteral struct {
    pos lex.Position
    end lex.Position
    props []StructProp
}

type Identifier struct {
    pos lex.Position
    name string
}

type BooleanLiteral struct {
    pos lex.Position
    val bool
}

type ArgumentList struct {
    pos lex.Position
    exprs []exprNode
}

type MemberExpression struct {
    pos lex.Position
    token lex.TokenType
    member exprNode
    x exprNode
}

type PrefixExpression struct {
    pos lex.Position
    op lex.TokenType
    x exprNode
}

type PostfixExpression struct {
    pos lex.Position
    end lex.Position
    op lex.TokenType
    x exprNode
}

type BinaryExpression struct {
    pos lex.Position
    op lex.TokenType
    left exprNode
    right exprNode
}

type AssignExpression struct {
    pos lex.Position
    left exprNode
    right exprNode
}

func (*NilLiteral) expr() {}
func (*NumericLiteral) expr() {}
func (*StringLiteral) expr() {}
func (*BooleanLiteral) expr() {}
func (*Identifier) expr() {}
func (*ArgumentList) expr() {}
func (*MemberExpression) expr() {}
func (*PrefixExpression) expr() {}
func (*PostfixExpression) expr() {}
func (*BinaryExpression) expr() {}
func (*AssignExpression) expr() {}

// support for visitor
func (l *NilLiteral) Accept(visitor Visitor) {
    visitor.VisitNilLiteral(l)
}

func (n *NumericLiteral) Accept(visitor Visitor) {
    visitor.VisitNumericLiteral(n)
}

func (s *StringLiteral) Accept(visitor Visitor) {
    visitor.VisitStringLiteral(s)
}

func (b *BooleanLiteral) Accept(visitor Visitor) {
    visitor.VisitBooleanLiteral(b)
}

func (i *Identifier) Accept(visitor Visitor) {
    visitor.VisitIdentifier(i)
}

func (a *ArgumentList) Accept(visitor Visitor) {
    visitor.VisitArgumentList(a)
}

func (m *MemberExpression) Accept(visitor Visitor) {
    visitor.VisitMemberExpression(m)
}

func (p *PrefixExpression) Accept(visitor Visitor) {
    visitor.VisitPrefixExpression(p)
}

func (p *PostfixExpression) Accept(visitor Visitor) {
    visitor.VisitPostfixExpression(p)
}

func (b *BinaryExpression) Accept(visitor Visitor) {
    visitor.VisitBinaryExpression(b)
}

func (a *AssignExpression) Accept(visitor Visitor) {
    visitor.VisitAssignExpression(a)
}

func (n *NilLiteral) Beg() lex.Position { return n.pos }
func (n *NilLiteral) End() lex.Position {
    return lex.Position{Col:n.pos.Col+3,Row:n.pos.Row}
}

func (n *NumericLiteral) Beg() lex.Position{ return n.pos }
func (n *NumericLiteral) End() lex.Position {
    return lex.Position{Col: n.pos.Col+len(n.val), Row: n.pos.Row}
}

func (s *StringLiteral) Beg() lex.Position { return s.pos }
func (s *StringLiteral) End() lex.Position {
    return lex.Position{Col:s.pos.Col+len(s.val), Row:s.pos.Row}
}

func (b *BooleanLiteral) Beg() lex.Position { return b.pos }
func (b *BooleanLiteral) End() lex.Position {
    if b.val {
        return lex.Position{Col:b.pos.Col+4, Row:b.pos.Row}
    } else {
        return lex.Position{Col:b.pos.Col+5, Row:b.pos.Row}
    }
}

func (i *Identifier) Beg() lex.Position { return i.pos }
func (i *Identifier) End() lex.Position {
    return lex.Position{Col: i.pos.Col+len(i.name), Row: i.pos.Row}
}

func (a *ArgumentList) Beg() lex.Position { return a.pos }
func (a *ArgumentList) End() lex.Position {
    if len(a.exprs) > 0 {
        last := a.exprs[len(a.exprs)-1]
        return lex.Position{
            Row: last.End().Row,
            Col: last.End().Col,
        }
    } else {
        return a.Beg()
    }
}

func (m *MemberExpression) Beg() lex.Position { return m.pos }
func (m *MemberExpression) End() lex.Position { return m.x.End() }

func (p *PrefixExpression) Beg() lex.Position { return p.pos }
func (p *PrefixExpression) End() lex.Position { return p.x.End() }

func (p *PostfixExpression) Beg() lex.Position { return p.pos }
func (p *PostfixExpression) End() lex.Position {
    return p.end
}

func (b *BinaryExpression) Beg() lex.Position { return b.left.Beg() }
func (b *BinaryExpression) End() lex.Position { return b.right.End() }

func (a *AssignExpression) Beg() lex.Position { return a.left.Beg() }
func (a *AssignExpression) End() lex.Position { return a.right.End() }

// visitor methods
