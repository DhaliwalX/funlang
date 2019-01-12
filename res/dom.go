// this file implements algorithm for building dominator tree for a graph
//  ref: A fast algorithm for finding dominators in a graph by Tarjan et. al.
package res

import (
    "bitbucket.org/dhaliwalprince/funlang/ssa"
    "fmt"
)

type domUtil struct {
    dom []int

    // holds semidominator information and before that will store dfs numbers
    sdom []int

    // holds the vertex number for a given dfs number
    // useful for reverse iterating dfs spanning tree
    vertex []int

    // holds parent of a vertex in spanning tree
    parent []int

    d dsu
    buckets [][]int
}

func (d *domUtil) eval(v int) int {
    return d.d.eval(v, d.sdom)
}

func (d *domUtil) link(u, w int) {
    d.d.link(u, w)
}

func (d *domUtil) String() string {
    return fmt.Sprintf(
        "SDOM: %v; DOM: %v; PARENT: %v; V: %v",
        d.sdom, d.dom, d.parent, d.vertex)
}

// main algorithm for computing dominators
func (d *domUtil) computeDominators(g graph) {
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
    d.dom = make([]int, total)
    d.buckets = make([][]int, total)


    for i := total-1; i > 0; i-- {
        v := d.vertex[i]

        // step 2: compute semi-dominators
        for _, pred := range g[v].pred {
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

        // step 3: dom
        for _, u := range d.buckets[d.parent[v]] {
            d.buckets[d.parent[v]] = d.buckets[d.parent[v]][1:]
            w := d.eval(u)
            if d.sdom[w] < d.sdom[u] {
                d.dom[u] = w
            } else {
                d.dom[u] = d.parent[v]
            }
        }
    }

    // step 4: fill dom implicitly
    for i := 1; i < total; i++ {
        w := d.vertex[i]
        if d.dom[w] != d.vertex[d.sdom[w]] {
            d.dom[w] = d.dom[d.dom[w]]
        }
    }

    // dom(root) == 0
    d.dom[0] = 0
}

func createGraph(blocks []*ssa.BasicBlock) graph {
    g := makeGraph(len(blocks))
    for _, block := range blocks {
        n := block.Index
        for _, succ := range block.Succs {
            g.addEdge(n, succ.Index)
        }
    }
    return g
}

func ComputeDominators(blocks []*ssa.BasicBlock) *domUtil {
    g := createGraph(blocks)
    d := &domUtil{}
    d.computeDominators(g)
    return d
}
