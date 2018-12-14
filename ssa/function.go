package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/types"
	"fmt"
	"strings"
)

type Argument struct {
	valueWithName
	valueWithUsers
	t types.Type
}

func (a *Argument) Type() types.Type {
	return a.t
}

func (a *Argument) Tag() ValueTag {
	return ARGUMENT
}

func (a *Argument) Uses() []Value {
	return []Value{}
}

func (a *Argument) String() string {
	return fmt.Sprintf("%s:%s", a.Name(), a.Type())
}

func (a *Argument) ShortString() string {
	return a.String()
}

type Function struct {
	valueWithName
	valueWithUsers

	bbs []*BasicBlock
	t types.Type
	args []*Argument
}

func (f *Function) BasicBlocks() []*BasicBlock {
	return f.bbs
}

func (f *Function) Uses() []Value {
	return []Value{}
}

func (f *Function) Tag() ValueTag {
	return FUNCTION
}

func (f *Function) Type() types.Type {
	return f.t
}

func (f *Function) Arguments() []*Argument {
	return f.args
}

func (f *Function) ShortString() string {
	return fmt.Sprintf("%%%s", f.Name())
}

func (f *Function) String() string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("%s %s(", types.ToFunctionType(f.t).ReturnType(), f.Name()))
	l := len(f.args)
	for i, arg := range f.args {
		builder.WriteString(arg.String())
		if i+1 == l {
			break
		}
		builder.WriteString(", ")
	}

	builder.WriteString(") {\n")
	for _, bb := range f.bbs {
		builder.WriteString(bb.String())
		builder.WriteString("\n")
	}

	builder.WriteString("}\n")
	return builder.String()
}
