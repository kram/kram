// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package kram

import (
	ins "github.com/kram/kram/src/instructions"
)

type Symbol struct {
	Function   SymbolFunction
	Importance int
}

type SymbolFunction func(ON) ins.Node
