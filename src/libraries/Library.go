// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"github.com/kram/kram/src/types"
	"github.com/kram/kram/src/types/builtin"
	"log"
)

type Library struct{}

func (self Library) Init(str string) {
	// Do nothing
}

func (self Library) InitWithParams(params []*types.Class) {
	log.Panic("This Type does not support InitWithParams()")
}

func (self Library) ToString() string {
	return "<nil>"
}

func (self Library) Params(in []*types.Class, out ...interface{}) {
	for k, v := range in {
		switch o := out[k].(type) {
		case *string:
			*o = v.ToString()
		default:
			log.Panic("Library.Params() can not handle this type?")
		}
	}
}

func (self Library) Bool(value bool) *types.Class {
	bl := builtin.Bool{}
	bl.Set(value)

	return self.fromLib(&bl)
}

func (self Library) Null() *types.Class {
	return self.fromLib(&builtin.Null{})
}

func (self Library) String(str string) *types.Class {
	return self.fromLib(&builtin.String{
		Value: str,
	})
}

func (self Library) Number(nr float64) *types.Class {
	return self.fromLib(&builtin.Number{
		Value: nr,
	})
}

func (self Library) InitMap() *builtin.Map {
	m := builtin.Map{}
	m.Init("")

	return &m
}

func (self Library) InitList() *builtin.List {
	l := builtin.List{}
	l.Init("")

	return &l
}

func (self Library) fromLib(lib types.Lib) *types.Class {
	class := types.Class{}
	class.InitWithLib(lib)

	return &class
}

func (self Library) FromGolangValues(in interface{}) *types.Class {
	return self.Null()
}
