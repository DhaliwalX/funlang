package types

import "fmt"

type pointerType struct {
    elemType Type
}

func (t *pointerType) Elem() Type {
    return t.elemType
}

func (t *pointerType) Field(string) Type {
    return nil
}

func (t *pointerType) Tag() TypeTag {
    return POINTER_TYPE
}

func (t *pointerType) Name() string {
    return fmt.Sprintf("{%s}", t.Elem().Name())
}

func (t *pointerType) String() string {
    return t.Name()
}

func ToPointerType(t Type) *pointerType {
    return t.(*pointerType)
}
