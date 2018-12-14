// defines interface which will be implemented by every instruction
package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/types"
	"fmt"
	"strings"
)

type Instruction interface {
	Value
	Operands() []Value
	Operand(i int) Value
}

type instrNode struct {}

func (i instrNode) Tag() ValueTag {
	return INSTRUCTION
}


type instrWithOperands struct {
	operands []Value
}

func (i instrWithOperands) AddOperand(operand Value) {
	i.operands = append(i.operands, operand)
}

func (in instrWithOperands) Operand(i int) Value {
	return in.operands[i]
}

func (i instrWithOperands) Operands() []Value {
	return i.operands
}

func (i instrWithOperands) Uses() []Value {
	return i.operands
}

// memory instructions
// 		x:*int = new int	alloc
//		a:int = @x:*int				load
//		@x:*int = 10:int				store
//

type AllocInstr struct {
	instrNode
	valueWithUsers
	valueWithName
	t types.Type
}

type LoadInstr struct {
	instrNode
	valueWithUsers
	instrWithOperands
	valueWithName
	t types.Type
	users []Value
}

type StoreInstr struct {
	instrNode
	instrWithOperands
}

func (a *AllocInstr) Uses() []Value {
	return []Value{}
}

func (a *AllocInstr) AddUse(use Value) {
	return
}

func (a *AllocInstr) Operands() []Value {
	return []Value{}
}

func (a *AllocInstr) AddOperand(operand Value){
	return
}

func (a *AllocInstr) Operand(i int) Value {
	return nil
}

func (a *AllocInstr) Type() types.Type {
	return typeFactory.PointerType(a.t)
}

func (a *AllocInstr) ShortString() string {
	return fmt.Sprintf("%s:%s", a.Name(), a.Type())
}

func (a *AllocInstr) String() string {
	return fmt.Sprintf("%s = new %s", a.ShortString(), a.t)
}

func (l *LoadInstr) Type() types.Type {
	return l.Operand(1).Type().Elem()
}

func (l *LoadInstr) ShortString() string {
	return fmt.Sprintf("%s:%s", l.Name(), l.Type())
}

func (l *LoadInstr) String() string {
	return fmt.Sprintf("%s:%s = %s", l.ShortString(), l.Operand(0).ShortString())
}

func (s *StoreInstr) Type() types.Type {
	return nil
}

func (s *StoreInstr) ShortString() string {
	return ""
}

func (s *StoreInstr) String() string {
	return fmt.Sprintf("%s = %s", s.Operand(0).ShortString(),
		s.Operand(1).ShortString())
}

func (s *StoreInstr) Users() []Value {
	return []Value{}
}

func (s *StoreInstr) AddUser(user Value) {
	return
}

func (s *StoreInstr) Name() string {
	return ""
}

// member access instruction
//
type MemberInstr struct {
	instrNode
	instrWithOperands
	valueWithName
	valueWithUsers
}

func (m *MemberInstr) Type() types.Type {
	return m.Operand(0).Type().Field(m.Operand(1).(*ConstantString).Value)
}

func (m *MemberInstr) ShortString() string {
	return fmt.Sprintf("%s:%s", m.Name(), m.Type())
}

func (m *MemberInstr) String() string {
	return fmt.Sprintf("%s = member %s, %s",
		m.ShortString(), m.Operand(0).ShortString(), m.Operand(1).ShortString())
}

type IndexInstr struct {
	instrNode
	instrWithOperands
	valueWithName
	valueWithUsers
}

func (i *IndexInstr) Type() types.Type {
	return i.Operand(0).Type().Elem()
}

func (i *IndexInstr) ShortString() string {
	return fmt.Sprintf("%s:%s", i.Name(), i.Type())
}

func (i *IndexInstr) String() string {
	return fmt.Sprintf("%s = index %s, %s",
		i.ShortString(), i.Operand(0).ShortString(), i.Operand(1).ShortString())
}

