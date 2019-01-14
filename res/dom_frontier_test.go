package res

import "testing"

func TestDomFrontierInfo_Compute(t *testing.T) {
	g := sample5()
	d := DomInfo{}
	d.computeDominators(g)
	t.Log(d.String())
	df := domFrontierInfo{}
	df.Compute(g, &d)
	t.Log(df.String())
}

func TestDomFrontierIterative(t *testing.T) {
	g := sample5()
	d := DomInfo{}
	d.computeDominators(g)
	t.Log(d.String())
	df := domFrontierInfo{}
	df.computeDomFrontierIterative(g, &d)
	t.Log(df.String())
}
