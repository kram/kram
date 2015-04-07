// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"github.com/zegl/Gus/src/types"
)

type Library struct{}

func (lib *Library) Init(str string)               {}
func (lib *Library) Instance() (types.Lib, string) { return &Library{}, lib.Type() }
func (lib *Library) Type() string                  { return "Library" }
func (lib *Library) ToString() string              { return lib.Type() }

func (lib *Library) TypeWithLib(l types.Lib) *types.Type {
	class := types.Type{}
	class.InitWithLib(l)

	return &class
}
