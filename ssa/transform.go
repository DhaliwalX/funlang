package ssa

func ReplaceInstr(i Instruction, v Value) {
	for _, u := range i.Users() {
		if t, ok := u.(Instruction); ok {
			for idx, o := range t.Operands() {
				if o == i {
					t.SetOperand(idx, v)
					v.AddUser(t)
				}
			}
		}
	}
}

func RemoveFromUsers(i Instruction) {
	for _, user := range i.Operands() {
		user.RemoveFromUsers(user)
	}
}

// Remove removes instruction from parent and returns next
func Remove(i Instruction) Instruction {
	next := i.Next()
	RemoveFromUsers(i)
	i.Parent().Remove(i)
	return next
}

// ReplaceOperand replace operand o of instruction i with value v
func ReplaceOperand(instr Instruction, o Value, v Value) {
	for i, operand := range instr.Operands() {
		if operand == o {
			instr.SetOperand(i, v)
			v.AddUser(instr)
		}
	}
}
