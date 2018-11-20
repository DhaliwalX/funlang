package ast

import "bitbucket.org/dhaliwalprince/funlang/lex"

// Node is the base interface for all the expression
// and statement types.
type Node interface {
    Beg() lex.Position
    End() lex.Position
}
