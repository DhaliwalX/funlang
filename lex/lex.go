package lex

import (
    "strings"
    "unicode"
)

type Lexer struct {
    // represents the current position in the source code
    // with col and row information
    position Position

    // source from which we are reading source code
    source Source

    current rune
}

func NewLexer(source Source) *Lexer {
    lex := &Lexer{source: source}
    lex.next()
    return lex
}

func (lexer *Lexer) Source() Source {
    return lexer.source
}

// nextLine moves the position pointer to next line
func (lexer *Lexer) nextLine() {
    lexer.position.Col = 0
    lexer.position.Row++
}

// nextLine increments the column counter
func (lexer *Lexer) nextCol() {
    lexer.position.Col++
}

// next moves the lexer position to next char
func (lexer *Lexer) next() {
    lexer.current = lexer.source.ReadChar()
    if lexer.current == '\n' {
        lexer.nextLine()
    } else {
        lexer.nextCol()
    }
}

// skips white space characters and returns first non whitespace character
func (lexer *Lexer) skipWSCharacters() rune {
    char := lexer.current
    if unicode.IsSpace(char) {
        lexer.next()
        return lexer.skipWSCharacters()
    }
    return char
}

// skips comments from the source code
func (lexer *Lexer) skipComments(isBlock bool) {
    lexer.next()
    char := lexer.current
    if char == -1 {
        return
    }
    if isBlock && char == '*' {
        lexer.next()
        char := lexer.current
        if char == '/' {
            lexer.next()
            return
        }
    }

    if !isBlock && char == '\n' {
        lexer.next()
        return
    }

    lexer.skipComments(isBlock)
}

func (lexer *Lexer) parseNumber(beg Position) Token {
    var t TokenType
    builder := strings.Builder{}
    builder.WriteRune(lexer.current)

    t = INT
    for {
        lexer.next()
        if unicode.IsDigit(lexer.current) {
            builder.WriteRune(lexer.current)
            continue
        }

        if lexer.current == '.' {
            t = FLOAT
            builder.WriteRune(lexer.current)
            continue
        }

        break
    }

    return Token{value: builder.String(), _type: t, beg: beg, end: lexer.position}
}

func (lexer *Lexer) parseIdentOrKeyword(beg Position) Token {
    builder := strings.Builder{}
    builder.WriteRune(lexer.current)

    for {
        lexer.next()
        if unicode.IsLetter(lexer.current) || unicode.IsDigit(lexer.current) {
            builder.WriteRune(lexer.current)
            continue
        }

        break
    }

    val := builder.String()
    t := Lookup(val)

    return Token{value: val, _type: t, beg: beg, end: lexer.position}
}

func (lexer *Lexer) illegal(beg Position) Token {
    return Token{_type: ILLEGAL, beg: beg, end: lexer.position}
}

func (lexer *Lexer) parseString(beg Position, t rune) Token {
    builder := strings.Builder{}

    for {
        lexer.next()
        if lexer.current == t {
            break
        }

        if lexer.current == -1 {
            return lexer.illegal(beg)
        }

        builder.WriteRune(lexer.current)
    }

    // skip " or '
    lexer.next()

    return Token{value: builder.String(), _type: STRING, beg: beg, end: lexer.position}
}

func (lexer *Lexer) Next() Token {
    beg := lexer.position

    lexer.skipWSCharacters()

    if unicode.IsDigit(lexer.current) {
        // parse the number
        return lexer.parseNumber(beg)
    }

    if unicode.IsLetter(lexer.current) {
        return lexer.parseIdentOrKeyword(beg)
    }

    if lexer.current == '\'' || lexer.current == '"' {
        return lexer.parseString(beg, lexer.current)
    }

    if lexer.current == '/' {
        lexer.next()
        if lexer.current == '*' {
            lexer.skipComments(true)
        } else if lexer.current == '/' {
            lexer.skipComments(false)
        } else {
            return Token{value: "/", _type: QUO, beg: beg, end: lexer.position}
        }
        return lexer.Next()
    }

    op := lexer.current
    var token Token
    switch op {
    case '+': token = Token{ value: "+", _type: ADD, beg: beg, end:lexer.position}
    case '-': {
        lexer.next()
        if lexer.current == '>' {
            token = Token{value: "->", _type: ARROW, beg:beg, end: lexer.position}
        } else {
            return Token{value: "-", _type: SUB, beg: beg , end: lexer.position}
        }
    }
    case '*': token = Token{value: "*", _type: MUL, beg: beg, end: lexer.position}
    case '%': token = Token{value: "%", _type: REM, beg: beg, end: lexer.position}
    case '^': token = Token{value: "^", _type: XOR, beg: beg, end: lexer.position}
    case '(': token = Token{value: "(", _type: LPAREN, beg: beg, end: lexer.position}
    case ')': token = Token{value: ")", _type: RPAREN, beg: beg, end: lexer.position}
    case '[': token = Token{ value: "[", _type: LBRACK, beg: beg, end: lexer.position}
    case ']': token = Token{ value: "]", _type: RBRACK, beg: beg, end: lexer.position}
    case '{': token = Token{value: "{", _type: LBRACE, beg: beg, end: lexer.position}
    case '}': token = Token{value: "}", _type: RBRACE, beg: beg, end: lexer.position}
    case '.': token = Token{value: ".", _type: PERIOD, beg: beg, end: lexer.position}
    case ',': token = Token{value: ",", _type: COMMA, beg: beg, end: lexer.position}
    case '!': {
        lexer.next()
        if lexer.current == '=' {
            token = Token{value: "!=", _type: NEQ, beg: beg, end: lexer.position}
        } else {
            return Token{value: "!", _type: NOT, beg: beg, end: lexer.position}
        }
    }

    case ';': token = Token{value: ";", _type: SEMICOLON, beg: beg, end: lexer.position}
    case ':': token = Token{value: ":", _type: COLON, beg: beg, end: lexer.position}

    case '=': {
        lexer.next()
        if lexer.current == '=' {
            token = Token{value: "==", _type: EQL, beg: beg, end: lexer.position}
        } else {
            // we return immediately as we do not need to call next now
            return Token{value: "=", _type: ASSIGN, beg: beg, end: lexer.position}
        }
    }

    case '&': {
        lexer.next()
        if lexer.current == '&' {
            token = Token{value: "&&", _type: LAND, beg: beg, end: lexer.position}
        } else {
            return Token{value: "&", _type: AND, beg: beg, end: lexer.position}
        }
    }

    case '|': {
        lexer.next()
        if lexer.current == '|' {
            token = Token{value: "||", _type: LOR, beg: beg, end: lexer.position}
        } else {
            return Token{value: "|", _type: OR, beg: beg, end: lexer.position}
        }
    }

    case '<': {
        lexer.next()
        if lexer.current == '=' {
            token = Token{value: "<=", _type: LEQ, beg: beg, end: lexer.position}
        } else {
            return Token{value: "<", _type: LSS, beg: beg, end: lexer.position}
        }
    }

    case '>': {
        lexer.next()
        if lexer.current == '=' {
            token = Token{value: ">=", _type: GEQ, beg: beg, end: lexer.position}
        } else {
            return Token{value: ">", _type: GTR, beg: beg, end: lexer.position}
        }
    }

    default: return lexer.illegal(beg)
    }

    lexer.next()
    return token
}
