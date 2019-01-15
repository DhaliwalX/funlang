package res

import (
    "testing"
)

func TestDOM(t *testing.T) {
    g := sample()
    d := DomInfo{}
    d.computeDominators(g)
    t.Log(d.String())
}

func TestDOM2(t *testing.T) {
    g := sample2()
    d := DomInfo{}
    d.computeDominators(g)
    t.Log(d.String())
}

func TestDOM3(t *testing.T) {
    g := sample3()
    d := DomInfo{}
    d.computeDominators(g)
    t.Log(d.String())
}

func TestDOM4(t *testing.T) {
    g := sample4()
    d := DomInfo{}
    d.computeDominators(g)
    t.Log(d.String())
}

func TestDOM5(t *testing.T) {
    g := sample5()
    d := DomInfo{}
    d.computeDominators(g)
    t.Log(d.String())
}

func TestDOM6(t *testing.T) {
    g := sample6()
    t.Log(g.Dot())
    d := DomInfo{}
    d.computeDominators(g)
    t.Log(d.String())
}
