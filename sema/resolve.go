// resolve.go - resolves all the names present in the ast
package sema

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"fmt"
)

// resolves all the symbols of the program and links them
type resolver struct {
	topScope *Scope
	unresolved []*ast.Identifier

	errs []error
}

func (r *resolver) openScope() {
	r.topScope = NewScope(r.topScope)
}

func (r *resolver) appendError(err error) {
	r.errs = append(r.errs, err)
}

func (r *resolver) closeScope() {
	r.topScope = r.topScope.outer
}

func (r *resolver) resolve(name string, object *ast.Object) *ast.Object {
	if object == nil {
		return r.topScope.Lookup(name)
	}
	o := r.topScope.PutStrict(name, object)
	if o != nil {
		r.appendError(fmt.Errorf("%s: %s already defined at %s", object.Pos, name, o.Pos))
	}
	return o
}

func ResolveProgram(program *ast.Program) []error {
	r := &resolver{}
	r.openScope()
	for _, decl := range program.Decls() {
		resolve2(r, decl)
	}

	for _, unresolved := range r.unresolved {
		o := r.resolve(unresolved.Name(), nil)
		if o == nil {
			r.appendError(fmt.Errorf("%s: %s undefined", unresolved.Beg(), unresolved.Name()))
		} else {
			unresolved.Object = o
		}
	}
	return r.errs
}

func resolve2(r *resolver, node ast.Node) {
	if node == nil {
		return
	}
	ast.Walk(r, node)
}

func makeObject(kind ast.ObjKind, data interface{}, position lex.Position) *ast.Object {
	switch kind {
	case ast.VAR:
		return &ast.Object{Kind:kind, Decl: data, Pos:position}

	case ast.TYPE:
		return &ast.Object{Kind:kind, Type:data, Pos:position}

	default:
		return &ast.Object{Kind:ast.DONT_KNOW}
	}
}

func (r *resolver) resolveMemberExpression(m *ast.MemberExpression) {
	resolve2(r, m.Expr())
	switch m.AccessKind() {
	case lex.PERIOD:
		if _, ok := m.Member().(*ast.Identifier); !ok {
			r.appendError(fmt.Errorf("%s:expected an identifier", m.Member().Beg()))
		}

	case lex.LBRACK:
		resolve2(r, m.Member())

	case lex.LPAREN:
		resolve2(r, m.Member())

	default:
		panic("undefined operation"+m.AccessKind().String())
	}
}

func (r *resolver) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Declaration:
		o := ast.Object{Type:n.Type(), Decl:n, Kind:ast.VAR, Name: n.Name(), Pos: n.Beg()}
		r.resolve(n.Name(), &o)

		if n.Init() != nil {
			resolve2(r, n.Init())
		}

		if n.Type() != nil {
			resolve2(r, n.Type())
		}

	case *ast.Identifier:
		if o := r.resolve(n.Name(), nil); o == nil {
			r.unresolved = append(r.unresolved, n)
		} else {
			n.Object = o
		}

	case *ast.TypeDeclaration:
		o := &ast.Object{Type:n.Type(), Pos:n.Beg(), Decl:n, Kind:ast.TYPE, Name: n.Name()}
		r.resolve(n.Name(), o)
		n.Ident().Object = o
		resolve2(r, n.Type())

	case *ast.StructType:
		for _, field := range n.Fields() {
			resolve2(r, field)
		}

	case *ast.Field:
		resolve2(r, n.Type())

	case *ast.ArrayType:
		resolve2(r, n.Type())

	case *ast.NilLiteral:
	case *ast.NumericLiteral:
	case *ast.StringLiteral:
	case *ast.BooleanLiteral:
		// do nothing
		func(){}()

	case *ast.ArgumentList:
		for _, arg := range n.Exprs() {
			resolve2(r, arg)
		}

	case *ast.MemberExpression:
		r.resolveMemberExpression(n)

	case *ast.PrefixExpression:
		resolve2(r, n.Expression())

	case *ast.PostfixExpression:
		resolve2(r, n.Expression())

	case *ast.BinaryExpression:
		resolve2(r, n.Left())
		resolve2(r, n.Right())

	case *ast.AssignExpression:
		resolve2(r, n.Left())
		resolve2(r, n.Right())

	case *ast.BlockStatement:
		r.openScope()
		for _, stmt := range n.Statements() {
			resolve2(r, stmt)
		}
		r.closeScope()

	case *ast.ForStatement:
		resolve2(r, n.Init())
		resolve2(r, n.Condition())
		resolve2(r, n.Body())

	case *ast.ExpressionStmt:
		resolve2(r, n.Expr())

	case *ast.FunctionStatement:
		o := &ast.Object{Kind:ast.VAR, Name: n.Proto().Name(), Decl:n, Type:n.Proto()}
		r.resolve(n.Proto().Name(), o)
		resolve2(r, n.Proto().Return())
		r.openScope()
		for _, param := range n.Proto().Params() {
			resolve2(r, param)
		}

		for _, stmt := range n.Body().Statements() {
			resolve2(r, stmt)
		}
		r.closeScope()

	case *ast.DeclarationStatement:
		resolve2(r, n.Decl())

	case *ast.IfElseStatement:
		resolve2(r, n.Condition())
		resolve2(r, n.Body())
		resolve2(r, n.ElseNode())

	case *ast.ReturnStatement:
		resolve2(r, n.Expression())

	default:
		panic(fmt.Sprintf("%T is not handled by resolver", node))
	}

	return nil
}

