/// type checking/resolving logic
package sema
//
//import (
//	"bitbucket.org/dhaliwalprince/funlang/ast"
//	"bitbucket.org/dhaliwalprince/funlang/lex"
//	"bitbucket.org/dhaliwalprince/funlang/types"
//	"fmt"
//)
//
//type typeChecker struct {
//	types map[ast.Node]types.Type
//
//	// since type definition can only exist on the global scope level
//	// so we need only one factory
//	factory types.Factory
//
//	// map name of the type to the type
//	typeNames map[string]types.Type
//
//	unresolved []ast.Node
//
//	topScope *Scope
//
//	errs []error
//}
//
//func (t *typeChecker) init() {
//}
//
//func (c *typeChecker) openScope() {
//	c.topScope = NewScope(c.topScope)
//}
//
//func (c *typeChecker) closeScope() {
//	c.topScope = c.topScope.outer
//}
//
//func (r *typeChecker) resolve(name string, object *ast.Object) *ast.Object {
//	if object == nil {
//		return r.topScope.Lookup(name)
//	}
//	o := r.topScope.PutStrict(name, object)
//	if o != nil {
//		r.appendError(nil, fmt.Sprintf("%s: %s already defined at %s", object.Pos, name, o.Pos))
//	}
//	return o
//}
//
//func (c *typeChecker) saveType(node ast.Node, t types.Type) {
//	c.types[node] = t
//}
//
//func (c *typeChecker) saveNamedType(node ast.Node, name string) {
//	t := c.getNamedType(name)
//	if t == nil {
//		t = c.factory.UnknownType()
//	}
//	c.saveType(node, t)
//}
//
//func (c *typeChecker) appendError(node ast.Node, message string) {
//	if node == nil {
//		c.errs = append(c.errs, fmt.Errorf("%s", message))
//	} else {
//		c.errs = append(c.errs, fmt.Errorf("%d:%d: %s", node.Beg().Row, node.Beg().Col, message))
//	}
//}
//
//func (c *typeChecker) getType(node ast.Node) types.Type {
//	switch n := node.(type) {
//	case *ast.NumericLiteral:
//		return c.factory.IntType()
//
//	case *ast.BooleanLiteral:
//		return c.factory.IntType()
//
//	case *ast.StringLiteral:
//		return c.factory.StringType()
//
//	case *ast.Identifier:
//		return n.Object().Type
//	}
//
//	// for everyone else
//	t, ok := c.types[node]
//	if !ok {
//		return nil
//	}
//
//	return t
//}
//
//func (c *typeChecker) getNamedType(name string) types.Type {
//	t, ok := c.typeNames[name]
//	if !ok { return nil }
//	return t
//}
//
//// implement ast.Visitor
//func (c *typeChecker) Visit(node ast.Node) ast.Visitor {
//	switch n := node.(type) {
//	case *ast.NilLiteral:
//		// not implemented yet
//
//	case *ast.NumericLiteral:
//		n.SetType(c.factory.IntType())
//
//	case *ast.BooleanLiteral:
//		n.SetType(c.factory.IntType())
//
//	case *ast.StringLiteral:
//		n.SetType(c.factory.StringType())
//
//	case *ast.Identifier:
//		// resolve this identifier
//		o := c.resolve(n.Name(), nil)
//		if o == nil {
//			c.unresolved = append(c.unresolved, n)
//		} else {
//			n.SetType(o.Type)
//		}
//
//
//	case *ast.ArgumentList:
//		// we are not handling this here
//		panic("unreachable")
//
//	case *ast.MemberExpression:
//		c.Visit(n.Expr())
//		if n.Expr().Type() == nil {
//			c.unresolved = append(c.unresolved, n)
//			break
//		}
//
//	case *ast.PrefixExpression:
//		c.Visit(n.Expression())
//		t := c.getType(n.Expression())
//		switch n.Op() {
//		case lex.AND:
//			c.saveType(n, c.factory.PointerType(t))
//
//		case lex.MUL:
//			c.saveType(n, types.ToPointerType(t).Elem())
//		}
//
//	case *ast.BinaryExpression:
//		c.Visit(n.Left())
//		c.Visit(n.Right())
//		t1 := c.getType(n.Left())
//		t2 := c.getType(n.Right())
//		if t1 != t2 || (t1.Tag() != types.INT_TYPE || t2.Tag() != types.INT_TYPE) {
//			c.appendError(n, "incompatible types for binary expression")
//		}
//
//		c.saveType(n, t1)
//
//	case *ast.AssignExpression:
//		c.Visit(n.Left())
//		c.Visit(n.Right())
//
//		t1 := c.getType(n.Left())
//		t2 := c.getType(n.Right())
//		if t1 != t2 {
//			c.appendError(n, "different types on either sides of assignment expression")
//		}
//		c.saveType(n, t1)
//
//	case *ast.StructType:
//		fields := make(map[string]types.Type)
//		for _, field := range n.Fields() {
//			c.Visit(field.TypeExpression())
//			fields[field.Name()] = t
//		}
//
//		c.saveType(n, c.factory.StructType(fields))
//
//	case *ast.ArrayType:
//		c.Visit(n.Type())
//		t := c.getType(n.Type())
//		arrayType := c.factory.ArrayType(t)
//		c.saveType(n, arrayType)
//
//	case *ast.BlockStatement:
//		c.openScope()
//		for _, stmt := range n.Statements() {
//			c.Visit(stmt)
//		}
//		c.closeScope()
//
//	case *ast.ForStatement:
//		c.Visit(n.Init())
//		c.Visit(n.Condition())
//		c.Visit(n.Body())
//
//	case *ast.ExpressionStmt:
//		c.Visit(n.Expr())
//
//	case *ast.FunctionStatement:
//
//	}
//
//	return nil
//}

