package ast

import "bitbucket.org/dhaliwalprince/funlang/lex"

// ASTNode is the base interface for all the expression
// and statement types.
type ASTNode interface {
    Beg() lex.Position
    End() lex.Position
}
