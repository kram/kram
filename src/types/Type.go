// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package types

type Value struct {
	Type      Value_Type
	Number    float64
	String    string
	Reference *Class
}

type Value_Type uint8

const (
	NULL   Value_Type = 1 << iota // No value needed
	BOOL                          // Stored as 1 or 0 in value_number
	NUMBER                        // Stored in value_number
	STRING                        // Stored in value_string
	CLASS                         // A reference / pointer? to the class in value_class
)
