package main

import (
	"fmt"
)

type IO struct {
	*Library
}

func (io *IO) Instance() (Lib, string) {
	return &IO{}, "IO"
}

func (io IO) Print(vm *VM, params []Type) {

	for _, param := range params {
		fmt.Print(param.ToString())
	}
}

func (io IO) Println(vm *VM, params []Type) {

	for _, param := range params {
		fmt.Print(param.ToString())
	}
}