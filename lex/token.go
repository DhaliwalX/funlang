package lex

// Token represents one lexeme
type Token struct {
	// type of the token
	_type TokenType

	// position where this token was found
	beg Position
	end Position

	// string value for this token found in the source
	value string
}

func (token Token) Type() TokenType {
	return token._type
}

func (toke Token) Begin() Position {
	return toke.beg
}

func (token Token) End() Position {
	return token.end
}

func (token Token) Value() string {
	return token.value
}
