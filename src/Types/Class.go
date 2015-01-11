package types

import (
	"../Instructions"
	"fmt"
	"log"
)

type Method struct {
	Method     bool
	Parameters []instructions.Parameter
	Body       instructions.Block
	IsStatic   bool
	IsPublic   bool
}

type Class struct {
	Class     string
	Methods   map[string]Method
	Variables map[string]Type
}

func (self *Class) Init(str string) {
	self.Class = str
	self.Methods = make(map[string]Method)
	self.Variables = make(map[string]Type)
}

func (self *Class) AddMethod(name string, method Method) {
	self.Methods[name] = method
}

func (self *Class) SetVariable(name string, value Type) {
	self.Variables[name] = value
}

func (self *Class) Type() string {
	return self.Class
}

func (self *Class) ToString() string {
	return self.Type() + "\n" + fmt.Sprint(self.Variables)
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
