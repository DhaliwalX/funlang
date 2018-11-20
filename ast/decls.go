package ast

import (
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"bitbucket.org/dhaliwalprince/funlang/types"
)

type declNode interface {
	Node
	decl()
}

type Declaration struct {
	pos lex.Position
	name string
	t TypeDeclaration
	init exprNode
}

type DeclarationList struct {
	pos lex.Position
	decls []Declaration
}

type TypeDeclaration struct {
	pos lex.Position
	end lex.Position
	t types.Type
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
