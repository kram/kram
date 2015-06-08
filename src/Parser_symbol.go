// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package gus

import (
	ins "github.com/zegl/Gus/src/instructions"
)

type Symbol struct {
	Function   SymbolFunction
	Importance int
}

type SymbolFunction func(ON) ins.Node
