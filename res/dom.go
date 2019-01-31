// this file implements algorithm for building dominator tree for a Graph
//  ref: A fast algorithm for finding dominators in a Graph by Tarjan et. al.
package res

import (
	"fmt"
	"funlang/ssa"
)

type DomInfo struct {
	Dominator []int

	// holds semidominator information and before that will store dfs numbers
	sdom []int

	// holds the vertex number for a given dfs number
	// useful for reverse iterating dfs spanning tree
	vertex []int

	// holds parent of a vertex in spanning tree
	parent []int

	d       dsu
	buckets [][]int
}

func (d *DomInfo) eval(v int) int {
	return d.d.eval(v, d.sdom)
}

func (d *DomInfo) link(u, w int) {
	d.d.link(u, w)
}

func (d *DomInfo) String() string {
	return fmt.Sprintf(
		"SDOM: %v; DOM: %v; PARENT: %v; V: %v",
		d.sdom, d.Dominator, d.parent, d.vertex)
}

// this will fill Graph nodes with Dominees information
func (d *DomInfo) fillInfo(g Graph) {
	for i, d := range d.Dominator {
		if i == 0 && d == 0 {
			continue
		}
		g[d].Dominees = append(g[d].Dominees, i)
	}
}

// main algorithm for computing dominators
func (d *DomInfo) computeDominators(g Graph) {
	// step 1: carry out dfs and initialize structure properly
	df := dfs(g)
	d.sdom = df.dfnum
	d.parent = df.parent
	total := len(g)
	d.vertex = make([]int, total)
	for i, v := range d.sdom {
		d.vertex[v] = i
	}
	d.d = makeDsu(total)
	d.Dominator = make([]int, total)
	d.buckets = make([][]int, total)

	for i := total - 1; i > 0; i-- {
		v := d.vertex[i]

		// step 2: compute semi-dominators
		for _, pred := range g[v].Preds {
			u := d.eval(pred)
			if d.sdom[u] < d.sdom[v] {
				d.sdom[v] = d.sdom[u]
			}
		}

		// add w to bucket(vertex(semi(w))
		vs := d.vertex[d.sdom[v]]
		d.buckets[vs] = append(d.buckets[vs], v)

		// LINK(parent(w), w)
		d.link(d.parent[v], v)

		// step 3: Dominator
		for _, u := range d.buckets[d.parent[v]] {
			d.buckets[d.parent[v]] = d.buckets[d.parent[v]][1:]
			w := d.eval(u)
			if d.sdom[w] < d.sdom[u] {
				d.Dominator[u] = w
			} else {
				d.Dominator[u] = d.parent[v]
			}
		}
	}

	// step 4: fill Dominator implicitly
	for i := 1; i < total; i++ {
		w := d.vertex[i]
		if d.Dominator[w] != d.vertex[d.sdom[w]] {
			d.Dominator[w] = d.Dominator[d.Dominator[w]]
		}
	}

	// Dominator(root) == 0
	d.Dominator[0] = 0

	d.fillInfo(g)
}

func CreateGraph(blocks []*ssa.BasicBlock) Graph {
	g := makeGraph(len(blocks))
	for _, block := range blocks {
		n := block.Index
		for _, succ := range block.Succs {
			g.addEdge(n, succ.Index)
		}
	}
	return g
}

func ComputeDominators(g Graph) *DomInfo {
	d := &DomInfo{}
	d.computeDominators(g)
	return d
}
