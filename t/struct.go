package t

import "strings"

// represents structs
type structType struct {
    elems map[string]Type
}

func (s *structType) Elem() Type {
    return nil
}

func (s *structType) Field(name string) Type {
    t, ok := s.elems[name]
    if !ok {
        return nil
    }

    return t
}

func (s *structType) Tag() TypeTag {
    return STRUCT_TYPE
}

func (s *structType) Name() string {
    builder := strings.Builder{}
    builder.WriteString("struct{")
    for name, t := range s.elems {
        builder.WriteString(name)
        builder.WriteString(":")
        builder.WriteString(t.Name())
        builder.WriteString(";")
    }

    return builder.String()
}

func (s *structType) LenFields() int {
    return len(s.elems)
}

func ToStructType(t Type) *structType {
    if t.Tag() != STRUCT_TYPE {
        return nil
    }

    return t.(*structType)
}
