// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package types

import (
	"fmt"
	ins "github.com/zegl/Gus/src/instructions"
	"log"
	"reflect"
)

type Lib interface {
	Init(string)
	InitWithParams([]*Class)
	Instance() (Lib, string)
	Type() string
	ToString() string
}

type MathLib interface {
	Math(string, *Class) *Class
	Compare(string, *Class) *Class
}

type Method struct {
	Method     bool
	Parameters []ins.Parameter
	Body       ins.Block
	IsStatic   bool
	IsPublic   bool
}

type Class struct {
	Class     string
	Methods   map[string]Method
	Variables map[string]*Value
	Extension Lib
	HasExtension bool
}

func (self Class) IsClass() bool {
	return true
}

func (self *Class) Init(str string) {
	self.Class = str
	self.Methods = make(map[string]Method)
	self.Variables = make(map[string]*Value)
}

func (self *Class) InitWithLib(lib Lib) {
	self.Init(lib.Type())
	self.Extension = lib	
	self.HasExtension = true
}

func (self *Class) AddMethod(name string, method Method) {
	self.Methods[name] = method
}

func (self *Class) SetVariable(name string, value *Value) {
	self.Variables[name] = value
}

func (self *Class) Invoke(vm VM, name string, params []*Value) *Value {

	if name == "Type" {
		return self.M_Type()
	}

	res, ok := self.InvokeNative(vm, name, params)

	if ok {
		return res
	}

	res, ok = self.InvokeExtension(vm, name, params)

	if ok {
		return res
	}

	log.Panicf("%s::%s, no such method", self.Type(), name)

	return self.CreateNull()
}

func (self *Class) InvokeExtension(vm VM, method string, params []*Value) (*Value, bool) {

	if !self.HasExtension {
		return self.CreateNull(), false
	}

	value := reflect.ValueOf(self.Extension).MethodByName("M_" + method)

	if !value.IsValid() {
		return self.CreateNull(), false
	}

	var res []reflect.Value

	// The list as a parameter
	// This should probaby be rewritten so that we can use parameters properly...
	// Eg. a parameter in Gus => a parameter in Go
	if value.Type().NumIn() == 1 {

		// Convert params to []*Class
		par := make([]*Class, len(params))

		for k, v := range params {
			par[k] = vm.GetAsClass(v)
		}

		inputs := make([]reflect.Value, 1)
		inputs[0] = reflect.ValueOf(par)
		res = value.Call(inputs)
	} else {
		inputs := make([]reflect.Value, 0)
		res = value.Call(inputs)
	}

	if len(res) > 0 {
		return vm.ConvertClassToValue(res[0].Interface().(*Class)), true
	}

	// Nothing was returned, but still valid
	return self.CreateNull(), true
}

func (self *Class) InvokeNative(vm VM, name string, params []*Value) (*Value, bool) {

	method, ok := self.Methods[name]

	if !ok {
		return self.CreateNull(), false
	}

	//param_count := len(method.Parameters)
	// optional_parameters := 0
	required_params := 0

	for _, par := range method.Parameters {
		if !par.HasDefault {
			required_params++
		}
	}

	if required_params > len(params) {
		fmt.Printf("Can not call %s.%s() (%d parameters) with %d parameters\n", self.ToString(), name, len(method.Parameters), len(params))

		return self.CreateNull(), true
	}

	vm.EnvironmentPush()

	// Define variables
	for i, param := range method.Parameters {
		ass := ins.Assign{}
		ass.Name = param.Name

		if len(params) > i {
			ass.Right = params[i]
		} else {
			ass.Right = param.Default
		}

		vm.Assign(ass)
	}

	body := vm.Block(method.Body, ON_METHOD_BODY)

	vm.EnvironmentPop()

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

func (self *Class) Math(vm VM, method string, right *Value) *Value {

	if lib, ok := self.Extension.(MathLib); ok {
		return vm.ConvertClassToValue(lib.Math(method, vm.GetAsClass(right)))
	}

	res, ok := self.InvokeNative(vm, method, []*Value{right})

	if ok {
		return res
	}

	log.Panicf("%s() is not implementing %s (General Math)", self.Type(), method)

	// This code will never be reached
	return self.CreateNull()
}

func (self *Class) Compare(vm VM, method string, right *Value) *Value {

	if lib, ok := self.Extension.(MathLib); ok {
		return vm.ConvertClassToValue(lib.Compare(method, vm.GetAsClass(right)))
	}

	res, ok := self.InvokeNative(vm, method, []*Value{right})

	if ok {
		return res
	}

	log.Panicf("%s() is not implementing %s (General Compare)", self.Type(), method)

	// This code will never be reached
	return self.CreateNull()
}

func (self Class) CreateNull() *Value {
	return &Value{
		Type: NULL,
	}
}

func (self *Class) M_Type() *Value {
	return &Value{
		Type: STRING,
		String: self.Type(),
	}
}