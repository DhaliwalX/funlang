package parse

import "bitbucket.org/dhaliwalprince/funlang/lex"

type Parser struct {
    lex *lex.Lexer
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
