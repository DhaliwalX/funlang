package sema

import "funlang/ast"

type AlreadyDefined struct{}

func (AlreadyDefined) Error() string {
	return "already defined"
}

// scope is part of ast
type Scope struct {
	outer   *Scope
	symbols map[string]*ast.Object
}

func NewScope(outer *Scope) *Scope {
	return &Scope{outer: outer, symbols: make(map[string]*ast.Object)}
}

func (scope *Scope) Outer() *Scope {
	return scope.outer
}

func (scope *Scope) Lookup(name string) *ast.Object {
	o, ok := scope.symbols[name]
	if !ok {
		if scope.outer != nil {
			return scope.outer.Lookup(name)
		} else {
			return nil
		}
	}
	return o
}

func (scope *Scope) Put(name string, o *ast.Object) {
	scope.symbols[name] = o
}

func (scope *Scope) PutStrict(name string, o *ast.Object) *ast.Object {
	if k := scope.Lookup(name); k != nil {
		return k
	}

	scope.Put(name, o)
	return nil
}
