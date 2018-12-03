package parse

import (
    "bitbucket.org/dhaliwalprince/funlang/ast"
    "bitbucket.org/dhaliwalprince/funlang/lex"
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
    init := parser.parseAssignExpression()
    return parser.builder.NewDeclaration(pos, v.String(), t, init)
}

func (parser *Parser) parseDeclaration() ast.DeclNode {
    if parser.current.Type() != lex.VAR {
        parser.errs.append(unexpectedToken(parser.current, lex.VAR))
        return nil
    }

    parser.advance()
    return parser.parseDeclarationEpilogue()
}
