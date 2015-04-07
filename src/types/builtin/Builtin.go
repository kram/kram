// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/zegl/Gus/src/types"
)

type Builtin struct{}

func (buil Builtin) Null() *types.Type {
	class := types.Type{}
	class.InitWithLib(&Null{})

	return &class
}
