// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/zegl/Gus/src/types"
	"log"
	"strings"
)

type String struct {
	Builtin
	Value string
}

func (self String) Instance() (types.Lib, string) { return &String{}, self.Type() }
func (self String) Type() string { return "String" }
func (self String) M_Type() *types.Class { return self.String(self.Type()) }

func (self *String) Init(str string) {
	self.Value = str
}

func (self *String) InitWithParams(params []*types.Class) {
	i := params[0]

	if i.Type() == "Number" {
		self.Init(i.ToString())
	} else {
		log.Panicf("Can not init String with %s", i.Type())
	}
}

func (self *String) ToString() string {
	return self.Value
}

func (self *String) Math(method string, right *types.Class) *types.Class {

	r, ok := right.Extension.(*String)

	if !ok {
		log.Panicf("You can not apply %s to a %s() with a %s()", method, self.Type(), right.Type())
	}

	// String concatenation
	if method == "+" {
		str := String{}
		str.Init(self.Value + r.Value)

		res := types.Class{}
		res.InitWithLib(&str)

		return &res
	}

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// This code will never be reached

	return &types.Class{}
}

func (self *String) Compare(method string, right *types.Class) *types.Class {

	r, ok := right.Extension.(*String)

	if !ok {
		log.Panicf("You can not compare a %s() with a %s()", self.Type(), right.Type())
	}

	b := false

	switch method {
	case ">":
		b = self.Value > r.Value
	case "<":
		b = self.Value < r.Value
	case ">=":
		b = self.Value >= r.Value
	case "<=":
		b = self.Value <= r.Value
	case "==":
		b = self.Value == r.Value
	case "!=":
		b = self.Value != r.Value
	default:
		log.Panicf("%s() is not implementing %s", self.Type(), method)
	}

	bl := Bool{}
	bl.Set(b)

	res := types.Class{}
	res.InitWithLib(&bl)

	return &res
}
	
// Constructor method
func (self *String) M_New(params []*types.Class) {
	self.InitWithParams(params)
}

func (self *String) M_Lower(params []*types.Class) *types.Class {
	return self.String(strings.ToLower(self.Value))
}

func (self *String) M_Upper(params []*types.Class) *types.Class {
	return self.String(strings.ToUpper(self.Value))
}