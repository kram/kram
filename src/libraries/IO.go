// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"fmt"
	"github.com/kram/kram/src/types"
)

type Library_IO struct {
	Library
}

func (self Library_IO) Instance() (types.Lib, string) { return &Library_IO{}, self.Type() }
func (self Library_IO) Type() string                  { return "IO" }
func (self Library_IO) M_Type() *types.Class          { return self.String(self.Type()) }

func (io Library_IO) M_Print(params []*types.Class) {
	for _, param := range params {
		fmt.Print(param.ToString())
	}
}

func (io Library_IO) M_Println(params []*types.Class) {
	for _, param := range params {
		fmt.Println(param.ToString())
	}
}
