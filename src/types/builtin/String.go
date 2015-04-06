package builtin

import (
	"github.com/zegl/Gus/src/types"
	"log"
)

type String struct {
	Value string
}

func (self String) Instance() (types.Lib, string) {
	return &String{}, self.Type()
}

func (self String) Type() string {
	return "String"
}

func (self *String) Init(str string) {
	self.Value = str
}

func (self *String) ToString() string {
	return self.Value
}

func (self *String) Math(method string, right *types.Type) *types.Type {

	r, ok := right.Extension.(*String)

	if !ok {
		log.Panicf("You can not apply %s to a %s() with a %s()", method, self.Type(), right.Type())
	}

	// String concatenation
	if method == "+" {
		str := String{}
		str.Init(self.Value + r.Value)

		res := types.Type{}
		res.InitWithLib(&str)

		return &res
	}

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// This code will never be reached

	return &types.Type{}
}

func (self *String) Compare(method string, right *types.Type) *types.Type {

	r, ok := right.Extension.(*String)

	if !ok {
		log.Panicf("You can not compare a %s() with a %s()", self.Type(), right.Type())
	}

	b := false

	switch method {
	case ">":
		b = self.Value > r.Value
	case "<":
		b = self.Value < r.Value
	case ">=":
		b = self.Value >= r.Value
	case "<=":
		b = self.Value <= r.Value
	case "==":
		b = self.Value == r.Value
	case "!=":
		b = self.Value != r.Value
	default:
		log.Panicf("%s() is not implementing %s", self.Type(), method)
	}

	bl := Bool{}
	bl.Set(b)

	res := types.Type{}
	res.InitWithLib(&bl)

	return &res
}
