package passes

import "bitbucket.org/dhaliwalprince/funlang/ssa"

// DeadInstructionElimination removes those instruction
// which are not used and does not cause any side effect
type DeadInstructionElimination struct {}

func (d *DeadInstructionElimination) Run(b *ssa.BasicBlock) bool {
	changed := false
	for _, i := range b.Instructions() {
		if ssa.IsSideEffect(i) {
			continue
		}
	}
	return changed
}

