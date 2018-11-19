package ast

import "bitbucket.org/dhaliwalprince/funlang/lex"

type exprNode interface {
    Node
    expr()
}

type nilLiteral struct {
    pos lex.Position // since nil has three characters only
}

type numericLiteral struct {
    pos lex.Position
    val string
    isFloating bool
}

type stringLiteral struct {
    pos lex.Position
    val string
}

type arrayLiteral struct {
    vals []exprNode
}

type structProp struct {
    pos lex.Position
    name string
    val exprNode
}

type structLiteral struct {
    pos lex.Position
    end lex.Position
    props []structProp
}

type identifier struct {
    pos lex.Position
    name string
}

type booleanLiteral struct {
    pos lex.Position
    val bool
}

type argumentList struct {
    pos lex.Position
    exprs []exprNode
}

type memberExpression struct {
    pos lex.Position
    token lex.TokenType
    member exprNode
    expr exprNode
}

type prefixExpression struct {
    pos lex.Position
    op lex.TokenType
    expr exprNode
}

type postfixExpression struct {
    pos lex.Position
    op lex.TokenType
    expr exprNode
}

type binaryExpression struct {
    pos lex.Position
    op lex.TokenType
    left exprNode
    right exprNode
}

type assignExpression struct {
    pos lex.Position
    left exprNode
    right exprNode
}
