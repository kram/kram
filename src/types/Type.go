package types

import (
	"fmt"
	"log"
	"reflect"
	ins "../instructions"
)

type Lib interface {
	Init(string)
	Instance() (Lib, string)
	Type() string
	ToString() string
}

type Method struct {
	Method     bool
	Parameters []ins.Parameter
	Body       ins.Block
	IsStatic   bool
	IsPublic   bool
}

type Type struct {
	Class     string
	Methods   map[string]Method
	Variables map[string]*Type
	Extension Lib
}

func (self *Type) Init(str string) {
	self.Class = str
	self.Methods = make(map[string]Method)
	self.Variables = make(map[string]*Type)
}

func (self *Type) InitWithLib(lib Lib) {
	self.Init(lib.Type())
	self.Extension = lib
}

func (self *Type) AddMethod(name string, method Method) {
	self.Methods[name] = method
}

func (self *Type) SetVariable(name string, value *Type) {
	self.Variables[name] = value
}

func (self *Type) Invoke(vm VM, name string, params []ins.Node) *Type {

	res, ok := self.InvokeNative(vm, name, params)

	if ok {
		return res
	}

	return self.InvokeExtension(vm, name, params)
}

func (self *Type) InvokeExtension(vm VM, method string, params []ins.Node) *Type {

	inputs := make([]reflect.Value, 1)

	param_type := make([]*Type, len(params))

	for i, param := range params {
		param_type[i] = vm.Operation(param, ON_NOTHING)
	}

	inputs[0] = reflect.ValueOf(param_type)

	res := reflect.ValueOf(self.Extension).MethodByName(method).Call(inputs)

	if len(res) > 0 {
		return res[0].Interface().(*Type)
	}

	return &Type{}
}

func (self *Type) InvokeNative(vm VM, name string, params []ins.Node) (*Type, bool) {

	method, ok := self.Methods[name]

	if !ok {
		return &Type{}, false
	}

	if len(method.Parameters) != len(params) {
		fmt.Printf("Can not call %s.%s() (%d parameters) with %d parameters\n", self.ToString(), name, len(method.Parameters), len(params))

		return &Type{}, true
	}

	vm.EnvironmentPush()

	// Define variables
	for i, param := range method.Parameters {
		ass := ins.Assign{}
		ass.Name = param.Name
		ass.Right = params[i]

		vm.OperationAssign(ass)
	}

	body := vm.OperationBlock(method.Body, ON_METHOD_BODY)

	vm.EnvironmentPop()

	return body, true
}

func (self *Type) Type() string {
	return self.Class
}

func (self *Type) ToString() string {

	if _, ok := self.Extension.(Lib); ok {
		return self.Extension.ToString()
	}

	return self.Class
}

func (self *Type) Math(method string, right *Type) *Type {

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// This code will never be reached

	return &Type{}
}

func (self *Type) Compare(method string, right *Type) *Type {

	log.Panicf("You can not compare a %s() with a %s()", self.Type(), right.Type())

	// Will never be reached

	log.Panicf("%s() is not implementing %s", self.Type(), method)

	// Will never be reached

	return &Type{}
}