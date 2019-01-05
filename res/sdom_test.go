package res

import "testing"

// [0 0 0 1 2]
func TestSDOM(t *testing.T) {
    g := sample()
    d := dfs(g)
    s := sdom(g, d)
    t.Log(s)
}

// [0 0 0 2 0 2]
func TestSDOM2(t *testing.T) {
    g := sample2()
    d := dfs(g)
    s := sdom(g, d)
    t.Log(s)
}

// [0 0 4 1 3 4]
func TestSDOM3(t *testing.T) {
    g := sample3()
    d := dfs(g)
    s := sdom(g, d)
    t.Log(s)
}
