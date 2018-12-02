package parse

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"errors"
)

func (parser *Parser) parseForStatement() ast.Statement {
	parser.advance()


}

func (parser *Parser) parseFunction() ast.Statement {
	pos := parser.current.Begin()
	parser.advance()
	name := ""
	if parser.current.Type() == lex.IDENT {
		name = parser.current.Value()
	}

	if parser.current.Type() != lex.LBRACK {
		parser.errs.append(newParseError(parser.current, "expected a ("))
		return nil
	}

	parser.advance()
	// parse arguments
	params := []ast.DeclNode{}
	stop := false
	for {
		decl := parser.parseDeclarationEpilogue()
		if decl == nil {
			return nil
		}
		if parser.current.Type() != lex.COMMA && parser.current.Type() == lex.RBRACK {
			parser.advance()
			stop = true
		}

		if parser.current.Type() != lex.COMMA {
			parser.errs.append(unexpectedToken(parser.current, lex.RBRACK))
			return nil
		}

		params = append(params, decl)

		if stop {
			break
		}
	}

	ret := parser.parseType()
	proto := parser.builder.NewFunctionProtoType(pos, name, params, ret)
	body := parser.parseBlockStatement()
	if body == nil {
		parser.errs.append(newParseError(parser.current, "expected function body"))
		return nil
	}

	return parser.builder.NewFunctionStatement(proto, body.(*ast.BlockStatement))
}

func (parser *Parser) parseIfStatement() ast.Statement {
	pos := parser.current.Begin()
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

}

func (parser *Parser) parseExpressionStatement() ast.Statement {
	if parser.current.Type() == lex.SEMICOLON {
		return nil
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

	return parser.builder.NewExpressionStatement(expr)
}

func (parser *Parser) parseBlockStatement() ast.Statement {

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
	default:
		return parser.parseExpressionStatement()
	}
}
