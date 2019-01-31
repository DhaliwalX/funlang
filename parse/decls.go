package parse

import (
	"funlang/ast"
	"funlang/lex"
)

func (parser *Parser) parseDeclarationEpilogue() ast.DeclNode {
	pos := parser.current.Begin()
	v := parser.parseIdentifier()
	if v == nil {
		parser.errs.append(newParseError(parser.current, "expected an identifier"))
		return nil
	}
	t := parser.parseType()
	if t == nil {
		if parser.current.Type() != lex.ASSIGN {
			parser.errs.append(newParseError(parser.current, "expected a type or ="))
			return nil
		}
	}

	if parser.current.Type() == lex.ASSIGN {
		parser.advance()
	}
	if parser.current.Type() == lex.SEMICOLON ||
		parser.current.Type() == lex.COMMA ||
		parser.current.Type() == lex.RPAREN {
		return parser.builder.NewDeclaration(pos, v.String(), t, nil)
	}
	init := parser.parseAssignExpression()
	node := parser.builder.NewDeclaration(pos, v.String(), t, init)
	return node
}

func (parser *Parser) parseDeclaration() ast.DeclNode {
	if parser.current.Type() != lex.VAR {
		parser.errs.append(unexpectedToken(parser.current, lex.VAR))
		return nil
	}

	parser.advance()
	return parser.parseDeclarationEpilogue()
}
