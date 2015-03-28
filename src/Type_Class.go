package main

import (
	"fmt"
	"log"
	"reflect"
)

type Method struct {
	Method     bool
	Parameters []Parameter
	Body       Block
	IsStatic   bool
	IsPublic   bool
}

type Class struct {
	Class     string
	Methods   map[string]Method
	Variables map[string]Type
	Extension Lib
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

func (self *Class) Invoke(vm *VM, name string, params []Node) Type {

	res, ok := self.InvokeNative(vm, name, params)

	if ok {
		return res
	}

	return self.InvokeExtension(vm, name, params)
}

func (self *Class) InvokeExtension(vm *VM, method string, params []Node) Type {

	inputs := make([]reflect.Value, 2)

	inputs[0] = reflect.ValueOf(vm)

	param_type := make([]Type, len(params))

	for i, param := range params {
		param_type[i] = vm.Operation(param, ON_NOTHING)
	}

	inputs[1] = reflect.ValueOf(param_type)

	res := reflect.ValueOf(self.Extension).MethodByName(method).Call(inputs)

	if len(res) > 0 {
		return res[0].Interface().(Type)
	}

	return &Null{}
}

func (self *Class) InvokeNative(vm *VM, name string, params []Node) (Type, bool) {

	method, ok := self.Methods[name]

	if !ok {
		return &Bool{}, false
	}

	if len(method.Parameters) != len(params) {
		fmt.Printf("Can not call %s.%s() (%d parameters) with %d parameters\n", self.ToString(), name, len(method.Parameters), len(params))

		return &Bool{}, true
	}

	vm.Environment = vm.Environment.Push()

	// Define variables
	for i, param := range method.Parameters {
		ass := Assign{}
		ass.Name = param.Name
		ass.Right = params[i]

		vm.OperationAssign(ass)
	}

	body := vm.OperationBlock(method.Body, ON_METHOD_BODY)

	vm.Environment = vm.Environment.Push()

	return body, true
}

func (self *Class) Type() string {
	return self.Class
}

func (self *Class) ToString() string {

	if _, ok := self.Extension.(Lib); ok {
		return self.Extension.ToString()
	}

	return self.Class
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
