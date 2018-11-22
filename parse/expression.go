package parse

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/lex"
)

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

func (parser *Parser) parseExpression() ast.Expression {
	return nil
}
