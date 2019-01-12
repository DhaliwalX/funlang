package codegen

import (
	"bitbucket.org/dhaliwalprince/funlang/ssa"
	"bitbucket.org/dhaliwalprince/funlang/types"
	"fmt"
	"strings"
)

type gofunction struct {
	blocks map[*ssa.BasicBlock]string
	decls map[string]string
}

func (f *gofunction) addDecl(name, t string) {
	if dt, _ := f.decls[name]; dt != t {
		panic("multiple decls for "+name+" t: "+t)
	}
	f.decls[name] = t
}

func newGoFunction() *gofunction {
	return &gofunction{blocks: make(map[*ssa.BasicBlock]string)}
}

type goBackendState struct {
	currentFunction *gofunction
}

type GoBackend struct {
	builder strings.Builder
	state goBackendState
}

func (backend *GoBackend) String() string {
	return backend.builder.String()
}

func (backend *GoBackend) write(f string, args ...interface{}) {
	backend.builder.WriteString(fmt.Sprintf(f, args...))
}

func (backend *GoBackend) writeln(f string, args ...interface{}) {
	backend.write(f+"\n", args...)
}

func (backend *GoBackend) addDecl(name, t string) {
	backend.writeln("%s %s;", t, name)
}

func (backend *GoBackend) writeDecl(a *ssa.AllocInstr) {
	backend.writeln("%s = (%s)malloc(sizeof(%s));", a.Name(), a.Type(), a.Type().Elem())
}

func (backend *GoBackend) writeArith(a *ssa.ArithInstr) {
	backend.writeln("%s = %s %s %s;", a.Name(), a.Operand(0).Name(), a.Op().String(), a.Operand(1).Name())
}

func (backend *GoBackend) writeCall(c *ssa.CallInstr) {
	retType := types.ToFunctionType(c.Operand(0).Type()).ReturnType()
	if retType == nil {
		backend.write("%s(", c.Operand(0).Name())
	} else {
		backend.write("%s = %s(", c.Name(), c.Operand(0).Name())
	}

	for i, arg := range c.Operands() {
		if i == 0 {
			continue
		}

		backend.write(arg.Name())
		if i != len(c.Operands()) - 1 {
			backend.write(",")
		}
	}

	backend.writeln(");")
}

func (backend *GoBackend) generateInstruction(i ssa.Instruction) {
	switch instr := i.(type) {
	case *ssa.AllocInstr:
		backend.writeDecl(instr)

	case *ssa.LoadInstr:
		backend.writeln("%s = *%s;", instr.Name(), instr.Operand(0).Name())

	case *ssa.StoreInstr:
		backend.writeln("*%s = %s;", instr.Operand(0).Name(), instr.Operand(1).Name())

	case *ssa.ArithInstr:
		backend.writeArith(instr)

	case *ssa.UnconditionalGoto:
		backend.writeln("goto %s;", removeDot(instr.Operand(0).Name()))

	case *ssa.ConditionalGoto:
		backend.writeln("if (%s) goto %s; else goto %s;",
			instr.Operand(0).Name(), removeDot(instr.Operand(1).Name()),
			removeDot(instr.Operand(2).Name()))

	case *ssa.RetInstr:
		backend.writeln("return %s;", instr.Operand(0).Name())

	case *ssa.MemberInstr:
		backend.writeln("%s = &%s.%s;", instr.Name(), instr.Operand(0).Name(), instr.Operand(1).Name())

	case *ssa.IndexInstr:
		backend.writeln("%s = &%s[%s];", instr.Name(), instr.Operand(0).Name(), instr.Operand(1).Name())

	case *ssa.CallInstr:
		backend.writeCall(instr)
	default:
		panic("unhandled instruction")
	}
}

func removeDot(s string) string {
	return strings.Join(strings.Split(s, "."), "_")
}

func shouldSkipInstr(i ssa.Instruction) bool {
	switch i.(type) {
	case *ssa.StoreInstr, *ssa.RetInstr, *ssa.ConditionalGoto, *ssa.UnconditionalGoto:
		return true

	case *ssa.CallInstr:
		if i.Type() == nil {
			return true
		}
	}
	return false
}

func (backend *GoBackend) collectDecls(f *ssa.Function) {
	for _, bb := range f.Blocks {
		for _, i := range bb.Instructions() {
			if shouldSkipInstr(i) {
				continue
			}

			backend.addDecl(i.Name(), i.Type().Name())
		}
	}
}

func (backend *GoBackend) generateBlock(bb *ssa.BasicBlock) {
	if len(bb.Users()) > 0 {
		backend.writeln("%s:", removeDot(bb.Name()))
	}
	for _, instr := range bb.Instructions() {
		backend.generateInstruction(instr)
	}
}

func (backend *GoBackend) generateFunction(f *ssa.Function) {
	if f.Extern {
		// we don't need to codegen extern function
		return
	}
	if retType := types.ToFunctionType(f.Type()).ReturnType(); retType != nil {
		backend.write(retType.Name())
	} else {
		backend.write("void")
	}
	backend.write(" %s(", f.Name())
	count := 0
	for _, arg := range f.Args {
		backend.write("%s %s", arg.Type().Name(), arg.Name())
		count++
		if count != len(f.Args) {
			backend.write(",")
		}
	}

	backend.write(")")

	backend.write("{\n")
	backend.collectDecls(f)
	backend.state.currentFunction = newGoFunction()
	for _, bb := range f.Blocks {
		backend.generateBlock(bb)
	}

	backend.writeln("}\n")
}

func (backend *GoBackend) Run(program *ssa.Program) bool {
	backend.writeln("#include <stdio.h>")
	for _, global := range program.Globals {
		switch v := global.(type) {
		case *ssa.Function:
			backend.generateFunction(v)

		default:
			panic(fmt.Sprintf("%T unknown type", v))
		}
	}

	return false
}
