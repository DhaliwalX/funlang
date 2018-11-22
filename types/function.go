package types

import (
    "fmt"
    "strings"
)

type functionType struct {
    returnType Type
    argsType []Type
}

func (t *functionType) Elem() Type {
    return nil
}

func (t *functionType) Tag() TypeTag {
    return FUNCTION_TYPE
}

func (t *functionType) Name() string {
    builder := strings.Builder{}
    builder.WriteString("func(")
    for _, argType := range t.argsType {
        builder.WriteString(argType.Name())
    }

    builder.WriteString(fmt.Sprintf("): %s", t.returnType.Name()))
    return builder.String()
}

func (t *functionType) Field(string) Type {
    return nil
}

func (t *functionType) ReturnType() Type {
    return t.returnType
}

func (t *functionType) LenArgs() int {
    return len(t.argsType)
}

func (t *functionType) ArgByIndex(i int) Type {
    if i >= len(t.argsType) {
        return nil
    }

    return t.argsType[i]
}

func ToFunctionType(t Type) *functionType {
    return t.(*functionType)
}
