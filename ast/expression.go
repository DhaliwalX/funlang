package ast

import "bitbucket.org/dhaliwalprince/funlang/lex"

type Expression struct {
    beg lex.Position
    end lex.Position
}

func (expr *Expression) Beg() lex.Position {
    return expr.beg
}

func (expr *Expression) End() lex.Position {
    return expr.end
}

type NumericLiteral struct {
    Expression
    val string
    isFloating bool
}

func (literal *NumericLiteral) Value() string {
    return literal.val
}

func (literal *NumericLiteral) IsFloating() bool {
    return literal.isFloating
}

type StringLiteral struct {
    Expression
    val string
}

func (literal *StringLiteral) Value() string {
    return literal.val
}
