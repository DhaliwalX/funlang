// mem2reg pass will try promote memory values to registers
package passes

import (
	"bitbucket.org/dhaliwalprince/funlang/ssa"
	"bitbucket.org/dhaliwalprince/funlang/ssa/analysis"
)

type Mem2RegPass struct {
	dom *analysis.DominatorAnalysisInfo
}

func (m *Mem2RegPass) IsAnalysisPass() bool {
	return false
}

func (m *Mem2RegPass) Run(f *ssa.Function) bool {
	// this expects that dominator analysis has been already run
	dominatorAnalysis := ssa.GetPass("dominator")
	m.dom = dominatorAnalysis.(*analysis.DominatorAnalysis).GetInfo().(*analysis.DominatorAnalysisInfo)
	changed := false


	return changed
}
