package ssa

import (
	"bitbucket.org/dhaliwalprince/funlang/types"
	"fmt"
	"strings"
)

type Program struct {
	Types map[string]types.Type
	Globals map[string]Value
}

func (p *Program) String() string {
	builder := strings.Builder{}
	builder.WriteString("# funlang bytecode\n\n")
	builder.WriteString("# types\n")

	for name, t := range p.Types {
		builder.WriteString(fmt.Sprintf("type %s %s\n", name, t))
	}

	builder.WriteString("\n")
	for _, global := range p.Globals {
		builder.WriteString(fmt.Sprintf("%s\n", global))
	}

	return builder.String()
}
