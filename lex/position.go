package lex

// Position is used for representing the location in source code
type Position struct {
    Col int
    Row int
}

var NO_POS = Position{Col:-1, Row: -1}