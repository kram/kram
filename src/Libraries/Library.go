package libraries

import (
	"../Instructions"
	"../Types"
)

type ON int

const (
	ON_NOTHING ON = 1 << iota // 1
	ON_CLASS                  // 2
)

type VirtualMachine struct {
	Operation func(instructions.Node, ON) types.Type
}

func false() types.Type {
	bl := types.Bool{}
	bl.Init("false")

	return &bl
}
