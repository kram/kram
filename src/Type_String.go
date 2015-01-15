package main

import (
	"log"
)

type String struct {
	String bool
	Value  string
}

func (self *String) Init(str string) {
	self.Value = str
}

func (self String) Type() string {
	return "String"
}

func (self *String) Math(method string, right Type) Type {

	r, ok := right.(*String)

	if !ok {
		log.Panicf("You can not apply %s to a %s() with a %s()", method, self.Type(), right.Type())
	}

	// String concatenation
	if method == "+" {
		str := String{}
		str.Init(self.Value + r.Value)

		return &str
	}

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// This code will never be reached

	res := Bool{}
	res.Init("false")

	return &res
}

func (self *String) Compare(method string, right Type) Type {

	r, ok := right.(*String)

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
	bl.Value = b

	return &bl
}

func (self *String) ToString() string {
	return self.Value
}
