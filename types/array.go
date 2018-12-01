package types

import "fmt"

// arrayType represents arrays
type arrayType struct {

    // array is of this type
    elemType Type
    // number of elements in this array, -1 for unknown size
    count int
}

func (t *arrayType) Elem() Type {
    return t.elemType
}

func (t *arrayType) Field(string) Type {
    return nil
}

func (t *arrayType) Tag() TypeTag {
    return ARRAY_TYPE
}

func (t *arrayType) Name() string {
    if t.count > 0 {
        return fmt.Sprintf("%s[%d]", t.elemType.Name(), t.count)
    } else {
        return fmt.Sprintf("%s[]", t.elemType.Name())
    }
}

func (t *arrayType) String() string {
    return t.Name()
}

func ToArrayType(t Type) *arrayType {
    return t.(*arrayType)
}
