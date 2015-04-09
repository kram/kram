// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"github.com/zegl/Gus/src/types"
	"github.com/zegl/Gus/src/types/builtin"
	"strings"
)

type Library_String struct {
	Library
}

func (self Library_String) Instance() (types.Lib, string) { return &Library_String{}, self.Type() }
func (self Library_String) Type() string { return "Strings" }
func (self Library_String) M_Type() *types.Type { return self.String(self.Type()) }

func (self Library_String) M_ToLower(params []*types.Type) *types.Type {
	str := builtin.String{}

	for _, param := range params {
		str.Init(strings.ToLower(param.ToString()))
		break
	}

	return self.fromLib(&str)
}

func (self Library_String) M_ToUpper(params []*types.Type) *types.Type {
	str := builtin.String{}

	for _, param := range params {
		str.Init(strings.ToUpper(param.ToString()))
		break
	}

	return self.fromLib(&str)
}
