package res

import (
	"fmt"
	"strings"
)

type domFrontierInfo struct {
	df [][]int
}

func (d *domFrontierInfo) init(g Graph) {
	d.df = make([][]int, len(g))
}

// Reference: https://www.ed.tus.ac.jp/j-mune/keio/m/ssa2.pdf
func (dfi *domFrontierInfo) computeDomFrontier(g Graph, util *DomInfo, n int) []int {
	df := []int{}
	for _, y := range g[n].Succs {
		if util.Dominator[y] != n {
			// Compute DF local
			df = append(df, y)
		}
	}

	for _, d := range g[n].Dominees {
		dfc := dfi.computeDomFrontier(g, util, d)
		for _, w := range dfc {
			if util.Dominator[w] != n {
				df = append(df, w)
			}
		}
	}
	return df
}

func (dfi *domFrontierInfo) computeDomFrontierIterative(g Graph, util *DomInfo) {
	dfi.init(g)
	for i, n := range g {
		if len(n.Preds) > 0 {
			for _, pred := range n.Preds {
				runner := pred
				for runner != util.Dominator[i] {
					dfi.df[runner] = append(dfi.df[runner], i)
					runner = util.Dominator[runner]
				}
			}
		}
	}
}

func (d *domFrontierInfo) Compute(g Graph, util *DomInfo) {
	d.init(g)
	for i := range g {
		df := d.computeDomFrontier(g, util, i)
		d.df[i] = df
	}
}

func (d *domFrontierInfo) String() string {
	builder := strings.Builder{}
	builder.WriteString("DomFrontierInfo:\n")
	for i, df := range d.df {
		builder.WriteString(fmt.Sprintf("%d: %v\n", i, df))
	}
	return builder.String()
}

func ComputeDominanceFrontiers(g Graph, util *DomInfo) [][]int {
	d := &domFrontierInfo{}
	d.computeDomFrontierIterative(g, util)
	return d.df
}
