package parse

import (
	"bitbucket.org/dhaliwalprince/funlang/ast"
	"bitbucket.org/dhaliwalprince/funlang/lex"
)

func (parser *Parser) parseField() *ast.Field {
	name := parser.parseIdentifier()
	t := parser.parseType()
	return parser.builder.NewField(name, t)
}

func (parser *Parser) parseStructType() *ast.StructType {
	pos := parser.current.Begin()
	parser.advance()
	fields := []*ast.Field{}
	for parser.current.Type() != lex.RBRACE {
		if parser.current.Type() != lex.IDENT {
			parser.errs.append(unexpectedToken(parser.current, lex.IDENT))
			parser.advanceTil(lex.RBRACE)
			break
		}

		field := parser.parseField()
		fields = append(fields, field)
	}

	return parser.builder.NewStructType(pos, fields)
}

func (parser *Parser) parseArrayType() ast.Expression {
	pos := parser.current.Begin()
	parser.advance()
	if parser.current.Type() == lex.RBRACK {
		parser.advance()
		return parser.parseType()
	}

	expr := parser.parseExpression()
	t := parser.parseType()
	return parser.builder.NewArrayType(pos, expr, t)
}

func (parser *Parser) parseType() ast.Expression {
	pos := parser.current.Begin()
	val := parser.current.Value()
	switch parser.current.Type() {
	case lex.IDENT:
		parser.advance()
		return parser.builder.NewIdentifier(pos, val)

	case lex.LBRACK:
		return parser.parseArrayType()

	case lex.STRUCT:
		parser.advance()
		return parser.parseStructType()

	default:
		return nil
	}
}

// parsing type declarations
// e.g.
// type Person struct {
//		firstName string
//		age		  int
// }
// type MyInt int
func (parser *Parser) parseTypeDeclaration() *ast.TypeDeclaration {
	if parser.current.Type() != lex.TYPE {
		parser.errs.append(unexpectedToken(parser.current, lex.TYPE))
		return nil
	}
	pos := parser.current.Begin()
	parser.advance()
	name := parser.parseIdentifier()
	t := parser.parseType()
	return parser.builder.NewTypeDeclaration(pos, name, t)
}
