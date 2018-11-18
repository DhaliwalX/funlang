// base.go contains the definition for base type interface
package t

type TypeTag int

const (
    VOID_TYPE TypeTag = iota
    INT_TYPE  // int
    FLOAT_TYPE // float
    STRING_TYPE // string
    FUNCTION_TYPE // func
    POINTER_TYPE // int*
    STRUCT_TYPE // struct
    ARRAY_TYPE // int[]
)

type Type interface {
    Name() string

    // returns this type's element type
    Elem() Type

    // returns the named field type
    Field(name string) Type

    // Tag returns the type tag
    Tag() TypeTag
}