// arithematic instructions
/*   x = a + b
 *   y = x * c
 *   +, -, *, /, %, ^, &, |, <, >, ==
 */
type ArithOpcode int
const (
	PLUS ArithOpcode = iota
	MINUS
	MUL
	DIV
	MOD
	XOR
	AND
	OR
	LT
	GT
	EQ
	LAND
	LOR
)

type ArithInstr struct {
	valueWithUsers
	instrNode
	instrWithOperands
	valueWithName
	opCode ArithOpcode
}

func (a *ArithInstr) Type() types.Type {
	return a.Operand(0).Type()
}

func (a *ArithInstr) ShortString() string {
	return fmt.Sprintf("%s:%s", a.Name(), a.Type())
}

func (a *ArithInstr) String() string {
	return fmt.Sprintf("%s = %s",a.Operand(0).ShortString(), a.Operand(1).ShortString())
}

// control flow instructions
//  if true:int goto $label else goto $label
//  r:int = %add(a:int, b:int)
//  goto $label
type ConditionalGoto struct {
	instrNode
	instrWithOperands
	valueWithNoName
}

func (c *ConditionalGoto) Users() []Value {
	return []Value{}
}

func (c *ConditionalGoto) AddUser(user Value) {
	return
}

func (c *ConditionalGoto) Type() types.Type {
	return nil
}

func (c *ConditionalGoto) ShortString() string {
	return c.String()
}

func (c *ConditionalGoto) String() string {
	return fmt.Sprintf("if %s goto %s else %s",
		c.Operand(0).ShortString(), c.Operand(1).ShortString(), c.Operand(2).ShortString())
}

type CallInstr struct {
	instrNode
	instrWithOperands
	valueWithUsers
	valueWithName
}

func (c *CallInstr) Type() types.Type {
	return types.ToFunctionType(c.Operand(0).Type()).ReturnType()
}

func (c *CallInstr) ShortString() string {
	return fmt.Sprintf("%s:%s", c.Name(), c.Type())
}

func (c *CallInstr) String() string {
	builder := strings.Builder{}
	builder.WriteString(c.ShortString())
	builder.WriteString(" = ")
	builder.WriteString(c.Operand(0).ShortString())
	builder.WriteString("(")

	l := len(c.operands)-1
	for i, arg := range c.Operands()[1:] {
		builder.WriteString(arg.ShortString())
		if l == i+1 {
			break
		}
		builder.WriteString(", ")
	}
	builder.WriteString(")")
	return builder.String()
}

type UnconditionalGoto struct {
	instrNode
	instrWithOperands
	valueWithNoName
}

func (u *UnconditionalGoto) Type() types.Type {
	return nil
}

func (u *UnconditionalGoto) ShortString() string {
	return u.String()
}

func (u *UnconditionalGoto) String() string {
	return fmt.Sprintf("goto %s", u.Operand(0).ShortString())
}

func (u *UnconditionalGoto) Users() []Value {
	return []Value{}
}

func (u *UnconditionalGoto) AddUser(user Value) {
	return
}

// phi node
// x = phi [x1 <- $b1, x2 <- $b2, ... ]
type PhiNode struct {
	instrNode
	instrWithOperands
	valueWithName
	valueWithUsers
}

func (p *PhiNode) Type() types.Type {
	return p.Operand(0).Type()
}

func (p *PhiNode) ShortString() string {
	return fmt.Sprintf("%s:%s", p.Name(), p.Type())
}

func (p *PhiNode) String() string {
	base := fmt.Sprintf("%s = phi [", p.ShortString())
	builder := strings.Builder{}
	builder.WriteString(base)
	l := len(p.Operands())
	for i, n := range p.Operands() {
		builder.WriteString(n.ShortString())
		if i+1 == l {
			break
		}

		builder.WriteString(", ")
	}

	builder.WriteString(")")
	return builder.String()
}
