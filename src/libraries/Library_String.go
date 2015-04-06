package libraries

import (
	"strings"
	"../types"
	"../types/builtin"
)

type Library_String struct {
	*Library
}

func (self *Library_String) Instance() (types.Lib, string) {
	return &Library_String{}, "String"
}

func (self Library_String) ToLower(params []*types.Type) *types.Type {
	str := builtin.String{}

	for _, param := range params {
		str.Init(strings.ToLower(param.ToString()))
		break
	}

	return self.TypeWithLib(&str)
}

func (self Library_String) ToUpper(params []*types.Type) *types.Type {
	str := builtin.String{}

	for _, param := range params {
		str.Init(strings.ToUpper(param.ToString()))
		break
	}

	return self.TypeWithLib(&str)
}
