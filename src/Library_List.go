package main

import (
	"strings"
)

type List struct {
	*Library
	Items []Type
}

func (list *List) Instance() (Lib, string) {
	return &List{}, "List"
}

func (list *List) Init(vm *VM, params []Type) {
	list.Items = make([]Type, 0)
	list.Push(vm, params)
}

func (list *List) Push(vm *VM, params []Type) {

	for _, param := range params {
		list.Items = append(list.Items, param)
	}
}

func (list *List) Pop(vm *VM, params []Type) Type {
	res := list.Items[len(list.Items)-1]
	list.Items = list.Items[:len(list.Items)-1]

	return res
}

func (list *List) ToString() string {

	out := make([]string, len(list.Items))

	for i, item := range list.Items {
		out[i] = item.ToString()
	}

	return "[" + strings.Join(out, ", ") + "]"
}
