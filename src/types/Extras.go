package types

import (
	"../instructions"
)

type ON int

const (
	ON_NOTHING    ON = 1 << iota // 1
	ON_CLASS                     // 2
	ON_METHOD_BODY               // 4
	ON_FOR_PART                  // 8
)

type VM interface {
	EnvironmentPush()
	EnvironmentPop()
	Operation(instructions.Node, ON) *Type
	OperationAssign(instructions.Assign) *Type
	OperationBlock(instructions.Block, ON) *Type
}