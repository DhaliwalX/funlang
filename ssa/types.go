package ssa

import (
	"funlang/context"
	"funlang/types"
)

var typeFactory *types.Factory

func init() {
	ctx := context.Context{}
	typeFactory = types.NewFactory(&ctx)
}
