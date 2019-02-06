package passes

import (
	"fmt"
	"funlang/ds"
	"funlang/ssa"
)

type Verifier struct {
	errs ds.ErrorList
}

func (v *Verifier) IsAnalysisPass() bool { return true }

func (v *Verifier) collectError(f string, args ...interface{}) {
	err := fmt.Errorf(f, args...)
	v.errs = append(v.errs, err)
}

func (v *Verifier) verifyInstruction(i ssa.Instruction) {
	// each operand should have this instruction as user
	for _, o := range i.Operands() {
		inUserList := 0
		for _, u := range o.Users() {
			if u == i {
				inUserList++
			}
		}
		if inUserList == 0 {
			v.collectError("%s is not in user list of %s", i.String(), o.String())
		}
	}

	// each user should have this instruction as operand
	for _, u := range i.Users() {
		inOperandList := 0
		// XXX: PhiNode behaves differently as it stores operands as pairs
		if phi, ok := u.(*ssa.PhiNode); ok {
			for _, edge := range phi.Edges {
				if edge.Value == i {
					inOperandList++
				}
			}
		} else {
			for _, o := range u.(ssa.Instruction).Operands() {
				if o == i {
					inOperandList++
				}
			}
		}

		if inOperandList == 0 {
			v.collectError("%s is not in operand list of %s, but is present in formers user list",
				i.String(), u.String())
		}
	}
}

func (v *Verifier) verifyFunction(f *ssa.Function) {
	for _, bb := range f.Blocks {
		for i := bb.First; i != nil; i = i.Next() {
			v.verifyInstruction(i)
		}
	}
}

func (v *Verifier) Run(p *ssa.Program) bool {
	for _, global := range p.Globals {
		if f, ok := global.(*ssa.Function); ok {
			v.verifyFunction(f)
		}
	}

	if len(v.errs) > 0 {
		fmt.Println(v.errs.Error())
	}
	return false
}

func init() {
	ssa.RegisterPass("verifier", &Verifier{})
}
