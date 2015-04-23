// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/zegl/Gus/src/types"
)

type Bool struct {
	Builtin
	Value bool
}

func (self Bool) Instance() (types.Lib, string) { return &Bool{}, self.Type() }
func (self Bool) Type() string { return "Bool" }
func (self Bool) M_Type() *types.Class { return self.String(self.Type()) }

func (self *Bool) Init(str string) {
	if str == "true" {
		self.Value = true
	} else {
		self.Value = false
	}
}

func (self *Bool) ToString() string {
	if self.Value {
		return "true"
	}

	return "false"
}

func (self *Bool) Set(bl bool) {
	self.Value = bl
}

func (self *Bool) IsTrue() bool {
	return self.Value
}

/*

func (self *Bool) Compare(method string, right Type) Type {

	r, ok := right.(*Bool)

	if !ok {
		log.Panicf("You can not compare a %s() with a %s()", self.Type(), right.Type())
	}

	if self.Value == r.Value {
		bl := Bool{}
		bl.Init("true")

		return &bl
	}

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// Will never be reached

	bl := Bool{}
	bl.Init("false")

	return &bl
}
*/
