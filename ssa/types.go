package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/types"
)

var typeFactory *types.Factory

func init() {
	ctx := context.Context{}
	typeFactory = types.NewFactory(&ctx)
}
