package parse

import (
	"bitbucket.org/dhaliwalprince/funlang/lex"
	"fmt"
	"strings"
)

type parseError struct {
	message string
	token lex.Token
}

func newParseError(token lex.Token, message string) parseError {
	return parseError{token:token, message:message}
}

func (e parseError) Error() string {
	return fmt.Sprintf("%d:%d:%s (%s)", e.token.Begin(), e.token.End(), e.message, e.token.Value())
}

func unexpectedToken(t lex.Token, expected lex.TokenType) error {
	return fmt.Errorf("ParseError:%d:%d: unexpected token %s (expected %s)", t.Begin().Col, t.End().Row, t.Type(), expected)
}

type errorList struct {
	list []error
}

func (list *errorList) append(err error) {
	list.list = append(list.list, err)
}

func (list *errorList) Error() string {
	builder := strings.Builder{}
	for _, err := range list.list {
		builder.WriteString(err.Error()+"\n")
	}

	return builder.String()
}
