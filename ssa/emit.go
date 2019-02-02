package ssa

import (
	"fmt"
	"funlang/ds"
	"strconv"

	"funlang/ast"
	"funlang/context"
	"funlang/lex"
	"funlang/types"
)

func NewPhiNode(edges []*PhiEdge, f *Function, bb *BasicBlock) *PhiNode {
	phi := &PhiNode{valueWithName: valueWithName{name: f.NextName()}, Edges: edges}
	phi.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: phi}
	for _, edge := range edges {
		edge.Block.AddUser(phi)

		// mem2reg pass will often generate nil value entries while placing
		// phis.
		if edge.Value != nil {
			edge.Value.AddUser(phi)
		}
	}

	phi.parent = bb

	return phi
}

// transform AST to SSA
type transformer struct {
	// program where the whole ssa will be stored
	program *Program

	// current function
	function *Function

	factory *types.Factory

	// name of temporaries
	counter int

	types map[string]types.Type

	address bool
}

func (t *transformer) constantInt(val int) *ConstantInt {
	return &ConstantInt{Value: val}
}

func (t *transformer) constantString(val string) *ConstantString {
	return &ConstantString{Value: val}
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
	return t.function.NextName()
}

func (t *transformer) operands(user Value, ops ...Value) instrWithOperands {
	for _, op := range ops {
		op.AddUser(user)
	}

	return instrWithOperands{operands: ops}
}

// emit load instruction for named value
func (t *transformer) load(val Value) *LoadInstr {
	if !isLVal(val) {
		panic("load will only work for lvalues")
	}

	instr := &LoadInstr{users: []Value{},
		valueWithName: valueWithName{name: t.nextTemp()}}
	instr.instrWithOperands = t.operands(instr, val)
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	instr.parent = t.function.current
	return instr
}

// *dst = src
func (t *transformer) store(dest Value, src Value) *StoreInstr {
	if !isLVal(dest) {
		panic("destination should be lvalue")
	}

	instr := &StoreInstr{}
	instr.instrWithOperands = t.operands(instr, dest, src)
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	instr.parent = t.function.current
	return instr
}

func (trans *transformer) alloc(t types.Type) *AllocInstr {
	instr := &AllocInstr{t: t, valueWithName: valueWithName{name: trans.nextTemp()}}
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	instr.parent = trans.function.current
	return instr
}

func (t *transformer) member(val Value, member *ConstantString) *MemberInstr {
	if !isLVal(val) || (val.Type().Elem().Tag() != types.STRUCT_TYPE) {
		panic("destination should be a struct with address")
	}

	meTy := t.factory.PointerType(val.Type().Elem().Field(member.Value))

	instr := &MemberInstr{t: meTy,
		instrWithOperands: instrWithOperands{operands: []Value{val, member}},
		valueWithName:     valueWithName{name: t.nextTemp()},
	}

	instr.parent = t.function.current
	instr.instrWithOperands = t.operands(instr, val, member)
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	return instr
}

func (t *transformer) index(val Value, member Value) *IndexInstr {
	if !isLVal(val) || (val.Type().Elem().Tag() != types.ARRAY_TYPE) {
		panic("destination should be an array with address")
	}

	instr := &IndexInstr{
		t:                 t.factory.PointerType(val.Type().Elem()),
		instrWithOperands: instrWithOperands{operands: []Value{val, member}},
		valueWithName:     valueWithName{name: t.nextTemp()},
	}

	instr.parent = t.function.current
	instr.instrWithOperands = t.operands(instr, val, member)
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	return instr
}

func (t *transformer) eval(op ArithOpcode, l, r int) int {
	switch op {
	case PLUS:
		return l + r
	case MINUS:
		return l - r
	case MUL:
		return l * r
	case DIV:
		return l / r
	case MOD:
		return l % r
	case XOR:
		return l ^ r
	case AND:
		return l & r
	case OR:
		return l | r
	case LT:
		if l < r {
			return 1
		} else {
			return 0
		}
	case GT:
		if l > r {
			return 1
		} else {
			return 0
		}
	case EQ:
		if l == r {
			return 1
		} else {
			return 0
		}
	}

	panic("unknown arithematic operation")
}

