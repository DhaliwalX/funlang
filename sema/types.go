// type checking logic
package sema

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/types"
)

type typeChecker struct {
	types map[ast.Node]types.Type
	factory types.Factory
}
