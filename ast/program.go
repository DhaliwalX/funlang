package ast

import (
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"fmt"
	"strings"
)

type Program struct {
	source lex.Source
	ctx *context.Context
	decls []Node
}

func NewProgram(ctx *context.Context, source lex.Source, decls []Node) *Program {
	return &Program{source:source, ctx:ctx, decls:decls}
}

func (p *Program) Decls() []Node {
	return p.decls
}

func (p *Program) String() string {
	builder := strings.Builder{}
	for _, decl := range p.decls {
		builder.WriteString(fmt.Sprint(decl))
		builder.WriteString("\n")
	}

	return builder.String()
}