func (t *transformer) arith(op ArithOpcode, l, r Value) Value {
	if l.Tag() == CONSTANT_INT || r.Tag() == CONSTANT_STRING {
		return t.constantInt(t.eval(op, l.(*ConstantInt).Value, r.(*ConstantInt).Value))
	}

	instr := &ArithInstr{
		opCode:            op,
		instrWithOperands: instrWithOperands{operands: []Value{l, r}},
		valueWithName:     valueWithName{name: t.nextTemp()},
	}
	instr.parent = t.function.current
	instr.instrWithOperands = t.operands(instr, l, r)
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	return instr
}

func (t *transformer) gotoif(condition Value, ontrue, onfalse *BasicBlock) *ConditionalGoto {
	instr := &ConditionalGoto{
		instrWithOperands: instrWithOperands{operands: []Value{condition, ontrue, onfalse}},
	}

	instr.parent = t.function.current
	condition.AddUser(instr)
	ontrue.AddUser(instr)
	onfalse.AddUser(instr)
	t.function.current.AddSucc(ontrue)
	t.function.current.AddSucc(onfalse)

	ontrue.AddPred(t.function.current)
	onfalse.AddPred(t.function.current)
	// not tracking users for basicblocks
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	return instr
}

func (t *transformer) goTo(block *BasicBlock) *UnconditionalGoto {
	instr := &UnconditionalGoto{
		instrWithOperands: instrWithOperands{operands: []Value{block}},
	}
	instr.parent = t.function.current
	block.AddUser(instr)
	t.function.current.AddSucc(block)
	block.AddPred(instr.parent)
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	return instr
}

func (t *transformer) call(f *Function, args ...Value) *CallInstr {
	operands := []Value{f}
	for _, arg := range args {
		operands = append(operands, arg)
	}
	instr := &CallInstr{
		instrWithOperands: instrWithOperands{operands: operands},
		valueWithName:     valueWithName{name: t.nextTemp()},
	}

	instr.parent = t.function.current
	for _, op := range instr.operands {
		op.AddUser(instr)
	}

	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	return instr
}

func (t *transformer) ret(val Value) *RetInstr {
	retType := types.ToFunctionType(t.function.Type()).ReturnType()
	if val == nil {
		if retType != nil {
			panic("return cannot be nil for functions with non-nil return type ")
		} else {
			instr := &RetInstr{}
			instr.parent = t.function.current
			instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
			return instr
		}
	}

	if val.Type() != retType {
		panic("return value type does not match with functions return type: " + t.function.name)
	} else {
		instr := &RetInstr{
			instrWithOperands: instrWithOperands{operands: []Value{val}},
		}

		instr.parent = t.function.current
		val.AddUser(instr)
		instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
		return instr
	}
}

func (t *transformer) phi(edges []*PhiEdge) *PhiNode {
	instr := &PhiNode{Edges: edges, valueWithName: valueWithName{name: t.nextTemp()}}
	for _, edge := range edges {
		edge.Value.AddUser(instr)
	}
	instr.parent = t.function.current
	instr.listElement = &ds.ListElement{Next: nil, Prev: nil, Value: instr}
	return instr
}

func (t *transformer) astError(node ast.Node, message string) {
	s := fmt.Sprintf("%s: %s", node.Beg(), message)
	panic(s)
}

func (t *transformer) valueOf(i *ast.Identifier) Value {
	v, ok := t.function.locals[i.Name()]
	if !ok {
		if i.Object == nil {
			panic("unreachable code")
		}

		if i.Object.Decl == nil {
			t.astError(i, i.Name()+" without declaration")
		}

		v, ok = t.program.Globals[i.Name()]
		if !ok {
			v, ok = t.function.Args[i.Name()]
			// create a local copy of the argument and store it in locals
			if !ok {
				t.astError(i, "undefined variable")
			}

			emit := t.emit(t.alloc(v.Type()))
			t.function.locals[i.Name()] = emit
			t.emit(t.store(emit, v))
			return emit
		}
	}

	return v
}

