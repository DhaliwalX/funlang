package ssa

import "bitbucket.org/dhaliwalprince/funlang/ast"

// transform AST to SSA
type transformer struct {

}

func (t *transformer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	default:
		panic(n)
	}

	return nil
}
