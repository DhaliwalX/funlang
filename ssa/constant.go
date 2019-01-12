package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/types"
	"fmt"
)

type ConstantInt struct {
	valueWithUsers
	Value int
}

func (c *ConstantInt) Name() string {
	return fmt.Sprint(c.Value)
}

func (c *ConstantInt) Uses() []Value {
	return []Value{}
}

func (c *ConstantInt) String() string {
	return fmt.Sprintf("%d", c.Value)
}

func (c *ConstantInt) ShortString() string {
	return c.String()
}

func (c *ConstantInt) Tag() ValueTag {
	return CONSTANT_INT
}

func (c *ConstantInt) Type() types.Type {
	return typeFactory.IntType()
}

type ConstantString struct {
	valueWithUsers
	Value string
}

func (c *ConstantString) Name() string {
	return "\"" + c.Value + "\""
}

func (c *ConstantString) String() string {
	return c.Value
}

func (c *ConstantString) ShortString() string {
	return c.String()
}

func (c *ConstantString) Uses() []Value {
	return []Value{}
}

func (c *ConstantString) Tag() ValueTag {
	return CONSTANT_STRING
}

func (c *ConstantString) Type() types.Type {
	return typeFactory.StringType()
}
