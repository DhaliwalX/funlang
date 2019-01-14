package res


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

func makeGraph(n int) Graph {
    return make([]Node, n)
}
