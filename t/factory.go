package t

type TypeFactory struct {
    intT *intType
    floatT *floatType
    stringT *stringType
    pointerTypes map[Type]*pointerType
    structTypes map[string]*structType
    arrayTypes map[Type]*arrayType
}

func (f *TypeFactory) IntType() *intType {
    return f.intT
}

func (f *TypeFactory) FloatType() *floatType {
    return f.floatT
}

func (f *TypeFactory) StringType() *stringType {
    return f.stringT
}

func (f *TypeFactory) PointerType(t Type) *pointerType {
    if pt, ok := f.pointerTypes[t]; ok {
        return pt
    }

    pt := &pointerType{elemType: t}
    f.pointerTypes[t] = pt
    return pt
}

func (f *TypeFactory) ArrayType(t Type) *arrayType {
    if at, ok := f.arrayTypes[t]; ok {
        return at
    }

    // TODO(@me): support for array types with known size
    at := &arrayType{elemType:t, count: -1}
    f.arrayTypes[t] = at
    return at
}

func (f *TypeFactory) StructType(elemTypes map[string]Type) {
    if st, ok := f.structTypes[elemTypes]
}

