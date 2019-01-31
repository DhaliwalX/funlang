package parse

import (
	"funlang/ast"
	"funlang/context"
	"funlang/lex"
)

type Parser struct {
	lex     *lex.Lexer
	current lex.Token
	errs    errorList
	ctx     *context.Context
	builder *ast.Builder
}

func (parser *Parser) advance() {
	parser.current = parser.lex.Next()
}

func (parser *Parser) advanceTil(t lex.TokenType) {
	for parser.current.Type() != t {
		parser.advance()
	}
}

func NewParserFromString(ctx *context.Context, source string) *Parser {
	src := lex.NewStringSource(source)
	lexer := lex.NewLexer(src)
	parser := Parser{lex: lexer, ctx: ctx, builder: ast.NewBuilder(ctx)}
	return &parser
}

func (parser *Parser) Parse() (*ast.Program, error) {
	parser.advance()

	decls := []ast.Node{}
	for {
		if parser.current.Type() == lex.ILLEGAL || parser.current.Type() == lex.EOF {
			break
		}
		a := parser.parseTopLevelNode()
		if len(parser.errs.list) > 0 {
			return nil, &parser.errs
		}

		decls = append(decls, a)
	}
	return ast.NewProgram(parser.ctx, parser.lex.Source(), decls), nil
}

func NewParserFromFile(ctx *context.Context, name string) *Parser {
	src := lex.NewFileSource(name)
	lexer := lex.NewLexer(src)
	parser := Parser{lex: lexer, ctx: ctx, builder: ast.NewBuilder(ctx)}
	return &parser
}
