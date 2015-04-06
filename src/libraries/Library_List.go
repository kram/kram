package main

import (
	"strings"
	"log"
	"math"
)

type Library_List struct {
	*Library
	Items []Type
}

func (list *Library_List) Instance() (Lib, string) {
	return &Library_List{}, "List"
}

func (list *Library_List) Init(params []Type) {
	list.Items = make([]Type, 0)
	list.Push(params)
}

func (list *Library_List) Push(params []Type) {
	for _, param := range params {
		list.Items = append(list.Items, param)
	}
}

func (list *Library_List) Pop(params []Type) Type {
	res := list.Items[len(list.Items)-1]
	list.Items = list.Items[:len(list.Items)-1]

	return res
}

func (list *Library_List) ItemAt(params []Type) Type {
	if len(params) != 1 {
		log.Panic("Library_List::ItemAt() expected only 1 parameter")
	}

	param := params[0]

	if num, ok := param.(*Number); ok {
		return list.ItemAtNumber(num)
	}

	if class, ok := param.(*Class); ok {
		if li, ok := class.Extension.(*Library_List); ok {
			return list.ItemAtList(li)
		}
	}

	log.Panic("Library_List::ItemAt() expected parameter 1 to be of type Number or List")

	// Will never be reached
    return &Null{}
}

func (list *Library_List) ItemAtNumber(num *Number) Type {
	// Use https://golang.org/pkg/math/#Trunc to make sure that the float
	// is a whole number
	key_float := math.Trunc(num.Value)

	if key_float != num.Value {
		log.Panic("Library_List::ItemAt() can only be used together with whole numbers")
	}

	if len(list.Items) > int(key_float) {
		return list.Items[int(key_float)]
	}

	log.Panic("Library_List::ItemAt() out of range!")

	// Will never be reached
    return &Null{}
}

func (list *Library_List) ItemAtList(li *Library_List) Type {

	res := Library_List{}
	res.Items = make([]Type, 0)

	for _, item := range li.Items {
		if num, ok := item.(*Number); ok {
			res.Items = append(res.Items, list.ItemAtNumber(num))
		}
	}

	class := Class{}
	class.Init("LIst")
	class.Extension = &res

	return &class
}

func (list *Library_List) ToString() string {
	out := make([]string, len(list.Items))

	for i, item := range list.Items {
		out[i] = item.ToString()
	}

	return "[" + strings.Join(out, ", ") + "]"
}

// Used when iterating over each object in the list
func (list *Library_List) Length() int {
	return len(list.Items)
}

// Used when iterating over each object in the list
func (list *Library_List) ItemAtPosition(pos int) Type {
	return list.Items[pos]
}
