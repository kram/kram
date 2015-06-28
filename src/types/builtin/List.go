// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package builtin

import (
	"github.com/kram/kram/src/types"
	"log"
	"math"
	"strings"
)

type List struct {
	Builtin
	Items []*types.Class
}

func (self List) Instance() (types.Lib, string) { return &List{}, self.Type() }
func (self List) Type() string                  { return "List" }
func (self List) M_Type() *types.Class          { return self.String(self.Type()) }

func (list *List) ToString() string {
	out := make([]string, len(list.Items))

	for i, item := range list.Items {
		out[i] = item.ToString()
	}

	return "[" + strings.Join(out, ", ") + "]"
}

// List can not be initialized with Init
// see InitWithParams
func (list *List) Init(str string) {
	list.Items = make([]*types.Class, 0)
}

func (list *List) InitWithParams(params []*types.Class) {
	list.Items = make([]*types.Class, 0)
	list.M_Push(params)
}

func (list *List) Push(data *types.Class) {
	params := make([]*types.Class, 1)
	params[0] = data
	list.M_Push(params)
}

func (list *List) M_Push(params []*types.Class) {
	for _, param := range params {
		list.Items = append(list.Items, param)
	}
}

func (list *List) M_Pop(params []*types.Class) *types.Class {
	res := list.Items[len(list.Items)-1]
	list.Items = list.Items[:len(list.Items)-1]

	return res
}

// Adressable from VM
func (list *List) ItemAt(params []*types.Class) *types.Class {
	return list.M_Get(params)
}

func (list *List) M_Get(params []*types.Class) *types.Class {
	if len(params) != 1 {
		log.Panic("List::Get() expected only 1 parameter")
	}

	param := params[0]

	if num, ok := param.Extension.(*Number); ok {
		return list.ItemAtNumber(num)
	}

	if li, ok := param.Extension.(*List); ok {
		return list.ItemAtList(li)
	}

	log.Panic("List::Get() expected parameter 1 to be of type Number or List")

	// Will never be reached
	return list.Null()
}

func (list *List) ItemAtNumber(num *Number) *types.Class {
	// Use https://golang.org/pkg/math/#Trunc to make sure that the float
	// is a whole number
	key_float := math.Trunc(num.Value)

	if key_float != num.Value {
		log.Panic("List::Get() can only be used together with whole numbers")
	}

	if len(list.Items) > int(key_float) {
		return list.Items[int(key_float)]
	}

	log.Panic("List::Get() out of range!")

	// Will never be reached
	return list.Null()
}

func (list *List) ItemAtList(li *List) *types.Class {

	res := List{}
	res.Items = make([]*types.Class, 0)

	for _, item := range li.Items {
		if num, ok := item.Extension.(*Number); ok {
			res.Items = append(res.Items, list.ItemAtNumber(num))
		}
	}

	class := types.Class{}
	class.Init("List")
	class.Extension = &res

	return &class
}

// Used when iterating over each object in the list
func (list *List) Length() int {
	return len(list.Items)
}

// Used when iterating over each object in the list
func (list *List) ItemAtPosition(pos int) *types.Class {
	return list.Items[pos]
}
