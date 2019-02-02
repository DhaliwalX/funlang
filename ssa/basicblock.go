package ssa

import (
	"fmt"
	"strings"

	"funlang/types"
)

type BasicBlock struct {
	valueWithUsers
	valueWithName
	First Instruction
	Last  Instruction

	Preds, Succs []*BasicBlock
	Parent       *Function
	Index        int
}

func (b *BasicBlock) AddSucc(s *BasicBlock) {
	b.Succs = append(b.Succs, s)
}

// TODO: optimize this
func (b *BasicBlock) AddPred(s *BasicBlock) {
	for _, b := range s.Preds {
		if b == s {
			return
		}
	}
	b.Preds = append(b.Preds, s)
}

func (b *BasicBlock) Instructions() []Instruction {
	var elements []Instruction
	for i := b.First; i != nil; i = i.Next() {
		elements = append(elements, i)
	}
	return elements
}

func (b *BasicBlock) PushFront(val Instruction) {
	if b.Last == nil || b.First == nil {
		b.appendInstr(val)
		return
	}
	val.Elem().Next = b.First.Elem()
	b.First.Elem().Prev = val.Elem()
	val.Elem().Prev = nil
	b.First = val
}

func (b *BasicBlock) AppendInstr(val Instruction) {
	b.appendInstr(val)
}

func (b *BasicBlock) appendInstr(val Instruction) {
	if b.Last == nil || b.First == nil {
		b.First = val
		b.Last = val
		return
	}
	b.Last.Elem().Next = val.Elem()
	val.Elem().Prev = b.Last.Elem()
	b.Last = val
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
	builder.WriteString(":\t\t\t" + fmt.Sprintf("%d <u:%d>", b.Index, len(b.Users())) + "\n")
	for i := b.First; i != nil; i = i.Next() {
		builder.WriteString("\t" + i.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func (b *BasicBlock) ShortString() string {
	return "$" + b.Name()
}

// need to optimize this
func (b *BasicBlock) Remove(i Instruction) {
	if i == b.First {
		b.First = i.Next()
		if b.First != nil {
			b.First.Elem().Prev = nil
		}
	} else if i == b.Last {
		b.Last = i.Prev()
		if b.Last != nil {
			b.Last.Elem().Next = nil
		}
	} else {
		i.Elem().Prev.Next = i.Elem().Next
		i.Elem().Next.Prev = i.Elem().Prev
	}
	i.Elem().Next = nil
	i.Elem().Prev = nil
}
