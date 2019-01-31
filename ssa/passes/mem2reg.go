// mem2reg pass will try promote memory values to registers
package passes

import (
	"fmt"
	"funlang/ssa"
	"funlang/ssa/analysis"
)

var debug = true

type allocInfo struct {
	parent      *ssa.BasicBlock
	singleStore *ssa.StoreInstr
	stores      []*ssa.StoreInstr
}

// returns true if we can promote this alloca
func isAllocaPromotable(a *ssa.AllocInstr) (bool, *allocInfo) {
	var store *ssa.StoreInstr
	var stores []*ssa.StoreInstr
	singleStore := false
	for _, user := range a.Users() {
		switch i := user.(type) {
		case *ssa.StoreInstr:
			// a is used as an operand to store instruction
			if i.Operand(1) == a {
				return false, nil
			}
			if store != nil {
				singleStore = false
			} else {
				store = i
				singleStore = true
			}
			stores = append(stores, i)

		case *ssa.IndexInstr:
			if i.Operand(0) == a {
				return false, nil
			}

		case *ssa.MemberInstr:
			if i.Operand(0) == a {
				return false, nil
			}
		}
	}
	info := &allocInfo{}
	if singleStore {
		info.singleStore = store
	}

	return true, info
}

type phiEntry struct {
	phi   *ssa.PhiNode
	alloc *ssa.AllocInstr
}

type phiMap map[int][]phiEntry

type Mem2RegPass struct {
	dom     *analysis.DominatorAnalysisInfo
	current *ssa.Function
	phiMap  phiMap
}

func (m *Mem2RegPass) IsAnalysisPass() bool {
	return false
}

func (m *Mem2RegPass) replaceLoads(a *ssa.AllocInstr, v ssa.Value) bool {
	changed := false
	for _, u := range a.Users() {
		if l, ok := u.(*ssa.LoadInstr); ok {
			if ssa.Remove(l) != l {
				panic("unable to remove instr: " + l.String())
			}

			ssa.ReplaceInstr(l, v)
			changed = true
		}
	}

	return changed
}

func (m *Mem2RegPass) placePhi(bb *ssa.BasicBlock) *ssa.PhiNode {
	edges := []*ssa.PhiEdge{}
	for _, pred := range bb.Preds {
		edges = append(edges, &ssa.PhiEdge{pred, nil})
	}

	phi := ssa.NewPhiNode(edges, m.current)
	bb.PushFront(phi)
	return phi
}

func (m *Mem2RegPass) insertPhi(frontiers []int, a *ssa.AllocInstr) {
	for _, frontier := range frontiers {
		frontierBB := m.current.Blocks[frontier]
		phi := m.placePhi(frontierBB)
		m.phiMap[frontier] = append(m.phiMap[frontier], phiEntry{phi, a})
	}
}

func (m *Mem2RegPass) promote(a *ssa.AllocInstr, info *allocInfo) bool {
	if info.singleStore != nil {
		v := info.singleStore.Operand(1)
		ssa.Remove(a)
		ssa.Remove(info.singleStore)
		return m.replaceLoads(a, v)
	}

	for _, store := range info.stores {
		domFrontiers := m.dom.Frontiers[store.Parent().Index]
		m.insertPhi(domFrontiers, a)
	}

	return false
}

func (m *Mem2RegPass) Run(f *ssa.Function) bool {
	// this expects that dominator analysis has been already run
	dominatorAnalysis := ssa.GetPass("dominators")
	m.dom = dominatorAnalysis.(*analysis.DominatorAnalysis).GetInfo().(*analysis.DominatorAnalysisInfo)
	m.current = f
	changed := false
	var allocas map[*ssa.AllocInstr]*allocInfo

	allocas = make(map[*ssa.AllocInstr]*allocInfo)
	// collect all allocas
	for _, block := range f.Blocks {
		for _, instr := range block.Instructions() {
			if i, ok := instr.(*ssa.AllocInstr); ok {
				ok, info := isAllocaPromotable(i)
				if ok {
					info.parent = block
					allocas[i] = info
				}
			}
		}
	}

	if debug {
		fmt.Println("Promotable allocas:")
		for i, bb := range allocas {
			fmt.Printf("%s\t\t%s: store: %v\n", i, bb.parent.Name(), bb.singleStore)
		}
	}

	for alloc, bb := range allocas {
		c := m.promote(alloc, bb)
		if !changed {
			changed = c
		}
	}

	return changed
}

func init() {
	ssa.RegisterPass("mem2reg", &Mem2RegPass{})
}
