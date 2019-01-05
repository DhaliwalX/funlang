package res

import (
    "testing"
)

func TestDOM(t *testing.T) {
    g := sample()
    d := domUtil{}
    d.computeDominators(g)
    t.Log(d.String())
}

func TestDOM2(t *testing.T) {
    g := sample2()
    d := domUtil{}
    d.computeDominators(g)
    t.Log(d.String())
}

func TestDOM3(t *testing.T) {
    g := sample3()
    d := domUtil{}
    d.computeDominators(g)
    t.Log(d.String())
}

func TestDOM4(t *testing.T) {
    g := sample4()
    d := domUtil{}
    d.computeDominators(g)
    t.Log(d.String())
}
