// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package gus

// ON is the datatype used by the parser to differentiate between different parts of the sourcecode
// Eg. the "static"-keyword is only used when ON_CLASS_BODY is active. And func() can be both a Call to a method
// and the start of a method, the difference is simply if ON is ON_PUSH_CLASS or ON_CLASS_BODY.
type ON int

const (
	ON_DEFAULT           ON = 1 << iota // 1
	ON_CLASS_BODY                       // 2
	ON_PUSH_CLASS                       // 4
	ON_METHOD_PARAMETERS                // 8
	ON_ARGUMENTS                        // 16
)
