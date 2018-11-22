package lex

import "testing"

func sourceFromString(source string) Source {
    return &FileSource{ source: source }
}

func TestLexer_Next(t *testing.T) {
    source := "notakeyword"
    lex := NewLexer(sourceFromString(source))

    k := lex.Next()
    if k.Value() != source {
        t.Error(k, "is not same as source")
    }

    if k.Type() != IDENT {
        t.Error(k, "is not ident")
    }
}

func TestLexer_Next2(t *testing.T) {
    source := "1230"
    lex := NewLexer(sourceFromString(source))

    k := lex.Next()
    if k.Value() != source {
        t.Error(k, "is not same as source")
    }

    if k.Type() != INT {
        t.Error(k, "is not an int")
    }
}

func TestLexer_Next3(t *testing.T) {
    source := "/* comment */ key"
    lex := NewLexer(sourceFromString(source))

    k := lex.Next()
    if k.Value() != "key" {
        t.Error(k, "is not same as key")
    }

    if k.Type() != IDENT {
        t.Error(k, "is not an ident")
    }

}

func TestLexer_Next4(t *testing.T) {
    source := `
    // some comment
    val // another comment
    x
    `
    lex := NewLexer(sourceFromString(source))
    k := lex.Next()
    if k.Type() != IDENT {
        t.Error(k, "is not an ident")
    }

    k = lex.Next()
    if k.Type() != IDENT {
        t.Error(k, "is not an ident")
    }
}

func TestLexerString(t *testing.T) {
    source := `
        "somestring"
`
    lex := NewLexer(sourceFromString(source))
    k := lex.Next()
    if k.Type() != STRING {
        t.Error(k, "is not a string")
    }

    if k.Value() != "somestring" {
        t.Error(k, "is not somestring")
    }
}

func TestLexerParseOperator(t *testing.T) {
    for i := operator_beg+1; i != operator_end; i++ {
        source := tokens[i]
        lex := NewLexer(sourceFromString(source))
        k :=  lex.Next()
        if k.Type() != i {
            t.Error(k, "is not ", tokens[i])
        }
        if k.Value() != tokens[i] {
            t.Error(k, "is not", tokens[i])
        }
    }
}

func TestLexerClose(t *testing.T) {
    source := "struct{}"
    lex := NewLexer(sourceFromString(source))
    tok := lex.Next()
    if tok.Type() != STRUCT {
        t.Error(tok, "is not", STRUCT)
    }

    tok = lex.Next()
    if tok.Type() != LBRACE {
        t.Error(tok, "is not", LBRACE)
    }

    tok = lex.Next()
    if tok.Type() != RBRACE {
        t.Error(tok, "is not", RBRACE)
    }
}
