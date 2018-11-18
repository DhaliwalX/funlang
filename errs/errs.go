package errs

import "github.com/pkg/errors"

func NewEOFError() error {
	return errors.New("eof")
}

func NewSyntaxError() error {
	return errors.New("syntax error")
}
