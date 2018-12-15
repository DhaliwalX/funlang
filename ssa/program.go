package ssa

import "bitbucket.org/dhaliwalprince/funlang/types"

type Program struct {
	Types map[string]types.Type
	Globals map[string]Value
}
