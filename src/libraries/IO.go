// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package libraries

import (
	"fmt"
	"github.com/zegl/Gus/src/types"
)

type Library_IO struct {
	*Library
}

func (io *Library_IO) Instance() (types.Lib, string) {
	return &Library_IO{}, "IO"
}

func (io Library_IO) Print(params []*types.Type) {
	for _, param := range params {
		fmt.Print(param.ToString())
	}
}

func (io Library_IO) Println(params []*types.Type) {
	for _, param := range params {
		fmt.Println(param.ToString())
	}
}
