package ast

import (
    "bitbucket.org/dhaliwalprince/funlang/lex"
    "fmt"
    "strings"
)

type Expression interface {
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
    vals []Expression
}

type StructProp struct {
    pos  lex.Position
    name string
    val  Expression
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
    exprs []Expression
}

func (list *ArgumentList) Exprs() []Expression {
    return list.exprs
}

type MemberExpression struct {
    pos    lex.Position
    token  lex.TokenType
    member Expression
    x      Expression
}

type PrefixExpression struct {
    pos lex.Position
    op  lex.TokenType
    x   Expression
}

func (p *PrefixExpression) Op() lex.TokenType {
    return p.op
}

func (p *PrefixExpression) Expression() Expression {
    return p.x
}

type PostfixExpression struct {
    pos lex.Position
    end lex.Position
    op  lex.TokenType
    x   Expression
}

func (p *PostfixExpression) Op() lex.TokenType {
    return p.op
}

func (p *PostfixExpression) Expression() Expression {
    return p.x
}

type BinaryExpression struct {
    pos   lex.Position
    op    lex.TokenType
    left  Expression
    right Expression
}

type AssignExpression struct {
    pos   lex.Position
    left  Expression
    right Expression
}

// type tree
type ArrayType struct {
    pos lex.Position
    size Expression
    t Expression
}

func (a *ArrayType) Type() Expression {
    return a.t
}

func (a *ArrayType) Size() Expression {
    return a.size
}

type Field struct {
    name Expression
    t Expression
}

func (f *Field) Name() string {
    return f.name.(*Identifier).Name()
}

func (f *Field) Type() Expression {
    return f.t
}

type StructType struct {
    pos lex.Position
    fields []*Field
}

func (s *StructType) Fields() []*Field {
    return s.fields
}

type FuncType struct {
    pos lex.Position
    params []Expression
    ret Expression
}

func (f *FuncType) Params() []Expression {
    return f.params
}

func (f *FuncType) Return() Expression {
    return f.ret
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
func (*ArrayType) expr() {}
func (*Field) expr() {}
func (*StructType) expr() {}
func (*FuncType) expr() {}

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
        return last.End()
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

func (a *ArrayType) Beg() lex.Position { return a.pos }
func (a *ArrayType) End() lex.Position { return a.t.End() }

func (f *Field) Beg() lex.Position { return f.name.Beg() }
func (f *Field) End() lex.Position { return f.t.End() }

func (s *StructType) Beg() lex.Position { return s.pos }
func (s *StructType) End() lex.Position {
    if len(s.fields) >0 {
        last := s.fields[len(s.fields)-1]
        return last.End()
    } else {
        return s.pos
    }
}

func (f *FuncType) Beg() lex.Position { return f.pos }
func (f *FuncType) End() lex.Position { return f.ret.End() }

func (*NilLiteral) String() string {
    return "nil"
}

func (n *NumericLiteral) String() string {
    return n.val
}

func (n *StringLiteral) String() string {
    return n.val
}

func (b *BooleanLiteral) String() string {
    return fmt.Sprint(b.val)
}

func (i *Identifier) String() string {
    return i.name
}

func (a *ArgumentList) String() string {
    builder := strings.Builder{}
    builder.WriteString("(")
    for _, arg := range a.exprs {
        builder.WriteString(fmt.Sprint(arg))
        builder.WriteString(", ")
    }
    builder.WriteString(")")
    return builder.String()
}

func (m *MemberExpression) String() string {
    builder := strings.Builder{}
    builder.WriteString(fmt.Sprintf("('%s' %s %s)", m.token, m.member, m.x))
    return builder.String()
}

func (p *PrefixExpression) String() string {
    return fmt.Sprintf("%s (%s)", p.op, p.x)
}

func (p *BinaryExpression) String() string {
    return fmt.Sprintf("(%s %s %s)", p.left, p.op, p.right)
}

func (p *AssignExpression) String() string {
    return fmt.Sprintf("%s = %s", p.left, p.right)
}

func (a *ArrayType) String() string {
    return fmt.Sprintf("[]%s", a.t)
}

func (s *StructType) String() string {
    builder := strings.Builder{}
    builder.WriteString("struct{ ")
    for _, field := range s.fields {
        builder.WriteString(fmt.Sprintf("%s %s, ", field.name, field.t))
    }

    builder.WriteString("}")
    return builder.String()
}

func (a *ArgumentList) List() []Expression {
    return a.exprs
}

func (m *MemberExpression) Member() Expression {
    return m.x
}

func (m *MemberExpression) Expr() Expression {
    return m.member
}

func (m *MemberExpression) AccessKind() lex.TokenType {
    return m.token
}

func (i *Identifier) Name() string {
    return i.name
}

func (b *BinaryExpression) Left() Expression {
    return b.left
}

func (b *BinaryExpression) Right() Expression {
    return b.right
}


func (a *AssignExpression) Left() Expression {
    return a.left
}

func (a *AssignExpression) Right() Expression {
    return a.right
}