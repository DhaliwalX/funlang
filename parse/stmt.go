package parse

import (
	"errors"
	"funlang/ast"
	"funlang/lex"
)

func (parser *Parser) parseForStatement() ast.Statement {
	pos := parser.current.Begin()
	parser.advance()
	var init, condition ast.Expression
	init = parser.parseExpression()
	if parser.current.Type() == lex.SEMICOLON {
		parser.advance()
		condition = parser.parseExpression()
	} else {
		condition = init
		init = nil
	}

	if parser.current.Type() != lex.LBRACE {
		parser.errs.append(unexpectedToken(parser.current, lex.LBRACE))
		return nil
	}

	body := parser.parseBlockStatement()
	if body == nil {
		parser.errs.append(errors.New("for body cannot be empty"))
		return nil
	}

	return parser.builder.NewForStatement(pos, init, condition, body)
}

func (parser *Parser) parseFunction() ast.Statement {
	pos := parser.current.Begin()
	parser.advance()
	name := ""
	if parser.current.Type() == lex.IDENT {
		name = parser.current.Value()
		parser.advance()
	}

	if parser.current.Type() != lex.LPAREN {
		parser.errs.append(newParseError(parser.current, "expected a ("))
		return nil
	}

	parser.advance()
	// parse arguments
	params := []ast.DeclNode{}
	for {
		if parser.current.Type() == lex.RPAREN {
			parser.advance()
			break
		}
		decl := parser.parseDeclarationEpilogue()
		if decl == nil {
			return nil
		}
		if parser.current.Type() == lex.RPAREN {
			params = append(params, decl)
			parser.advance()
			break
		}

		if parser.current.Type() != lex.COMMA {
			parser.errs.append(unexpectedToken(parser.current, lex.RPAREN))
			return nil
		}

		parser.advance()

		params = append(params, decl)
	}

	ret := parser.parseType()
	proto := parser.builder.NewFunctionProtoType(pos, name, params, ret)

	var body *ast.BlockStatement
	if parser.current.Type() != lex.SEMICOLON {
		b := parser.parseBlockStatement()
		if b == nil {
			parser.errs.append(newParseError(parser.current, "expected a function body"))
			return nil
		}
		body = b.(*ast.BlockStatement)
	} else {
		parser.advance()
	}

	fun := parser.builder.NewFunctionStatement(proto, body)
	return fun
}

func (parser *Parser) parseIfStatement() ast.Statement {
	parser.advance()
	condition := parser.parseExpression()
	if condition == nil {
		parser.errs.append(errors.New("malformed if statement"))
		return nil
	}

	body := parser.parseStatement()
	if body == nil {
		parser.errs.append(errors.New("wrong if structure"))
		return nil
	}

	var elseNode ast.Statement = nil
	if parser.current.Type() == lex.ELSE {
		parser.advance()
		elseNode = parser.parseStatement()
		if elseNode == nil {
			parser.errs.append(errors.New("malformed else node"))
			return nil
		}
	}

	return parser.builder.NewIfStatement(condition, body, elseNode)
}

func (parser *Parser) parseReturnStatement() ast.Statement {
	if parser.current.Type() != lex.RETURN {
		parser.errs.append(newParseError(parser.current, "expected 'return'"))
		return nil
	}
	pos := parser.current.Begin()
	parser.advance()
	expr := parser.parseExpression()
	if parser.current.Type() != lex.SEMICOLON {
		parser.errs.append(newParseError(parser.current, "expected a semicolon"))
		return nil
	}
	parser.advance()
	return parser.builder.NewReturnStatement(pos, expr)
}

func (parser *Parser) parseExpressionStatement() ast.Statement {
	if parser.current.Type() == lex.SEMICOLON {
		parser.advance()
		return parser.builder.NewExpressionStatement(nil)
	}

	expr := parser.parseExpression()
	if expr == nil {
		parser.errs.append(newParseError(parser.current, "expected an expression statement"))
		parser.advanceTil(lex.SEMICOLON)
	}

	if parser.current.Type() != lex.SEMICOLON {
		parser.errs.append(newParseError(parser.current, "expected a semi-colon ;"))
		return nil
	}

	parser.advance()

	return parser.builder.NewExpressionStatement(expr)
}

func (parser *Parser) parseBlockStatement() ast.Statement {
	list := []ast.Statement{}
	parser.advance()
	for {
		if parser.current.Type() == lex.RBRACE {
			parser.advance()
			break
		}

		stmt := parser.parseStatement()
		if stmt == nil {
			return nil
		}
		list = append(list, stmt)
	}

	return parser.builder.NewBlockStatement(list)
}

func (parser *Parser) parseDeclarationStatement() ast.Statement {
	declNode := parser.parseDeclaration()
	if declNode == nil {
		return nil
	}

	if parser.current.Type() != lex.SEMICOLON {
		parser.errs.append(newParseError(parser.current, "expected a semicolon"))
		return nil
	}

	parser.advance()
	return parser.builder.NewDeclarationStatement(declNode)
}

func (parser *Parser) parseStatement() ast.Statement {
	switch parser.current.Type() {
	case lex.FOR:
		return parser.parseForStatement()

	case lex.FUNC:
		return parser.parseFunction()

	case lex.IF:
		return parser.parseIfStatement()

	case lex.RETURN:
		return parser.parseReturnStatement()

	case lex.LBRACE:
		return parser.parseBlockStatement()

	case lex.VAR:
		return parser.parseDeclarationStatement()

	default:
		return parser.parseExpressionStatement()
	}
}

func (parser *Parser) parseTopLevelNode() ast.Node {
	switch parser.current.Type() {
	case lex.FUNC:
		return parser.parseFunction()

	case lex.VAR:
		return parser.parseDeclarationStatement()

	case lex.TYPE:
		return parser.parseTypeDeclaration()
	default:
		parser.errs.append(newParseError(parser.current, "expected a function or var or type declaration"))
		return nil
	}
}
