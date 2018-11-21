package ast

import "bitbucket.org/dhaliwalprince/funlang/lex"

// Node is the base interface for all the expression
// and statement types.
type Node interface {
    VisitorAcceptor
    Beg() lex.Position
    End() lex.Position
}
