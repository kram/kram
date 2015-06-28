// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/kram/kram/src/types"
	"log"
)

type Builtin struct{}

func (self Builtin) InitWithParams(params []*types.Class) {
	log.Panic("This Type does not support InitWithParams()")
}

func (self Builtin) Null() *types.Class {
	return self.fromLib(&Null{})
}

func (self Builtin) String(str string) *types.Class {
	return self.fromLib(&String{
		Value: str,
	})
}

func (self Builtin) Bool(value bool) *types.Class {
	b := Bool{}
	b.Set(value)

	return self.fromLib(&b)
}

func (self Builtin) fromLib(lib types.Lib) *types.Class {
	class := types.Class{}
	class.InitWithLib(lib)

	return &class
}
