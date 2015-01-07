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

func (class *Class) Init(str string) {
	class.Class = str
	class.Methods = make(map[string]Method)
}

func (class *Class) AddMethod(name string, method Method) {
	class.Methods[name] = method
}

func (class *Class) Type() string {
	return class.Class
}

func (class *Class) Math(method string, right Type) Type {

	log.Panicf("Class() is not implementing %s", method)

	// This code will never be reached

	res := Bool{}
	res.Init("false")

	return &res
}

func (class *Class) toString() string {
	return "Class"
}
