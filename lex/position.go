package lex

import "fmt"

// Position is used for representing the location in source code
type Position struct {
    Col int
    Row int
}

func (p Position) String() string {
    return fmt.Sprintf("%d:%d", p.Row+1, p.Col+1)
}

var NO_POS = Position{Col:-1, Row: -1}