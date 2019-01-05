package ssa

import (
	"fmt"
	"strings"

	"bitbucket.org/dhaliwalprince/funlang/types"
)

type BasicBlock struct {
	valueWithUsers
	valueWithName
	instrs []Value

	Preds, Succs []*BasicBlock
	Parent       *Function
	index        int
}

func (b *BasicBlock) Instructions() []Value {
	return b.instrs
}

func (b *BasicBlock) appendInstr(val Value) {
	b.instrs = append(b.instrs, val)
}

func (b *BasicBlock) Uses() []Value {
	return []Value{}
}

func (b *BasicBlock) Tag() ValueTag {
	return BASIC_BLOCK
}

func (b *BasicBlock) Type() types.Type {
	return nil
}

func (b *BasicBlock) String() string {
	builder := strings.Builder{}
	builder.WriteString(b.Name())
	builder.WriteString(":\t\t\t" + fmt.Sprint(b.index) + "\n")
	for _, instr := range b.instrs {
		builder.WriteString("\t" + instr.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func (b *BasicBlock) ShortString() string {
	return "$" + b.Name()
}
