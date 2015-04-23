// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package types

type Literal struct {}

func (lit Literal) IsClass() bool {
	return false
}

type LiteralString struct {
	Literal
	String string
}

type LiteralNumber struct {
	Literal
	Number float64
}

type LiteralBool struct {
	Literal
	Bool bool
}

type LiteralNull struct {
	Literal
	Null bool
}