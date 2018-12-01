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

type DeclarationList struct {
	pos lex.Position
	decls []Declaration
}

type TypeDeclaration struct {
	pos lex.Position
	end lex.Position
	t Expression
	name Expression
}

func (*Declaration) decl() {}
func (*DeclarationList) decl() {}
func (*TypeDeclaration) decl() {}

func (d *Declaration) Accept(visitor Visitor) {
	visitor.VisitDeclaration(d)
}

func (d *DeclarationList) Accept(visitor Visitor) {
	visitor.VisitDeclarationList(d)
}

func (t *TypeDeclaration) Accept(visitor Visitor) {
	visitor.VisitTypeDeclaration(t)
}

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
		return fmt.Sprintf("var %s %s", d.name, d.init)
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
