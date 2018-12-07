package parse

import (
    "bitbucket.org/dhaliwalprince/funlang/ast"
    "bitbucket.org/dhaliwalprince/funlang/context"
    "bitbucket.org/dhaliwalprince/funlang/lex"
    "fmt"
)

type Parser struct {
    lex *lex.Lexer
    current lex.Token
    errs errorList
    ctx *context.Context
    builder *ast.Builder

    topScope *ast.Scope
    scopes map[ast.Node]*ast.Scope
    unresolved []*ast.Identifier
}

func (parser *Parser) openScope() {
    parser.topScope = ast.NewScope(parser.topScope)
}

func (parser *Parser) closeScope() {
    parser.topScope = parser.topScope.Outer()
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
    parser.openScope()

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
    parser.closeScope()
    return ast.NewProgram(parser.ctx, parser.lex.Source(), decls), nil
}

func newAlreadyDefinedError(o *ast.Object, n ast.Node) error {
    return fmt.Errorf("%s is already defined at %s (duplicate definition at %s)", o.Name, o.Pos, n.Beg())
}

func (parser *Parser) resolve(node ast.Node) *ast.Object {
    ident, isIdent := node.(*ast.Identifier)
    if !isIdent {
        return nil
    }

    o := &ast.Object{Name:ident.Name(), Kind:ast.VAR, Decl:ident}
    old := parser.topScope.PutStrict(ident.Name(), o)
    if old != nil {
        return old
    }
    return o
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
