// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/zegl/Gus/src/types"
)

type Builtin struct {}

func (self Builtin) Null() *types.Type {
	return self.fromLib(&Null{})
}

func (self Builtin) String(str string) *types.Type {
	return self.fromLib(&String{
		Value: str,
	})
}

func (self Builtin) fromLib(lib types.Lib) *types.Type {
	class := types.Type{}
	class.InitWithLib(lib)

	return &class
}