package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"bitbucket.org/dhaliwalprince/funlang/types"
	"fmt"
	"strconv"
)

// transform AST to SSA
type transformer struct {
	// program where the whole ssa will be stored
	program *Program

	// current function
	function *Function

	factory *types.Factory

	// name of temporaries
	counter int

	locals map[*ast.Identifier]Value

	address bool
}

func (t *transformer) constantInt(val int) *ConstantInt {
	return &ConstantInt{Value:val}
}

func (t *transformer) constantString(val string) *ConstantString {
	return &ConstantString{Value:val}
}

func isLVal(val Value) bool {
	return val.Type().Tag() == types.POINTER_TYPE
}

func isStruct(val Value) bool {
	return val.Type().Tag() == types.STRUCT_TYPE
}

func isArray(val Value) bool {
	return val.Type().Tag() == types.ARRAY_TYPE
}

func (t *transformer) nextTemp() string {
	return "t"+fmt.Sprint(t.counter)
}

// emit load instruction for named value
func (t *transformer) load(val Value) *LoadInstr {
	if !isLVal(val) {
		panic("load will only work for lvalues")
	}

	return &LoadInstr{users: []Value{},
		instrWithOperands: instrWithOperands{operands:[]Value{val}},
		valueWithName:valueWithName{name:t.nextTemp()}}
}

// *dst = src
func (t *transformer) store(dest Value, src Value) *StoreInstr {
	if !isLVal(dest) {
		panic("destination should be lvalue")
	}

	return &StoreInstr{
		instrWithOperands:instrWithOperands{operands:[]Value{dest, src}},
	}
}

func (trans *transformer) alloc(t types.Type) *AllocInstr {
	return &AllocInstr{t:t, valueWithName:valueWithName{name:trans.nextTemp()}}
}

func (t *transformer) member(val Value, member *ConstantString) *MemberInstr {
	if !isLVal(val) || (val.Type().Elem().Tag() != types.STRUCT_TYPE) {
		panic("destination should be a struct with address")
	}

	meTy := val.Type().Elem().Field(member.Value)

	return &MemberInstr{t: meTy,
		instrWithOperands: instrWithOperands{operands: []Value{val, member}},
		valueWithName:valueWithName{name:t.nextTemp()},
	}
}

func (t *transformer) index(val Value, member Value) *IndexInstr {
	if !isLVal(val) || (val.Type().Elem().Tag() != types.ARRAY_TYPE) {
		panic("destination should be a array with address")
	}

	return &IndexInstr{
		instrWithOperands:instrWithOperands{operands:[]Value{val, member}},
		valueWithName:valueWithName{name:t.nextTemp()},
	}
}

func (t *transformer) eval(op ArithOpcode, l, r int) int {
	switch op {
	case PLUS:
		return l+r
	case MINUS:
		return l-r
	case MUL:
		return l*r
	case DIV:
		return l/r
	case MOD:
		return l%r
	case XOR:
		return l^r
	case AND:
		return l&r
	case OR:
		return l|r
	case LT:
		if l < r { return 1 } else { return 0 }
	case GT:
		if l > r { return 1 } else { return 0 }
	case EQ:
		if l == r { return 1 } else { return 0 }
	}

	panic("unknown arithematic operation")
}

func (t *transformer) arith(op ArithOpcode, l, r Value) Value {
	if l.Tag() == CONSTANT_INT || r.Tag() ==  CONSTANT_STRING {
		return t.constantInt(t.eval(op, l.(*ConstantInt).Value, r.(*ConstantInt).Value))
	}

	return &ArithInstr{
		opCode:            op,
		instrWithOperands: instrWithOperands{operands: []Value{l, r}},
		valueWithName: valueWithName{name:t.nextTemp()},
	}
}

func (t *transformer) gotoif(condition Value, ontrue, onfalse *BasicBlock) *ConditionalGoto {
	return &ConditionalGoto{
		instrWithOperands:instrWithOperands{operands:[]Value{condition, ontrue, onfalse}},
	}
}

func (t *transformer) goTo(block *BasicBlock) *UnconditionalGoto {
	return &UnconditionalGoto{
		instrWithOperands:instrWithOperands{operands:[]Value{block}},
	}
}

func (t *transformer) call(f *Function, args ...Value) *CallInstr {
	operands := []Value{f}
	for _, arg := range args {
		operands = append(operands, arg)
	}
	return &CallInstr{
		instrWithOperands:instrWithOperands{operands:operands},
		valueWithName:valueWithName{name:t.nextTemp()},
	}
}

func (t *transformer) phi(edges []*PhiEdge) *PhiNode {
	return &PhiNode{Edges:edges, valueWithName:valueWithName{name:t.nextTemp()}}
}

