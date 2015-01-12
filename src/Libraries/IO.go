package libraries

import (
	"../Instructions"
	"../Types"
	"fmt"
)

type IO struct{}

func (self IO) Print(params ...instructions.Node) types.Type {

	for _, param := range params {
		// fmt.Print(vm.Operation(param, ON_NOTHING).ToString())
		fmt.Print(param)
	}

	return false()
}

func (self IO) Println(params ...instructions.Node) types.Type {

	for _, param := range params {
		// fmt.Println(vm.Operation(param, ON_NOTHING).ToString())
		fmt.Print(param)
	}

	return false()
}
