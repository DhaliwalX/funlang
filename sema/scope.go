package sema

import (
	"bitbucket.org/dhaliwalprince/funlang/types"
	"context"
)

type EntryType int

const (
	TYPE_NAME EntryType = iota
	VAR_NAME
)

type ScopeEntry struct {
	t EntryType
}

// scope represents one level
type Scope struct {
	// types declared at this scope
	factory *types.Factory

	// symbols declared at this scope
	symbols map[string]ScopeEntry
	parent *Scope
}

func NewScope(parent *Scope, ctx *context.Context) *Scope {
	return &Scope{parent:parent, symbols:make(map[string]ScopeEntry), factory: types.NewFactory(ctx)}
}
