package ssa

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
