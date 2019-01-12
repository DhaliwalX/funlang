package ssa

// tells whether this instruction is responsible
// for side effects
func IsSideEffect(i Instruction) bool {
	switch i.(type) {
	case *CallInstr, *ConditionalGoto, *UnconditionalGoto,
	*StoreInstr, *RetInstr:
		return true

	default:
		return false
	}
}