func (t *transformer) astError(node ast.Node, message string) {
	s := fmt.Sprintf("%s: %s", node.Beg(), message)
	panic(s)
}

func (t *transformer) valueOf(i *ast.Identifier) Value {
	v, ok := t.locals[i]
	if !ok {
		v, ok = t.program.Globals[i.Name()]
		if !ok {
			t.astError(i, "undefined variable")
		}
	}

	return v
}

func (t *transformer) emit(val Value) Value {
	t.function.current.appendInstr(val)
	return val
}

func (t *transformer) emitLiteral(node ast.Node) Value {
	switch n := node.(type) {
	case *ast.NumericLiteral:
		i, err := strconv.Atoi(n.String())
		if err != nil {
			t.astError(n, "bad number representation")
		}
		return t.constantInt(i)

	case *ast.StringLiteral:
		return t.constantString(n.String())

	case *ast.Identifier:
		if t.address {
			return t.valueOf(n)
		} else {
			return t.load(t.valueOf(n))
		}

	case *ast.BooleanLiteral:
		if n.String() == "true" {
			return t.constantInt(1)
		} else {
			return t.constantInt(0)
		}
	}

	t.astError(node, "unknown ast node")
	return nil
}

func (t *transformer) emitArgs(x ast.Expression) []Value {
	old := t.address
	var args []Value
	t.address = false
	for _, arg := range x.(*ast.ArgumentList).Exprs() {
		a := t.emitExpression(arg)
		args = append(args, a)
	}

	t.address = old
	return args
}

func (t *transformer) emitMemberExpresion(m *ast.MemberExpression) Value {
	old := t.address

	var v Value
	switch m.AccessKind() {
	case lex.PERIOD:
		t.address = true
		val := t.emitExpression(m.Expr())
		member := m.Member().(*ast.Identifier).Name()
		v = t.member(val, t.constantString(member))

	case lex.LBRACK:
		t.address = true
		val := t.emitExpression(m.Expr())
		t.address = false
		index := t.emitExpression(m.Member())
		v = t.index(val, index)

	case lex.LPAREN:
		t.address = false
		val := t.emitExpression(m.Expr())
		if val.Tag() != FUNCTION {
			t.astError(m.Expr(), "not a function")
		}
		args := t.emitArgs(m.Member())
		v = t.call(val.(*Function), args...)
	}

	t.address = old
	return v
}

func (t *transformer) emitPrefixExpression(p *ast.PrefixExpression) Value {
	old := t.address
	var v Value
	switch (p.Op()) {
	case lex.AND:
		t.address = true
		v = t.emitExpression(p.Expression())

	case lex.MUL:
		t.address = false
		val := t.emitExpression(p.Expression())
		l := t.load(val)
		v = t.emit(l)
	}

	t.address = old
	return v
}

func (t *transformer) mapOp(op lex.TokenType) ArithOpcode {
	switch op {
	case lex.ADD:
		return PLUS

	case lex.SUB:
		return MINUS

	case lex.MUL:
		return MUL
	case lex.QUO:
		return DIV

	case lex.REM:
		return MOD

	case lex.AND:
		return AND
	case lex.OR:
		return OR

	case lex.XOR:
		return XOR
	case lex.LSS:
		return LT

	case lex.GTR:
		return GT

	case lex.EQL:
		return EQ
	}

	panic("unknown binary operator: "+op.String())
}

func (t *transformer) emitBinaryExpression(e *ast.BinaryExpression) Value {
	old := t.address
	t.address = false
	l := t.emitExpression(e.Left())
	r := t.emitExpression(e.Right())
	t.address = old
	return t.emit(t.arith(t.mapOp(e.Op()), l, r))
}

func (t *transformer) emitAssignExpression(e *ast.AssignExpression) Value {
	old := t.address
	t.address = false
	r := t.emitExpression(e.Left())
	t.address = true
	l := t.emitExpression(e.Right())
	t.address = old
	return t.emit(t.store(l, r))
}

func (t *transformer) emitExpression(e ast.Expression) Value {
	switch n := e.(type) {
	case *ast.NumericLiteral:
	case *ast.StringLiteral:
	case *ast.BooleanLiteral:
	case *ast.Identifier:
		return t.emitLiteral(e)
	case *ast.MemberExpression:
		return t.emitMemberExpresion(n)

	case *ast.PrefixExpression:
		return t.emitPrefixExpression(n)

	case *ast.BinaryExpression:
		return t.emitBinaryExpression(n)

	case *ast.AssignExpression:
		return t.emitAssignExpression(n)
	}

	panic(fmt.Sprintf("Unknown expr type: %T", e))
}

func (t *transformer) emitDeclaration(e ast.DeclNode) Value {
	switch n := e.(type) {
	case *ast.Declaration:

	}
}

func (t *transformer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	default:
		panic(n)
	}

	return nil
}
