package res

import "testing"

func sample() Graph {
    g := makeGraph(5)
    g.addEdge(0, 1)
    g.addEdge(0, 2)
    g.addEdge(1, 2)
    g.addEdge(1, 3)
    g.addEdge(2, 3)
    g.addEdge(2, 4)
    g.addEdge(3, 4)
    return g
}

func sample2() Graph {
    g := makeGraph(6)
    g.addEdge(0, 1)
    g.addEdge(0, 2)
    g.addEdge(0, 4)
    g.addEdge(2, 3)
    g.addEdge(2, 4)
    g.addEdge(2, 5)
    g.addEdge(4, 5)
    g.addEdge(3, 1)
    return g
}


func sample3() Graph {
    g := makeGraph(6)
    g.addEdge(0, 1)
    g.addEdge(1, 3)
    g.addEdge(2, 1)
    g.addEdge(3, 4)
    g.addEdge(4, 2)
    g.addEdge(4, 5)
    return g
}


func sample4() Graph {
    g := makeGraph(6)
    g.addEdge(0, 1)
    g.addEdge(1, 2)
    g.addEdge(1, 3)
    g.addEdge(2, 4)
    g.addEdge(4, 2)
    g.addEdge(3, 5)
    g.addEdge(4, 5)
    return g
}

func sample5() Graph {
    g := makeGraph(6)
    g.addEdge(0, 1)
    g.addEdge(0, 2)
    g.addEdge(1, 3)
    g.addEdge(3, 1)
    g.addEdge(2, 3)
    g.addEdge(2, 4)
    g.addEdge(3, 5)
    g.addEdge(4, 5)
    return g
}

func TestDfs(t *testing.T) {
    g := sample3()
    d := dfs(g)
    t.Log(d.dfnum)
}
