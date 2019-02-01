package passes

import (
	"math/big"

	"funlang/ssa"
)

type BlockSet struct {
	set *big.Int
}

func (b *BlockSet) Add(bb *ssa.BasicBlock) bool {
	if b.set.Bit(bb.Index) > 0 {
		return false
	}
	b.set.SetBit(b.set, bb.Index, 1)
	return true
}

func (b *BlockSet) Remove(bb *ssa.BasicBlock) {
	b.set.SetBit(b.set, bb.Index, 0)
}
