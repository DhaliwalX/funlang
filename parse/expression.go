package parse

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/lex"
)

func (parser *Parser) parseNumber() *ast.NumericLiteral {
	pos := parser.current.Begin()
	val := parser.current.Value()
	parser.advance()
	return parser.builder.NewNumericLiteral(pos, val, true)
}

func (parser *Parser) parseString() *ast.StringLiteral {
	pos := parser.current.Begin()
	val := parser.current.Value()
	parser.advance()
	return parser.builder.NewStringLiteral(pos, val)
}

func (parser *Parser) parseIdentifier() *ast.Identifier {
	if parser.current.Type() != lex.IDENT {
		parser.errs.append(unexpectedToken(parser.current, lex.IDENT))
		return nil
	}

	pos := parser.current.Begin()
	val := parser.current.Value()
	parser.advance()
	return parser.builder.NewIdentifier(pos, val)
}

func (parser *Parser) parsePrimaryExpression() ast.Expression {
	switch parser.current.Type() {
	case lex.INT:
		return parser.parseNumber()
	case lex.FLOAT:
		return parser.parseNumber()
	case lex.STRING:
		return parser.parseString()
	case lex.IDENT:
		return parser.parseIdentifier()
	case lex.TYPE:
		return parser.parseType()
	case lex.LPAREN:
		parser.advance()
		x := parser.parseExpression()
		if x == nil {
			parser.advanceTil(lex.RPAREN)
		} else if parser.current.Type() != lex.RPAREN {
			parser.errs.append(newParseError(parser.current, "expected a )"))
		}
		parser.advance()
		return x
	default:
		parser.errs.append(newParseError(parser.current, "unexpected token at primary"))
		return nil
	}
}

func (parser *Parser) parseArgumentList() ast.Expression {
	pos := parser.current.Begin()
	var x []ast.Expression
	for {
		parser.advance()
		if parser.current.Type() == lex.RPAREN {
			parser.advance()
			break
		}
		x = append(x, parser.parseExpression())
		if parser.current.Type() == lex.COMMA {
			parser.advance()
		} else {
			parser.errs.append(newParseError(parser.current, "expected a ,"))
			parser.advanceTil(lex.RPAREN)
			return nil
		}
		if parser.current.Type() == lex.RPAREN {
			parser.errs.append(newParseError(parser.current, "unexpected )"))
			parser.advance()
			x = append(x, nil)
			break
		}
	}

	return parser.builder.NewArgumentList(pos, x)
}

func (parser *Parser) parseMemberExpression() ast.Expression {
	primary := parser.parsePrimaryExpression()
	if primary == nil {
		return primary
	}

	for {
		switch parser.current.Type() {
		case lex.PERIOD:
			parser.advance()
			member := parser.parsePrimaryExpression()
			primary = parser.builder.NewMemberExpression(lex.PERIOD, primary, member)
		case lex.ARROW:
			parser.advance()
			member := parser.parsePrimaryExpression()
			primary = parser.builder.NewMemberExpression(lex.ARROW, primary, member)
		case lex.LBRACK:
			parser.advance()
			member := parser.parseExpression()
			x := parser.builder.NewMemberExpression(lex.LBRACK, primary, member)
			if parser.current.Type() != lex.RBRACK {
				parser.errs.append(newParseError(parser.current, "expected a ]"))
				return nil
			}
			parser.advance()
			primary = x

		case lex.LPAREN:
			member := parser.parseArgumentList()
			primary = parser.builder.NewMemberExpression(lex.LPAREN, primary, member)
		default:
			return primary
		}
	}
}

func (parser *Parser) parsePrefixExpression() ast.Expression {
	tok := parser.current
	switch tok.Type() {
	case lex.AND:
	case lex.MUL:
		parser.advance()
		x := parser.parseMemberExpression()
		if x == nil {
			return nil
		}
		return parser.builder.NewPrefixExpression(tok.Begin(), tok.Type(), x)

	case lex.SUB:
		parser.advance()
		x := parser.parseMemberExpression()
		if x == nil {
			return nil
		}
		return parser.builder.NewBinaryExpression(tok.Begin(), lex.MUL, parser.builder.NewInt(-1),
			parser.parseMemberExpression())
	default:
		return parser.parseMemberExpression()
	}

	return nil
}

func (parser *Parser) parsePostfixExpression() ast.Expression {
	return parser.parsePrefixExpression()
}

func (parser *Parser) parseBinaryRHS(lowestPrec int, tok lex.Token, left ast.Expression) ast.Expression {
	for {
		if tok.Type().Precedence() <= lowestPrec {
			return left
		}

		parser.advance()
		right := parser.parsePrefixExpression()

		if parser.current.Type().Precedence() >= tok.Type().Precedence() {
			right = parser.parseBinaryRHS(lowestPrec+1, parser.current, right)
		}

		left = parser.builder.NewBinaryExpression(tok.Begin(), tok.Type(), left, right)
		tok = parser.current
	}
}

func (parser *Parser) parseBinaryExpression() ast.Expression {
	left := parser.parsePostfixExpression()
	tok := parser.current
	return parser.parseBinaryRHS(lex.LowestPrec, tok, left)
}

func (parser *Parser) parseAssignExpression() ast.Expression {
	left := parser.parseBinaryExpression()
	if parser.current.Type() == lex.ASSIGN {
		parser.advance()
		right := parser.parseAssignExpression()
		return parser.builder.NewAssignExpression(left.Beg(), left, right)
	}
	return left
}

func (parser *Parser) parseExpression() ast.Expression {
	return parser.parseAssignExpression()
}
