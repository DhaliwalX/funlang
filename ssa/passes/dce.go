package passes

import "bitbucket.org/dhaliwalprince/funlang/ssa"

// Dead Code Eliminator

type DCE struct {}

func (d *DCE) IsAnalysisPass() bool {
	return false
}

func (d *DCE) Run(b *ssa.BasicBlock) bool {
	branchSeen := false
	changed := false
	for i := b.First; i != nil; i = i.Next() {
		if branchSeen {
			ssa.Remove(i)
			changed = true
		} else if _, ok := i.(ssa.TerminatingInstr); ok {
			branchSeen = true
		}
	}

	die := ssa.GetPass("die")
	for die.(*DeadInstructionElimination).Run(b) {}
	return changed
}

func init() {
	ssa.RegisterPass("dce", &DCE{})
}
