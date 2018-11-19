package types

import "bitbucket.org/dhaliwalprince/funlang/context"

type typeFactory struct {
    ctx *context.Context
    intT *intType
    floatT *floatType
    stringT *stringType
    pointerTypes map[Type]*pointerType
    structTypes map[string]*structType
    functionTypes map[string]*functionType
    arrayTypes map[Type]*arrayType

    // this stores all the types with their names
    types map[string]Type
}

var factories map[*context.Context]*typeFactory

func Factory(ctx *context.Context) *typeFactory {
    if f, ok := factories[ctx]; ok {
        return f
    }

    f := &typeFactory{
        ctx: ctx,
        intT: &intType{},
        floatT: &floatType{},
        stringT: &stringType{},
        pointerTypes: make(map[Type]*pointerType),
        structTypes: make(map[string]*structType),
        functionTypes: make(map[string]*functionType),
        arrayTypes: make(map[Type]*arrayType),
        types: make(map[string]Type),
    }

    f.installBaseTypes()
    factories[ctx] = f
    return f
}

func (f *typeFactory) installBaseTypes() {
    f.types["int"] = f.intT
    f.types["float"] = f.floatT
    f.types["string"] = f.stringT
}

func (f *typeFactory) IntType() *intType {
    return f.intT
}

func (f *typeFactory) FloatType() *floatType {
    return f.floatT
}

func (f *typeFactory) StringType() *stringType {
    return f.stringT
}

func (f *typeFactory) PointerType(t Type) *pointerType {
    if pt, ok := f.pointerTypes[t]; ok {
        return pt
    }

    pt := &pointerType{elemType: t}
    f.pointerTypes[t] = pt
    return pt
}

func (f *typeFactory) ArrayType(t Type) *arrayType {
    if at, ok := f.arrayTypes[t]; ok {
        return at
    }

    // TODO(@me): support for array types with known size
    at := &arrayType{elemType:t, count: -1}
    f.arrayTypes[t] = at
    return at
}

// TODO(@me): need to find a better way
func (f *typeFactory) StructType(elemTypes map[string]Type) *structType {
    if st, ok := f.structTypes[structType{ elems: elemTypes}.Name()]; ok {
        return st
    }
    st := &structType{elems: elemTypes}
    f.structTypes[st.Name()] = st
    return st
}
// TODO(@me,sameAsPrevious)
func (f *typeFactory) FunctionType(retType Type, argsTypes []Type) *functionType {
    if st, ok := f.functionTypes[functionType{returnType:retType, argsType:argsTypes}.Name()]; ok {
        return st
    }

    st := &functionType{returnType: retType, argsType: argsTypes}
    f.functionTypes[st.Name()] = st
    return st
}

// Named returns the type referenced by a name known to source code
func (f *typeFactory) Named(name string) Type {
    if t, ok := f.types[name]; ok {
        return t
    }

    return nil
}

// AddTypename assigns name to a type
func (f *typeFactory) AddTypename(name string, t Type) {
    f.types[name] = t
}
