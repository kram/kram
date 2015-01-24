package main

import (
	"strings"
)

type Library_String struct {
	*Library
}

func (self *Library_String) Instance() (Lib, string) {
	return &Library_String{}, "String"
}

func (self Library_String) ToLower(vm *VM, params []Type) Type {

	str := String{}

	for _, param := range params {
		str.Init(strings.ToLower(param.ToString()))
		return &str
	}

	return &str
}

func (self Library_String) ToUpper(vm *VM, params []Type) Type {

	str := String{}

	for _, param := range params {
		str.Init(strings.ToUpper(param.ToString()))
		return &str
	}

	return &str
}