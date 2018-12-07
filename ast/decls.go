package ast

import (
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"fmt"
)

type DeclNode interface {
	Node
	decl()
}

type Declaration struct {
	pos  lex.Position
	name string
	t    Expression
	init Expression
}

func (d *Declaration) Name() string {
	return d.name
}

func (d *Declaration) Type() Expression {
	return d.t
}

func (d *Declaration) Init() Expression {
	return d.init
}

type DeclarationList struct {
	pos lex.Position
	decls []*Declaration
}

type TypeDeclaration struct {
	pos lex.Position
	end lex.Position
	t Expression
	name Expression
}

func (t *TypeDeclaration) Type() Expression {
	return t.t
}

func (t *TypeDeclaration) Name() string {
	return t.name.(*Identifier).Name()
}

func (*Declaration) decl() {}
func (*DeclarationList) decl() {}
func (*TypeDeclaration) decl() {}

func (d *DeclarationList) Beg() lex.Position { return d.pos }
func (d *DeclarationList) End() lex.Position {
	if len(d.decls) > 0 {
		last := d.decls[len(d.decls)-1]
		return lex.Position{
			Row: last.pos.Row,
			Col: last.pos.Col,
		}
	}

	return d.Beg()
}

func (d *Declaration) Beg() lex.Position { return d.pos }
func (d *Declaration) End() lex.Position { return d.t.End() }

func (t *TypeDeclaration) Beg() lex.Position { return t.pos }
func (t *TypeDeclaration) End() lex.Position { return t.end }

func (d *Declaration) String() string {
	if d.init == nil {
		return fmt.Sprintf("var %s %s", d.name, d.t)
	} else {
		if d.t == nil {
			return fmt.Sprintf("var %s = %s", d.name, d.init)
		} else {
			return fmt.Sprintf("var %s %s = %s", d.name, d.t, d.init)
		}
	}
}

func (t *TypeDeclaration) String() string {
	return fmt.Sprintf("type %s %s", t.name, t.t)
}
