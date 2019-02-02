package analysis

import (
	"fmt"
	"funlang/res"
	"funlang/ssa"
)

type DominatorAnalysisInfo struct {
	Util      *res.DomInfo
	Frontiers [][]int
	Graph     res.Graph
}

// this analysis will compute dominators information about the graph
type DominatorAnalysis struct {
	info *DominatorAnalysisInfo
}

func (d *DominatorAnalysis) GetInfo() AnalysisInfo {
	return d.info
}

func (d *DominatorAnalysis) IsAnalysisPass() bool {
	return true
}

func (d *DominatorAnalysis) Run(f *ssa.Function) bool {
	fmt.Println(f.Name(), "--")
	if len(f.Blocks) == 1 {
		return false
	}
	g := res.CreateGraph(f.Blocks)
	fmt.Print(g.Dot())
	info := res.ComputeDominators(g)
	d.info = &DominatorAnalysisInfo{Util: info}
	fmt.Println(d.info.Util)
	frontiers := res.ComputeDominanceFrontiers(g, d.info.Util)
	d.info.Frontiers = frontiers
	fmt.Println("Frontiers", d.info.Frontiers)
	return false
}

func init() {
	ssa.RegisterPass("dominators", &DominatorAnalysis{})
}
