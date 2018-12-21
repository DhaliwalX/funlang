package ssa

import "bitbucket.org/dhaliwalprince/funlang/types"

type ValueTag int

const (
	CONSTANT_INT ValueTag = iota
	CONSTANT_STRING
	ARGUMENT
	BASIC_BLOCK
	INSTRUCTION
	FUNCTION
)

// Value represents any constant, variable, method or instruction in ssa
type Value interface {
	// uses represent the list of values which this value uses
	Uses() []Value

	// users represent the list of values which uses this value
	Users() []Value

	// append a user to users list
	AddUser(user Value)
	Tag() ValueTag
	Type() types.Type
	Name() string
	// short representation
	ShortString() string
	// large representation
	String() string
}

type valueWithUsers struct {
	users []Value
}

func (i valueWithUsers) Users() []Value {
	return i.users
}

func (i valueWithUsers) AddUser(user Value) {
	i.users = append(i.users, user)
}

type valueWithNoName struct{}

func (v valueWithNoName) Name() string { return "" }

type valueWithName struct {
	name string
}

func (i valueWithName) Name() string {
	return i.name
}
