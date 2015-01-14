package gus

import (
	"fmt"
)

type IO struct{}

func (self IO) Print(params ...Node) Type {

	for _, param := range params {
		// fmt.Print(vm.Operation(param, ON_NOTHING).ToString())
		fmt.Print(param)
	}

	return DefaultReturn()
}

func (self IO) Println(params ...Node) Type {

	for _, param := range params {
		// fmt.Println(vm.Operation(param, ON_NOTHING).ToString())
		fmt.Print(param)
	}

	return DefaultReturn()
}
