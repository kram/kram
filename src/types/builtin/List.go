package builtin

import (
	"strings"
	"log"
	"math"
	"../" // types
)

type List struct {
	*Builtin
	Items []*types.Type
}

func (list *List) Instance() (types.Lib, string) {
	return &List{}, list.Type()
}

func (list *List) Type() string {
	return "List"
}

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
	list.Items = make([]*types.Type, 0)
}

func (list *List) InitWithParams(params []*types.Type) {
	list.Items = make([]*types.Type, 0)
	list.Push(params)
}

func (list *List) Push(params []*types.Type) {
	for _, param := range params {
		list.Items = append(list.Items, param)
	}
}

func (list *List) Pop(params []*types.Type) *types.Type {
	res := list.Items[len(list.Items)-1]
	list.Items = list.Items[:len(list.Items)-1]

	return res
}

func (list *List) ItemAt(params []*types.Type) *types.Type {
	if len(params) != 1 {
		log.Panic("List::ItemAt() expected only 1 parameter")
	}

	param := params[0]

	if num, ok := param.Extension.(*Number); ok {
		return list.ItemAtNumber(num)
	}

	if li, ok := param.Extension.(*List); ok {
		return list.ItemAtList(li)
	}


	log.Panic("List::ItemAt() expected parameter 1 to be of type Number or List")

	// Will never be reached
    return list.Null()
}

func (list *List) ItemAtNumber(num *Number) *types.Type {
	// Use https://golang.org/pkg/math/#Trunc to make sure that the float
	// is a whole number
	key_float := math.Trunc(num.Value)

	if key_float != num.Value {
		log.Panic("List::ItemAt() can only be used together with whole numbers")
	}

	if len(list.Items) > int(key_float) {
		return list.Items[int(key_float)]
	}

	log.Panic("List::ItemAt() out of range!")

	// Will never be reached
    return list.Null()
}

func (list *List) ItemAtList(li *List) *types.Type {

	res := List{}
	res.Items = make([]*types.Type, 0)

	for _, item := range li.Items {
		if num, ok := item.Extension.(*Number); ok {
			res.Items = append(res.Items, list.ItemAtNumber(num))
		}
	}

	class := types.Type{}
	class.Init("List")
	class.Extension = &res

	return &class
}

// Used when iterating over each object in the list
func (list *List) Length() int {
	return len(list.Items)
}

// Used when iterating over each object in the list
func (list *List) ItemAtPosition(pos int) *types.Type {
	return list.Items[pos]
}
