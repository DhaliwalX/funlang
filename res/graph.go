package res


type node struct {
    ns []int
    pred []int
}

func (n *node) addNeighbour(i int) {
    n.ns = append(n.ns, i)
}

func (n *node) addPred(p int) {
    n.pred = append(n.pred, p)
}

type graph []node

func (g graph) addEdge(s, e int) {
    g[s].addNeighbour(e)
    g[e].addPred(s)
}

func makeGraph(n int) graph {
    return make([]node, n)
}
