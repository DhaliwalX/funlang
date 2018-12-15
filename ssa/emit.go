package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"golang.org/x/tools/go/ssa"
)

// transform AST to SSA
type transformer struct {

}

func (t *transformer) transform(node ast.Node) ssa.Value

func (t *transformer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	default:
		panic(n)
	}

	return nil
}
