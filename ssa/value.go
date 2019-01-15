package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/ds"
	"bitbucket.org/dhaliwalprince/funlang/types"
)

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
	RemoveFromUsers(v Value)

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

func (i *valueWithUsers) Users() []Value {
	return i.users
}

func (i *valueWithUsers) AddUser(user Value) {
	i.users = append(i.users, user)
}

func (i *valueWithUsers) RemoveFromUsers(v Value) {
	usersI := ds.RemoveFromSlice(ds.ToInterfaceSlice(i.users), v)
	users := make([]Value, len(usersI))
	for i, userI := range usersI {
		users[i] = userI.(Value)
	}
	i.users = users
}

type valueWithNoName struct{}

func (v *valueWithNoName) Name() string { return "" }

type valueWithName struct {
	name string
}

func (i *valueWithName) Name() string {
	return i.name
}
