package types

import "bitbucket.org/dhaliwalprince/funlang/context"

type Factory struct {
    ctx *context.Context
    intT *intType
    floatT *floatType
    stringT *stringType
    pointerTypes map[Type]*pointerType
    structTypes map[string]*structType
    functionTypes map[string]*functionType
    arrayTypes map[Type]*arrayType
}

var defaultFactory *Factory

func NewFactory(ctx *context.Context) *Factory {
    if defaultFactory == nil {
        defaultFactory = &Factory{
            ctx: ctx,
            intT: &intType{},
            floatT: &floatType{},
            stringT: &stringType{},
            pointerTypes: make(map[Type]*pointerType),
            structTypes: make(map[string]*structType),
            functionTypes: make(map[string]*functionType),
            arrayTypes: make(map[Type]*arrayType),
        }
    }

    return defaultFactory
}

func (f *Factory) IntType() *intType {
    return f.intT
}

func (f *Factory) FloatType() *floatType {
    return f.floatT
}

func (f *Factory) StringType() *stringType {
    return f.stringT
}

func (f *Factory) PointerType(t Type) *pointerType {
    if pt, ok := f.pointerTypes[t]; ok {
        return pt
    }

    pt := &pointerType{elemType: t}
    f.pointerTypes[t] = pt
    return pt
}

func (f *Factory) ArrayType(t Type) *arrayType {
    if at, ok := f.arrayTypes[t]; ok {
        return at
    }

    // TODO(@me): support for array types with known size
    at := &arrayType{elemType:t, count: -1}
    f.arrayTypes[t] = at
    return at
}

// TODO(@me): need to find a better way
func (f *Factory) StructType(elemTypes map[string]Type) *structType {
    if st, ok := f.structTypes[(&structType{ elems: elemTypes}).Name()]; ok {
        return st
    }
    st := &structType{elems: elemTypes}
    f.structTypes[st.Name()] = st
    return st
}
// TODO(@me,sameAsPrevious)
func (f *Factory) FunctionType(retType Type, argsTypes []Type) *functionType {
    if st, ok := f.functionTypes[(&functionType{returnType:retType, argsType:argsTypes}).Name()]; ok {
        return st
    }

    st := &functionType{returnType: retType, argsType: argsTypes}
    f.functionTypes[st.Name()] = st
    return st
}

func (f *Factory) UnknownType() *placeHolder {
    return &placeHolder{}
}
