package res

import "math"

type sdomUtil struct {
    sdom []int

    // reverse preorder of the nodes
    preOrdered []int
}

func swap(c []int, i, j int) {
    t := c[i]
    c[i] = c[j]
    c[j] = t
}

// partition the array to find the correct position of c[i]
func partition(c []int, b []int, s, e int) int {
    if e - s == 1 {
        if c[s] > c[e] {
        }
    }
    t := c[s]
    j := s
    k := e+1

    for j < k && (j <= e && k >= s) {
        for {
            j++
            if j > e || c[j] <= t {
                break
            }
        }

        for {
            k--

            if k < s || c[k] >= t {
                break
            }
        }

        if j < k {
            swap(c, j, k)
            swap(b, j, k)
        }
    }

    swap(c, k, s)
    swap(b, k, s)
    return k
}

func qsort(a []int, b []int, s, e int) {
    if s < e {
        p := partition(a, b, s, e)
        qsort(a, b, s, p-1)
        qsort(a, b, p+1, e)
    }
}

func sort(a []int, b []int) {
    c := make([]int, len(b))
    copy(c, b)
    qsort(a, c, 0, len(b)-1)
}

func (s *sdomUtil) init(g graph, d *dfsUtil) {
    s.preOrdered = make([]int, len(d.dfnum))
    s.sdom = make([]int, len(d.dfnum))
    for i := range g  {
        s.preOrdered[i] = i
        s.sdom[i] = -1
    }

    s.sdom[0] = 0

    sort(s.preOrdered, d.dfnum)
}

func (s *sdomUtil) evalNode(g graph, d *dfsUtil, r int) int {
    min := math.MaxInt32
    if r == 0 {
        return 0
    }

    if len(g[r].pred) == 1 {
        s.sdom[r] = g[r].pred[0]
        return s.sdom[r]
    }
    for _, p := range g[r].pred {
        if d.dfnum[p] > d.dfnum[r] {
            if s.sdom[p] == -1 {
                s.sdom[p] = s.evalNode(g, d, p)
            }

            if min > s.sdom[p] {
                min = s.sdom[p]
            }
            continue
        } else if p == 0 {
            // if root is predecessor, root is sdom
            min = 0
            break
        } else if min > p {
            min = p
        }
    }

    s.sdom[r] = min
    return min
}

func (s *sdomUtil) computeSdom(g graph, d *dfsUtil) {
    for i := range s.preOrdered {
        n := s.preOrdered[i]

        if s.sdom[n] == -1 {
            s.evalNode(g, d, n)
        }
    }
}

func sdom(g graph, d *dfsUtil) []int {
    s := sdomUtil{}
    s.init(g, d)
    s.computeSdom(g, d)
    return s.sdom
}
