package res


// disjoint set kind of data structure, with two methods
// Link(u, w) sets parent[w] = u
// Eval(x) is described as below,
//   return x if x is root
//   return u such that sdom(u) is minimum for all the nodes
//        in the path from r --> u --> x in tree

type dsu struct {
    ancestor []int
}

func makeDsu(n int) dsu {
    d := dsu{ancestor:make([]int, n)}
    for i := 0; i < n; i++ {
        d.link(i, i)
    }
    return d
}

func (d dsu) link(u, w int) {
    d.ancestor[w] = u
}

// this is the basic algorithm for eval,
// XXX: Optimize this algorithm using path compression
func (d dsu) eval(u int, semi []int) int {
    if d.ancestor[u] == u {
        return u
    }

    min := semi[u]
    minIdx := u
    for u != d.ancestor[u] {
        u = d.ancestor[u]
        if semi[u] < min {
            min = semi[u]
            minIdx = u
        }
    }

    return minIdx
}
