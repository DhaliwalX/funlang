package passes

import "bitbucket.org/dhaliwalprince/funlang/ssa"

type FunctionPass interface {
	// returns true if something is changed, otherwise false
	Run(f *ssa.Function) bool
}

type BBPass interface {
	Run(b *ssa.BasicBlock) bool
}

type ProgramPass interface {
	Run(p *ssa.Program) bool
}
