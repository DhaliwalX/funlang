package passes

import "bitbucket.org/dhaliwalprince/funlang/ssa"

// DeadInstructionElimination removes those instruction
// which are not used and does not cause any side effect
type DeadInstructionElimination struct {}

func canRemove(i ssa.Instruction) bool {
	if ssa.IsSideEffect(i) {
		return false
	}

	if len(i.Users()) == 0 {
		return true
	}

	return false
}

func (d *DeadInstructionElimination) IsAnalysisPass() bool {
	return false
}

func (d *DeadInstructionElimination) Run(b *ssa.BasicBlock) bool {
	changed := false
	for _, i := range b.Instructions() {
		if canRemove(i) {
			changed = true
			ssa.Remove(i)
		}
	}
	return changed
}

func init() {
	ssa.RegisterPass("die", &DeadInstructionElimination{})
}

