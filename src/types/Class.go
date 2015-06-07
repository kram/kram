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

type Argument struct {
	IsNamed bool
	Name    string
	Val     *Value
}

type Class struct {
	Class     string
	Methods   map[string]Method
	Variables map[string]*Value
	Extension Lib
	HasExtension bool
	IsInstance bool
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

// A safe version of Invoke() that doesn't panic
func (self *Class) InvokeSafe(vm VM, name string, arguments []Argument) *Value {
	return self.invoke(vm, name, arguments, false)
}

// Invoke and panik if method not found
func (self *Class) Invoke(vm VM, name string, arguments []Argument) *Value {
	return self.invoke(vm, name, arguments, true)
}

// Private method for invoking a Gus-method
func (self *Class) invoke(vm VM, name string, arguments []Argument, panic bool) *Value {

	// Special method
	if name == "Type" {
		return self.M_Type()
	}

	res, ok := self.InvokeNative(vm, name, arguments)

	if ok {
		return res
	}

	res, ok = self.InvokeExtension(vm, name, arguments)

	if ok {
		return res
	}

	if panic {
		log.Panicf("%s::%s, no such method", self.Type(), name)
	}

	return self.CreateNull()
}

func (self *Class) InvokeExtension(vm VM, method string, arguments []Argument) (*Value, bool) {

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
		par := make([]*Class, len(arguments))

		for k, v := range arguments {
			par[k] = vm.GetAsClass(v.Val)
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

func (self *Class) InvokeNative(vm VM, name string, arguments []Argument) (*Value, bool) {

	method, ok := self.Methods[name]

	if !ok {
		return self.CreateNull(), false
	}

	if !method.IsStatic && !self.IsInstance {
		log.Panicf("%s::%s is not a static method and needs to be called from a class instance", self.Type(), name)
	}

	required_params := 0

	for _, par := range method.Parameters {
		if !par.HasDefault {
			required_params++
		}
	}

	if required_params > len(arguments) {
		fmt.Printf("Can not call %s.%s() (%d parameters) with %d parameters\n", self.ToString(), name, len(method.Parameters), len(arguments))

		return self.CreateNull(), true
	}

	named_arguments := 0

	for _, arg := range arguments {
		if arg.IsNamed {
			named_arguments++
		}
	}

	if named_arguments > 0 && named_arguments != len(arguments) {
		fmt.Printf("Can not call %s.%s(), when using named arguments you need to name all arguments\n", self.ToString(), name)

		return self.CreateNull(), true
	}

	vm.EnvironmentPush()

	// Not using named arguments
	if named_arguments == 0 {

		// Define variables
		for i, param := range method.Parameters {
			ass := ins.Assign{}
			ass.Name = param.Name

			if len(arguments) > i {
				ass.Right = arguments[i].Val
			} else {
				ass.Right = param.Default
			}

			vm.Assign(ass)
		}
	} else {

		// Named arguments
		for _, param := range method.Parameters {
			ass := ins.Assign{}
			ass.Name = param.Name

			// Find name
			found_name := false
			for _, v := range arguments {
				if v.Name == ass.Name {
					ass.Right = v.Val
					found_name = true
					break
				}
			}

			if !found_name {
				if param.HasDefault {
					ass.Right = param.Default
				} else {
					log.Panic("Something something")
				}
			}

			vm.Assign(ass)
		}
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

	res, ok := self.InvokeNative(vm, method, []Argument{Argument{
 		Val: right,
	}})

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

	res, ok := self.InvokeNative(vm, method, []Argument{Argument{
		Val: right,
	}})

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