func (t *transformer) emit(val Instruction) Value {
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
			return t.emit(t.load(t.valueOf(n)))
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
		t.address = true
		val := t.emitExpression(m.Expr())
		if val.Tag() != FUNCTION {
			t.astError(m.Expr(), "not a function")
		}
		t.address = false
		args := t.emitArgs(m.Member())
		v = t.call(val.(*Function), args...)
	}

	t.address = old

	if t.address {
		return t.emit(v.(Instruction))
	} else {
		if v.Type() != nil && v.Type().Tag() != types.POINTER_TYPE {
			return t.emit(v.(Instruction))
		}
		if _, ok := v.(Instruction); ok {
			t.emit(v.(Instruction))
		}

		// v.Type() is nil only when function doesn't have return type
		if v.Type() == nil {
			return nil
		}
		return t.emit(t.load(v))
	}
}

func (t *transformer) emitPrefixExpression(p *ast.PrefixExpression) Value {
	old := t.address
	var v Value
	switch p.Op() {
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

	panic("unknown binary operator: " + op.String())
}

func (t *transformer) emitLogicalExpression(e *ast.BinaryExpression) Value {
	old := t.address
	t.address = false
	l := t.emitExpression(e.Left())

	x := t.function.current
	next := &BasicBlock{Parent: t.function,
		valueWithName: valueWithName{name: t.nextTemp()}, Index: len(t.function.Blocks)}

	final := &BasicBlock{Parent: t.function,
		valueWithName: valueWithName{name: t.nextTemp()}, Index: len(t.function.Blocks) + 1}
	t.function.Blocks = append(t.function.Blocks, next)
	t.function.Blocks = append(t.function.Blocks, final)
	if e.Op() == lex.LAND {
		t.emit(t.gotoif(l, next, final))
	} else if e.Op() == lex.LOR {
		t.emit(t.gotoif(l, final, next))
	} else {
		panic("illegal logical operator")
	}
	t.function.current = next
	r := t.emitExpression(e.Right())
	t.function.current = final

	// emit phi instruction
	edges := []*PhiEdge{&PhiEdge{Block: x, Value: l}, &PhiEdge{Block: next, Value: r}}
	v := t.emit(t.phi(edges))
	t.address = old
	return v
}

func (t *transformer) emitBinaryExpression(e *ast.BinaryExpression) Value {
	old := t.address
	t.address = false
	l := t.emitExpression(e.Left())
	r := t.emitExpression(e.Right())
	t.address = old

	if e.Op() == lex.LAND || e.Op() == lex.LOR {
		return t.emitLogicalExpression(e)
	}
	v := t.arith(t.mapOp(e.Op()), l, r)
	if v.Tag() != INSTRUCTION {
		return v
	}
	return t.emit(v.(Instruction))
}

func (t *transformer) emitAssignExpression(e *ast.AssignExpression) Value {
	old := t.address
	t.address = false
	r := t.emitExpression(e.Right())
	t.address = true
	l := t.emitExpression(e.Left())
	t.address = old
	return t.emit(t.store(l, r))
}

func (t *transformer) emitExpression(e ast.Expression) Value {
	switch n := e.(type) {
	case *ast.NumericLiteral:
		return t.emitLiteral(n)
	case *ast.StringLiteral:
		return t.emitLiteral(n)
	case *ast.BooleanLiteral:
		return t.emitLiteral(n)
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

func (t *transformer) resolveType(typeExpr ast.Expression) types.Type {
	switch n := typeExpr.(type) {
	case *ast.ArrayType:
		return t.factory.ArrayType(t.resolveType(n.Type()))

	case *ast.StructType:
		fields := make(map[string]types.Type)
		for _, field := range n.Fields() {
			fieldTy := t.resolveType(field.Type())
			fields[field.Name()] = fieldTy
		}

		return t.factory.StructType(fields)

	case *ast.FuncType:
		retType := t.resolveType(n.Return())
		argTypes := []types.Type{}
		for _, argType := range n.Params() {
			argTy := t.resolveType(argType)
			argTypes = append(argTypes, argTy)
		}

		return t.factory.FunctionType(retType, argTypes)

	case *ast.Identifier:
		switch n.Name() {
		case "string":
			return t.factory.StringType()

		case "int":
			return t.factory.IntType()

		case "bool":
			return t.factory.IntType()

		default:
			if ty, ok := t.types[n.Name()]; ok {
				return ty
			} else if n.Object == nil || n.Object.Type == nil {
				panic("unable to resolve type for " + n.Name() + " at " + n.Beg().String())
			} else {
				// try to resolve this type
				tr := t.resolveType(n.Object.Type.(ast.Expression))
				t.types[n.Name()] = tr
				return tr
			}
		}
	}

	panic(fmt.Sprint("unknown type found: %T", typeExpr))
}

func (t *transformer) emitDeclaration(e ast.DeclNode) Value {
	switch n := e.(type) {
	case *ast.Declaration:
		if n.Type() == nil {
			old := t.address
			t.address = false
			init := t.emitExpression(n.Init())
			t.address = old
			alloc := t.emit(t.alloc(init.Type()))
			t.function.locals[n.Name()] = alloc
			t.emit(t.store(alloc, init))
		} else {
			old := t.address
			t.address = false
			ty := t.resolveType(n.Type())
			alloc := t.emit(t.alloc(ty))
			if n.Init() != nil {
				init := t.emitExpression(n.Init())
				t.emit(t.store(alloc, init))
			}
			t.function.locals[n.Name()] = alloc
			t.address = old
		}
		return nil

	case *ast.TypeDeclaration:
		ty := t.resolveType(n.Type())
		t.program.Types[n.Name()] = ty
		return nil
	}

	panic("unknown declaration" + fmt.Sprintf("%T", e))
}

func (t *transformer) resolveFunctionSignature(f *ast.FunctionProtoType) types.Type {
	argTypes := []types.Type{}
	for _, arg := range f.Params() {
		argTypes = append(argTypes, t.resolveType(arg.(*ast.Declaration).Type()))
	}

	var retType types.Type
	if f.Return() != nil {
		retType = t.resolveType(f.Return())
	}
	return t.factory.FunctionType(retType, argTypes)
}

func (t *transformer) emitFunction(f *ast.FunctionStatement) {
	f.Proto()
	fun := &Function{}
	fun.name = f.Proto().Name()
	fun.t = t.resolveFunctionSignature(f.Proto())
	fun.locals = make(map[string]Value)
	args := make(map[string]*Argument)

	entryBlock := &BasicBlock{Parent: fun, valueWithName: valueWithName{name: "entry." + fun.name}}
	fun.Blocks = []*BasicBlock{entryBlock}
	fun.current = entryBlock
	entryBlock.Index = 0
	t.program.Globals[fun.name] = fun
	t.function = fun
	t.counter = 0
	for _, ar := range f.Proto().Params() {
		decl := ar.(*ast.Declaration)
		argType := t.resolveType(decl.Type())
		arg := &Argument{valueWithName: valueWithName{name: decl.Name()}, t: argType}
		args[decl.Name()] = arg

		// we need to copy the arguments to local variables, later optimisations can remove
		// those local variables and replace them with original arguments
		emit := t.emit(t.alloc(argType))
		t.function.locals[decl.Name()] = emit
		t.emit(t.store(emit, arg))
	}

	fun.Args = args
	// emit function body
	if f.Body() == nil {
		fun.Extern = true
		return
	}
	t.Visit(f.Body())
}

func (t *transformer) emitReturn(x ast.Expression) {
	old := t.address
	t.address = false
	var val Value
	if x != nil {
		val = t.emitExpression(x)
	}
	t.address = old
	t.emit(t.ret(val))
}

func (t *transformer) emitIfElseStatement(e *ast.IfElseStatement) {
	cond := t.emitExpression(e.Condition())
	label := t.nextTemp()
	onTrue := &BasicBlock{Parent: t.function,
		valueWithName: valueWithName{name: "if.true." + label}, Index: len(t.function.Blocks)}
	t.function.Blocks = append(t.function.Blocks, onTrue)
	onFalse := &BasicBlock{Parent: t.function,
		valueWithName: valueWithName{name: "if.false." + label}, Index: len(t.function.Blocks)}
	t.function.Blocks = append(t.function.Blocks, onFalse)
	done := &BasicBlock{Parent: t.function,
		valueWithName: valueWithName{name: "if.done." + label}, Index: len(t.function.Blocks)}
	t.function.Blocks = append(t.function.Blocks, done)

	t.emit(t.gotoif(cond, onTrue, onFalse))
	t.function.current = onTrue
	t.Visit(e.Body())
	t.emit(t.goTo(done))

	t.function.current = onFalse
	if e.ElseNode() != nil {
		t.Visit(e.ElseNode())
	}
	t.emit(t.goTo(done))

	t.function.current = done
}

func (t *transformer) emitForStatement(f *ast.ForStatement) {
	if f.Init() != nil {
		t.emitExpression(f.Init())
	}
	label := t.nextTemp()
	condBlock := &BasicBlock{Parent: t.function,
		valueWithName: valueWithName{name: "for.cond." + label}, Index: len(t.function.Blocks)}
	t.function.Blocks = append(t.function.Blocks, condBlock)
	bodyBlock := &BasicBlock{Parent: t.function,
		valueWithName: valueWithName{name: "for.body." + label}, Index: len(t.function.Blocks)}
	t.function.Blocks = append(t.function.Blocks, bodyBlock)
	done := &BasicBlock{Parent: t.function,
		valueWithName: valueWithName{name: "for.done." + label}, Index: len(t.function.Blocks)}
	t.function.Blocks = append(t.function.Blocks, done)

	t.emit(t.goTo(condBlock))

	t.function.current = condBlock
	cond := t.emitExpression(f.Condition())
	t.emit(t.gotoif(cond, bodyBlock, done))

	t.function.current = bodyBlock
	t.Visit(f.Body())
	t.emit(t.goTo(condBlock))

	t.function.current = done
}

func (t *transformer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.DeclarationStatement:
		t.emitDeclaration(n.Decl())

	case *ast.ExpressionStmt:
		t.emitExpression(n.Expr())

	case *ast.FunctionStatement:
		t.emitFunction(n)

	case *ast.BlockStatement:
		for _, stmt := range n.Statements() {
			t.Visit(stmt)
		}

	case *ast.ReturnStatement:
		t.emitReturn(n.Expression())

	case *ast.IfElseStatement:
		t.emitIfElseStatement(n)

	case *ast.ForStatement:
		t.emitForStatement(n)
	default:
		panic(n)
	}

	return nil
}

func Emit(program *ast.Program, ctx *context.Context) *Program {
	p := &Program{Types: make(map[string]types.Type), Globals: make(map[string]Value)}
	t := transformer{program: p,
		factory: types.NewFactory(ctx),
		types:   make(map[string]types.Type),
	}
	for _, decl := range program.Decls() {
		switch n := decl.(type) {
		case *ast.DeclarationStatement:
			if td, ok := n.Decl().(*ast.TypeDeclaration); !ok {
				panic("only type declaration and functions are allowed on global scope")
			} else {
				t.emitDeclaration(td)
			}

		case *ast.TypeDeclaration:
			t.emitDeclaration(n)

		case *ast.FunctionStatement:
			t.emitFunction(n)

		default:
			t.astError(n, "unknown expression on global level "+fmt.Sprint(n))
		}
	}

	return p
}
