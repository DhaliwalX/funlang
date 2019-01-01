package res

import "testing"

func TestSDOM(t *testing.T) {
    g := sample()
    d := dfs(g)
    s := sdom(g, d)
    t.Log(s)
}

func TestSDOM2(t *testing.T) {
    g := sample2()
    d := dfs(g)
    s := sdom(g, d)
    t.Log(s)
}

func TestSDOM3(t *testing.T) {
    g := sample3()
    d := dfs(g)
    s := sdom(g, d)
    t.Log(s)
}
