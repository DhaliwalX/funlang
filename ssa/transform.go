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

func Remove(i Instruction) Instruction {
	RemoveFromUsers(i)
	i.Parent().Remove(i)
	return i
}
