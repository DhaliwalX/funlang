package ssa

import (
	"container/list"
	"fmt"
)

type Pass interface {
	IsAnalysisPass() bool
}

type FunctionPass interface {
	Pass
	// returns true if something is changed, otherwise false
	Run(f *Function) bool
}

type BBPass interface {
	Pass
	Run(b *BasicBlock) bool
}

type ProgramPass interface {
	Pass
	Run(p *Program) bool
}

type PassRunner struct {
	runList *list.List
	p       *Program
}

func NewPassRunner(program *Program) *PassRunner {
	p := &PassRunner{runList: list.New(), p: program}
	return p
}

func (runner *PassRunner) Next() Pass {
	if runner.runList.Len() == 0 {
		return nil
	}
	p := runner.runList.Remove(runner.runList.Front())
	return p.(Pass)
}

func (runner *PassRunner) AddNext(p Pass) {
	runner.runList.PushFront(p)
}

func (runner *PassRunner) Add(p Pass) {
	runner.runList.PushBack(p)
}

func (runner *PassRunner) runFunctionPass(p FunctionPass) {
	for _, global := range runner.p.Globals {
		if f, ok := global.(*Function); ok {
			changed := p.Run(f)
			if changed {
				fmt.Printf("%T: changed %s\n", p, f.Name())
			}
		}
	}
}

type helperBBPass struct{ b BBPass }

func (h *helperBBPass) IsAnalysisPass() bool { return false }
func (h *helperBBPass) Run(f *Function) bool {
	changed := false
	for _, bb := range f.Blocks {
		changed = h.b.Run(bb)
	}
	return changed
}

func (runner *PassRunner) runBBPass(p BBPass) {
	runner.AddNext(&helperBBPass{b: p})
}

func (runner *PassRunner) runPass(p Pass) {
	fmt.Printf("::== Running pass: %T\n", p)
	switch pass := p.(type) {
	case FunctionPass:
		runner.runFunctionPass(pass)
		// fmt.Print(runner.p)

	case BBPass:
		runner.runBBPass(pass)

	case ProgramPass:
		pass.Run(runner.p)
	default:
		panic("don't know how run passes of kind: " + fmt.Sprintf("%T", pass))
	}
}

func (runner *PassRunner) RunPasses() {
	for p := runner.Next(); p != nil; p = runner.Next() {
		runner.runPass(p)
	}
}
