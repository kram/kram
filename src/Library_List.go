package main

import (
	"strings"
)

type Library_List struct {
	*Library
	Items []Type
}

func (list *Library_List) Instance() (Lib, string) {
	return &Library_List{}, "List"
}

func (list *Library_List) Init(vm *VM, params []Type) {
	list.Items = make([]Type, 0)
	list.Push(vm, params)
}

func (list *Library_List) Push(vm *VM, params []Type) {

	for _, param := range params {
		list.Items = append(list.Items, param)
	}
}

func (list *Library_List) Pop(vm *VM, params []Type) Type {
	res := list.Items[len(list.Items)-1]
	list.Items = list.Items[:len(list.Items)-1]

	return res
}

func (list *Library_List) ToString() string {

	out := make([]string, len(list.Items))

	for i, item := range list.Items {
		out[i] = item.ToString()
	}

	return "[" + strings.Join(out, ", ") + "]"
}

func (list *Library_List) Length() int {
	return len(list.Items)
}

func (list *Library_List) ItemAtPosition(pos int) Type {
	return list.Items[pos]
}