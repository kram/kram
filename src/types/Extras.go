// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package types

import (
	"github.com/zegl/Gus/src/instructions"
)

type ON int

const (
	ON_NOTHING     ON = 1 << iota // 1
	ON_CLASS                      // 2
	ON_METHOD_BODY                // 4
	ON_FOR_PART                   // 8
)

type VM interface {
	EnvironmentPush()
	EnvironmentPop()
	Operation(instructions.Node, ON) *Type
	OperationAssign(instructions.Assign) *Type
	OperationBlock(instructions.Block, ON) *Type
}
