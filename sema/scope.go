package sema

import "bitbucket.org/dhaliwalprince/funlang/lex"

type ObjKind int

const (
	DONT_KNOW ObjKind = iota
	TYPE
	VAR
	FUNC
)


type Object struct {
	Kind ObjKind
	Name string
	Type interface{}
	Decl interface{}
	Func interface{}
	Pos lex.Position
}

type AlreadyDefined struct {}

func (AlreadyDefined) Error() string {
	return "already defined"
}

// scope is part of ast
type Scope struct {
	outer *Scope
	symbols map[string]*Object
}

func NewScope(outer *Scope) *Scope {
	return &Scope{ outer: outer, symbols: make(map[string]*Object)}
}

func (scope *Scope) Outer() *Scope {
	return scope.outer
}

func (scope *Scope) Lookup(name string) *Object {
	o, ok := scope.symbols[name]
	if !ok {
		return nil
	}
	return o
}

func (scope *Scope) Put(name string, o *Object) {
	scope.symbols[name] = o
}

func (scope *Scope) PutStrict(name string, o *Object) *Object {
	if k := scope.Lookup(name); k != nil {
		return k
	}

	scope.Put(name, o)
	return nil
}

func resolve(scope *Scope, name string) *Object {
	o := scope.Lookup(name)
	if o != nil {
		return o
	}

	if scope.outer != nil {
		return resolve(scope.outer, name)
	}

	return nil
}
