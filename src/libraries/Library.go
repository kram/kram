// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"github.com/zegl/Gus/src/types"
	"github.com/zegl/Gus/src/types/builtin"
	"log"
)

type Library struct{}

func (self Library) Init(str string) {
	// Do nothing
}

func (self Library) InitWithParams(params []*types.Type) {
	log.Panic("This Type does not support InitWithParams()")
}

func (self Library) ToString() string {
	return "<nil>"
}

func (self Library) Bool(value bool) *types.Type {

	bl := builtin.Bool{}
	bl.Set(value)

	return self.fromLib(&bl)
}

func (self Library) Null() *types.Type {
	return self.fromLib(&builtin.Null{})
}

func (self Library) String(str string) *types.Type {
	return self.fromLib(&builtin.String{
		Value: str,
	})
}

func (self Library) fromLib(lib types.Lib) *types.Type {
	class := types.Type{}
	class.InitWithLib(lib)

	return &class
}