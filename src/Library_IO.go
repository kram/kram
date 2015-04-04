package main

import (
	"fmt"
)

type Library_IO struct {
	*Library
}

func (io *Library_IO) Instance() (Lib, string) {
	return &Library_IO{}, "IO"
}

func (io Library_IO) Print(params []Type) {
	for _, param := range params {
		fmt.Print(param.ToString())
	}
}

func (io Library_IO) Println(params []Type) {
	for _, param := range params {
		fmt.Println(param.ToString())
	}
}
