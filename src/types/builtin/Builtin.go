package builtin

import (
	"../"
)

type Builtin struct {}

func (buil Builtin) Null() *types.Type {
	class := types.Type{}
	class.InitWithLib(&Null{})

	return &class
}