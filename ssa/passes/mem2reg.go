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

	defs []ssa.Instruction
	// our phis
	phis map[*ssa.PhiNode]bool
}

func (info *allocInfo) ourPhi(phi *ssa.PhiNode) bool {
	_, yes := info.phis[phi]
	return yes
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
	info := &allocInfo{stores: stores, phis: make(map[*ssa.PhiNode]bool)}
	if singleStore {
		info.singleStore = store
	}

	return true, info
}

type phiEntry struct {
	phi   *ssa.PhiNode
	alloc *ssa.AllocInstr
	store ssa.Instruction
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
			ssa.Remove(l)

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

	phi := ssa.NewPhiNode(edges, m.current, bb)
	bb.PushFront(phi)
	return phi
}

func (m *Mem2RegPass) insertPhi(frontiers []int, a *ssa.AllocInstr, s ssa.Instruction, info *allocInfo) {
	var seen, work BlockSet
	for _, frontier := range frontiers {
		work.add(frontier)
	}
	current := work.take()
	for ; current != -1; current = work.take() {
		frontierBB := m.current.Blocks[current]
		seen.add(current)
		phi := m.placePhi(frontierBB)
		m.phiMap[current] = append(m.phiMap[current], phiEntry{phi, a, s})
		info.phis[phi] = true
	}
}

func (m *Mem2RegPass) promote(a *ssa.AllocInstr, info *allocInfo) bool {
	if info.singleStore != nil {
		v := info.singleStore.Operand(1)
		ssa.Remove(a)
		ssa.Remove(info.singleStore)
		return m.replaceLoads(a, v)
	}

	// compute defs
	// alloc itself is a store instruction
	info.defs = append(info.defs, a)

	for _, store := range info.stores {
		info.defs = append(info.defs, store)
	}

	// place phis
	for _, def := range info.defs {
		domFrontiers := m.dom.Frontiers[def.Parent().Index]
		m.insertPhi(domFrontiers, a, def, info)
	}

	// rename phase
	m.rename(a, info)

	return false
}

func (m *Mem2RegPass) rename(a *ssa.AllocInstr, info *allocInfo) {
	var currentDef ssa.Value
	var phiNode *ssa.PhiNode
	defs := make([]ssa.Value, len(m.current.Blocks))
	for id, block := range m.current.Blocks {
		for instr := block.First; instr != nil; instr = instr.Next() {

		again:
			switch instr.(type) {
			case *ssa.AllocInstr:
				if a == instr {
					currentDef = ssa.NewConstant(instr.Type().Elem())
				}

			case *ssa.StoreInstr:
				if instr.Operand(0) == a {
					currentDef = instr.Operand(1)
					instr = ssa.Remove(instr)
					goto again
				}

			case *ssa.LoadInstr:
				if instr.Operand(0) == a {
					ssa.ReplaceInstr(instr, currentDef)
					instr = ssa.Remove(instr)
					goto again
				}

			case *ssa.PhiNode:
				phi := instr.(*ssa.PhiNode)
				if info.ourPhi(phi) {
					if phiNode != nil {
						// redundant as a block can have only one phi

						instr = ssa.Remove(instr)
						goto again
					}

					for _, edge := range phi.Edges {
						edge.Value = defs[edge.Block.Index]
					}
					currentDef = instr

					// to remove redundant phis
					phiNode = phi
				}
			}
		}
		// save the last def for this block
		defs[id] = currentDef
	}
}

// Run this pass
func (m *Mem2RegPass) Run(f *ssa.Function) bool {
	// this expects that dominator analysis has been already run
	m.phiMap = make(phiMap)
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
