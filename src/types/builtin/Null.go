// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/zegl/Gus/src/types"
)

type Null struct {
	Builtin
}

func (self Null) Instance() (types.Lib, string) { return &Null{}, self.Type() }
func (self Null) Type() string                  { return "Null" }
func (self Null) M_Type() *types.Class          { return self.String(self.Type()) }

func (self Null) Init(str string)   {}
func (self *Null) ToString() string { return "null" }
