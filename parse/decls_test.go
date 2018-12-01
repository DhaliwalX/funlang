package parse

import (
    "bitbucket.org/dhaliwalprince/funlang/ast"
    "testing"
)

func TestParseDeclarationWithoutInit(t *testing.T) {
    parser := newParser("var a int")
    parser.advance()
    n := parser.parseDeclaration()
    if n == nil {
        t.Error("parsing errors", parser.errs.Error())
    }
    if d, ok := n.(*ast.Declaration); ok {
        t.Log(d)
    } else {
        t.Error("parsing expected a declaration", parser.errs.Error())
    }
}
