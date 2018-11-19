package types

import "fmt"

// placeHolder represents a type which has not been resolved yet
type placeHolder struct {
	resolvedType Type
}

func (p *placeHolder) Name() string {
	if p.resolvedType != nil {
		return fmt.Sprintf("resolved{%s}", p.resolvedType.Name())
	}

	return "unresolved"
}

func (p *placeHolder) Elem() Type {
	return p.resolvedType
}

func (p *placeHolder) Field(string) Type {
	return nil
}

func (p *placeHolder) Tag() TypeTag {
	return UNKNOWN_TYPE
}

// Resolve resolves an unknown type to resolved type
// in case of placeHolder
func Resolve(t Type) Type {
	if t.Tag() == UNKNOWN_TYPE {
		return t.(*placeHolder).resolvedType
	}

	// in other cases just return t
	return t
}
