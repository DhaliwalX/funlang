package parse

import (
    "bitbucket.org/dhaliwalprince/funlang/ast"
    "bitbucket.org/dhaliwalprince/funlang/context"
    "bitbucket.org/dhaliwalprince/funlang/lex"
)

type Parser struct {
    lex *lex.Lexer
    current lex.Token
    errs errorList
    ctx *context.Context
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
    parser := Parser{lex:lexer, ctx: ctx, builder: ast.NewBuilder(ctx)}
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

/*
    M(NullLiteral)          \
    M(ThisHolder)           \
    M(IntegralLiteral)      \
    M(StringLiteral)        \
    M(ArrayLiteral)         \
    M(ObjectLiteral)        \
    M(Identifier)           \
    M(BooleanLiteral)       \
    M(ArgumentList)         \
    M(CallExpression)       \
    M(MemberExpression)     \
    M(PrefixExpression)     \
    M(PostfixExpression)    \
    M(BinaryExpression)     \
    M(AssignExpression)     \
    M(Declaration)          \
    M(DeclarationList)      \
    M(IfStatement)          \
    M(IfElseStatement)      \
    M(ForStatement)         \
    M(WhileStatement)       \
    M(DoWhileStatement)     \
    M(BlockStatement)       \
    M(FunctionPrototype)    \
    M(FunctionStatement)    \
    M(ReturnStatement)
*/
