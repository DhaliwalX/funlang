package passes

import (
	"math/big"
)

type BlockSet struct {
	big.Int
}

func (b *BlockSet) add(id int) bool {
	if b.Bit(id) > 0 {
		return false
	}
	b.SetBit(&b.Int, id, 1)
	return true
}

func (b *BlockSet) take() int {
	for i, l := 0, b.BitLen(); i < l; i++ {
		if b.Bit(i) != 0 {
			b.SetBit(&b.Int, i, 0)
			return i
		}
	}

	return -1
}
