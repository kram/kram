package main

import (
	"log"
)

type Method struct {
	Method     bool
	Parameters []Parameter
	Body       Block
	IsStatic   bool
}

type Class struct {
	Class   string
	Methods map[string]Method
}

func (self *Class) Init(str string) {
	self.Class = str
	self.Methods = make(map[string]Method)
}

func (self *Class) AddMethod(name string, method Method) {
	self.Methods[name] = method
}

func (self *Class) Type() string {
	return self.Class
}

func (self *Class) toString() string {
	return "Class"
}

func (self *Class) Math(method string, right Type) Type {

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// This code will never be reached

	res := Bool{}
	res.Init("false")

	return &res
}

func (self *Class) Compare(method string, right Type) Type {

	log.Panicf("You can not compare a %s() with a %s()", self.Type(), right.Type())

	// Will never be reached

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// Will never be reached

	bl := Bool{}
	bl.Init("false")

	return &bl
}
