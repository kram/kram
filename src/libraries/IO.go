package libraries

import (
	"fmt"
	"github.com/zegl/Gus/src/types"
)

type Library_IO struct {
	*Library
}

func (io *Library_IO) Instance() (types.Lib, string) {
	return &Library_IO{}, "IO"
}

func (io Library_IO) Print(params []*types.Type) {
	for _, param := range params {
		fmt.Print(param.ToString())
	}
}

func (io Library_IO) Println(params []*types.Type) {
	for _, param := range params {
		fmt.Println(param.ToString())
	}
}
