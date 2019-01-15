package res

import (
    "fmt"
    "strings"
)

type Node struct {
    Succs    []int
    Preds    []int
    Dominees []int
}

func (n *Node) addNeighbour(i int) {
    n.Succs = append(n.Succs, i)
}

func (n *Node) addPred(p int) {
    n.Preds = append(n.Preds, p)
}

type Graph []Node

func (g Graph) addEdge(s, e int) {
    g[s].addNeighbour(e)
    g[e].addPred(s)
}

func (g Graph) Dot() string {
    builder := strings.Builder{}
    builder.WriteString(fmt.Sprintf("digraph cfg {\n"))
    for i := range g {
        for _, s := range g[i].Succs {
            builder.WriteString(fmt.Sprintf("\t%d -> %d\n", i, s))
        }
    }
    builder.WriteString("}\n")
    return builder.String()
}

func makeGraph(n int) Graph {
    return make([]Node, n)
}
