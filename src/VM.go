// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package gus

import (
	"log"

	"github.com/zegl/Gus/src/environment"
	ins "github.com/zegl/Gus/src/instructions"
	lib "github.com/zegl/Gus/src/libraries"
	"github.com/zegl/Gus/src/types"
	"github.com/zegl/Gus/src/types/builtin"
)

type VM struct {
	// Contains variables
	env *environment.Environment

	// The current stack of methods, used to know where to define a method
	Classes []*types.Type

	Debug bool

	ShouldReturn []bool
}

func (vm *VM) Run(tree ins.Block) {

	// Set empty environment
	vm.env = &environment.Environment{}
	vm.env.Init()

	// Create empty lists
	vm.ShouldReturn = make([]bool, 0)
	vm.Classes = make([]*types.Type, 0)

	vm.Libraries()

	vm.Operation(tree, types.ON_NOTHING)
}

func (vm *VM) EnvironmentPush() {
	vm.env = vm.env.Push()
}

func (vm *VM) EnvironmentPop() {
	vm.env = vm.env.Pop()
}

func (vm *VM) Libraries() {

	libs := make([]types.Lib, 0)

	libs = append(libs, &lib.Library_IO{})
	libs = append(libs, &lib.Library_String{})
	libs = append(libs, &lib.Library_File{})

	for _, li := range libs {

		instance, name := li.Instance()

		class := types.Type{}
		class.Init(name)
		class.Extension = instance

		vm.env.Set(name, &class)
	}
}

func (vm *VM) Operation(node ins.Node, on types.ON) *types.Type {

	if assign, ok := node.(ins.Assign); ok {
		return vm.OperationAssign(assign)
	}

	if math, ok := node.(ins.Math); ok {
		return vm.OperationMath(math)
	}

	if literal, ok := node.(ins.Literal); ok {
		return vm.OperationLiteral(literal)
	}

	if variable, ok := node.(ins.Variable); ok {

		if on == types.ON_CLASS {
			return vm.ClassOperationVariable(variable)
		}

		return vm.OperationVariable(variable)
	}

	if set, ok := node.(ins.Set); ok {

		if on == types.ON_CLASS {
			return vm.ClassOperationSet(set)
		}

		return vm.OperationSet(set)
	}

	if i, ok := node.(ins.If); ok {
		return vm.OperationIf(i)
	}

	if block, ok := node.(ins.Block); ok {
		return vm.OperationBlock(block, on)
	}

	if call, ok := node.(ins.Call); ok {
		return vm.OperationCall(call)
	}

	if pushClass, ok := node.(ins.PushClass); ok {
		return vm.OperationPushClass(pushClass)
	}

	if defineClass, ok := node.(ins.DefineClass); ok {
		return vm.OperationDefineClass(defineClass)
	}

	if defineMethod, ok := node.(ins.DefineMethod); ok {
		return vm.OperationDefineMethod(defineMethod)
	}

	if instance, ok := node.(ins.Instance); ok {
		return vm.OperationInstance(instance)
	}

	if list, ok := node.(ins.ListCreate); ok {
		return vm.OperationListCreate(list)
	}

	if access, ok := node.(ins.AccessChildItem); ok {
		return vm.OperationAccessChildItem(access)
	}

	if m, ok := node.(ins.MapCreate); ok {
		return vm.OperationMapCreate(m)
	}

	if ret, ok := node.(ins.Return); ok {
		return vm.OperationReturn(ret)
	}

	if f, ok := node.(ins.For); ok {
		return vm.OperationFor(f)
	}

	return vm.CreateType(&builtin.Null{})
}

func (vm *VM) OperationBlock(block ins.Block, on types.ON) (last *types.Type) {

	// Create new scope
	if on != types.ON_FOR_PART && block.Scope == true {
		vm.env = vm.env.Push()
	}

	if on == types.ON_METHOD_BODY {
		vm.ShouldReturn = append(vm.ShouldReturn, false)
	}

	for _, body := range block.Body {
		last = vm.Operation(body, types.ON_NOTHING)

		// Return statement
		if len(vm.ShouldReturn) > 0 && vm.ShouldReturn[len(vm.ShouldReturn)-1] {
			break
		}
	}

	// Pop
	if on == types.ON_METHOD_BODY {
		vm.ShouldReturn = vm.ShouldReturn[:len(vm.ShouldReturn)-1]
	}

	// Restore scope
	if on != types.ON_FOR_PART && block.Scope == true {
		vm.env = vm.env.Pop()
	}

	return last
}

func (vm *VM) OperationAssign(assign ins.Assign) *types.Type {

	var value *types.Type

	// Assign.Right is already a Type{}
	// Used in a ForIn for example
	if t, ok := assign.Right.(*types.Type); ok {
		value = t
	} else {
		value = vm.Operation(assign.Right, types.ON_NOTHING)
	}

	vm.env.Set(assign.Name, value)

	return value
}

func (vm *VM) OperationMath(math ins.Math) *types.Type {

	left := vm.Operation(math.Left, types.ON_NOTHING)
	right := vm.Operation(math.Right, types.ON_NOTHING)

	/*fmt.Println(left)
	fmt.Println(math.Method)
	fmt.Println(right)
	*/
	if math.IsComparision {
		return left.Compare(vm, math.Method, right)
	}

	return left.Math(vm, math.Method, right)
}

func (vm *VM) OperationLiteral(literal ins.Literal) *types.Type {

	if literal.Type == "number" {
		number := builtin.Number{}
		number.Init(literal.Value)

		return vm.CreateType(&number)
	}

	if literal.Type == "string" {
		str := builtin.String{}
		str.Init(literal.Value)

		return vm.CreateType(&str)
	}

	if literal.Type == "bool" {
		bl := builtin.Bool{}
		bl.Init(literal.Value)

		return vm.CreateType(&bl)
	}

	if literal.Type == "null" {
		return vm.CreateType(&builtin.Null{})
	}

	log.Panicf("Not able to handle Literal %s", literal)

	return vm.CreateType(&builtin.Null{})
}

func (vm *VM) OperationVariable(variable ins.Variable) *types.Type {

	if res, ok := vm.env.Get(variable.Name); ok {
		return res
	}

	log.Print("Undefined variable, " + variable.Name)

	return vm.CreateType(&builtin.Null{})
}

func (vm *VM) ClassOperationVariable(variable ins.Variable) *types.Type {

	class := vm.Classes[len(vm.Classes)-1]

	if res, ok := class.Variables[variable.Name]; ok {
		return res
	}

	log.Print("Undefined variable, " + class.Type() + "." + variable.Name)

	return vm.CreateType(&builtin.Null{})
}

func (vm *VM) OperationSet(set ins.Set) *types.Type {

	l, ok := vm.env.Get(set.Name)

	if !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right, types.ON_NOTHING)

	if l.Type() != value.Type() {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, l.Type(), value.ToString(), value.Type())
	}

	vm.env.Set(set.Name, value)

	return value
}

func (vm *VM) ClassOperationSet(set ins.Set) *types.Type {

	class := vm.Classes[len(vm.Classes)-1]

	l, ok := class.Variables[set.Name]

	if !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right, types.ON_NOTHING)

	if l.Type() != "Null" && l.Type() != value.Type() {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, l.Type(), value.ToString(), value.Type())
	}

	class.SetVariable(set.Name, value)

	return value
}

func (vm *VM) OperationIf(i ins.If) *types.Type {

	con := vm.Operation(i.Condition, types.ON_NOTHING)

	if con.Type() != "Bool" {
		log.Panicf("Expecing bool in condition, %s (%s)", con.ToString(), con.Type())
	}

	if con.ToString() == "true" {
		return vm.Operation(i.True, types.ON_NOTHING)
	}

	return vm.Operation(i.False, types.ON_NOTHING)
}

func (vm *VM) OperationCall(call ins.Call) *types.Type {

	params := make([]*types.Type, len(call.Parameters))

	for i, param := range call.Parameters {
		params[i] = vm.Operation(param, types.ON_NOTHING)
	}

	return vm.Classes[len(vm.Classes)-1].Invoke(vm, vm.Operation(call.Left, types.ON_NOTHING).ToString(), params)
}

func (vm *VM) OperationDefineClass(def ins.DefineClass) *types.Type {

	class := types.Type{}
	class.Init(def.Name)

	// Push
	vm.Classes = append(vm.Classes, &class)
	vm.env = vm.env.Push()

	for _, body := range def.Body.Body {

		/*if assign, ok := body.(Assign); ok {
			class.SetVariable(assign.Name, vm.Operation(assign.Right, ON_NOTHING))
			continue
		}*/

		// Fallback
		vm.Operation(body, types.ON_NOTHING)
	}

	// Pop
	vm.Classes = vm.Classes[:len(vm.Classes)-1]
	vm.env = vm.env.Pop()

	vm.env.Set(def.Name, &class)

	return vm.CreateType(&builtin.Null{})
}

func (vm *VM) OperationDefineMethod(def ins.DefineMethod) *types.Type {

	if len(vm.Classes) == 0 {
		log.Panic("Unable to define method, not in class")
	}

	method := types.Method{}
	method.Parameters = def.Parameters
	method.Body = def.Body
	method.IsStatic = def.IsStatic
	method.IsPublic = def.IsPublic

	vm.Classes[len(vm.Classes)-1].AddMethod(def.Name, method)

	return vm.CreateType(&builtin.Null{})
}

func (vm *VM) OperationPushClass(pushClass ins.PushClass) *types.Type {

	left := vm.Operation(pushClass.Left, types.ON_NOTHING)
	name := left.ToString()

	class, ok := vm.env.Get(name)

	if name == "self" {
		class = vm.Classes[len(vm.Classes)-1]
		ok = true
	}

	// There is no such class, use left
	if !ok {
		class = left
	}

	// Push
	vm.Classes = append(vm.Classes, class)

	res := vm.Operation(pushClass.Right, types.ON_CLASS)

	// Pop
	vm.Classes = vm.Classes[:len(vm.Classes)-1]

	return res
}

func (vm *VM) OperationMapCreate(m ins.MapCreate) *types.Type {
	ma := builtin.Map{}

	params := make([]*types.Type, 0)

	for i, key := range m.Keys {
		params = append(params, vm.Operation(key, types.ON_NOTHING))
		params = append(params, vm.Operation(m.Values[i], types.ON_NOTHING))
	}

	ma.InitWithParams(params)

	return vm.CreateType(&ma)
}

func (vm *VM) OperationListCreate(list ins.ListCreate) *types.Type {
	l := builtin.List{}

	params := make([]*types.Type, len(list.Items))

	for i, item := range list.Items {
		params[i] = vm.Operation(item, types.ON_NOTHING)
	}

	l.InitWithParams(params)

	return vm.CreateType(&l)
}

func (vm *VM) OperationAccessChildItem(access ins.AccessChildItem) *types.Type {

	// Extract the List or Map
	item := vm.Operation(access.Item, types.ON_NOTHING)

	// Is Map
	if item.Type() == "Map" {
		return vm.OperationAccessChildItemMap(access, item)
	}

	if item.Type() != "List" {
		log.Panicf("Expected List or Map in [], got %s", item.Type())
	}

	library, ok := item.Extension.(*builtin.List)

	if !ok {
		log.Panic("Expected class to be of types.Type *builtin.List")
	}

	// Get position to access from the list
	position := vm.Operation(access.Right, types.ON_NOTHING)

	return library.ItemAt([]*types.Type{position})
}

func (vm *VM) OperationAccessChildItemMap(access ins.AccessChildItem, item *types.Type) *types.Type {
	library, ok := item.Extension.(*builtin.Map)

	if !ok {
		log.Panic("Expected class to be of types.Type *builtin.Map")
	}

	// Get position to access from the list
	position := vm.Operation(access.Right, types.ON_NOTHING)

	return library.Get([]*types.Type{position})
}

func (vm *VM) OperationReturn(ret ins.Return) *types.Type {

	vm.ShouldReturn[len(vm.ShouldReturn)-1] = true

	return vm.Operation(ret.Statement, types.ON_NOTHING)
}

func (vm *VM) OperationInstance(instance ins.Instance) *types.Type {

	in, ok := vm.env.Get(instance.Left)

	if !ok {
		log.Panicf("No such class, %s", instance.Left)
	}

	return vm.Clone(in)
}

//
// for before; condition; each { body }
//
func (vm *VM) OperationFor(f ins.For) *types.Type {

	if f.IsForIn {
		return vm.OperationForIn(f)
	}

	// Create variable scope
	vm.env = vm.env.Push()

	// Execute before part
	vm.Operation(f.Before, types.ON_FOR_PART)

	for {

		// Test condition
		res := vm.Operation(f.Condition, types.ON_FOR_PART)

		condition, is_bool := res.Extension.(*builtin.Bool)

		if !is_bool {
			log.Panicf("Expected bool in for, got %s", res.Type())
		}

		if !condition.IsTrue() {
			break
		}

		// Execute body
		vm.Operation(f.Body, types.ON_NOTHING)

		// Execute part after each run
		vm.Operation(f.Each, types.ON_FOR_PART)
	}

	// Restore scope
	vm.env = vm.env.Pop()

	return vm.CreateType(&builtin.Null{})
}

// for var item in 1..2
// for var item in ["first", "second"]
// for var item in list
func (vm *VM) OperationForIn(f ins.For) *types.Type {

	// Create variable scope
	vm.env = vm.env.Push()

	// Convert Before to an assign object
	assign, assign_ok := f.Before.(ins.Assign)

	// Get iterator object
	each := vm.Operation(f.Each, types.ON_NOTHING)

	if each.Type() != "List" {
		log.Panic("Expected List in for ... in, got %s", each.Type())
	}

	list, ok := each.Extension.(*builtin.List)

	if !ok {
		log.Panic("Expected class to be of types.Type *builtin.List")
	}

	length := list.Length()

	for key := 0; key < length; key++ {

		// Update variable
		if assign_ok {
			item := list.ItemAtPosition(key)
			assign.Right = item
			vm.Operation(assign, types.ON_NOTHING)
		}

		// Run body
		vm.Operation(f.Body, types.ON_NOTHING)
	}

	// Restore scope
	vm.env = vm.env.Pop()

	return vm.CreateType(&builtin.Null{})
}

//
// Clones a type, returns the new one
//
func (vm *VM) Clone(in *types.Type) (out *types.Type) {
	res := types.Type{}
	res.Class = in.Class
	res.Methods = in.Methods
	res.Extension, _ = in.Extension.Instance()
	res.Variables = make(map[string]*types.Type)

	for name, def := range in.Variables {
		res.Variables[name] = vm.Clone(def)
	}

	return &res
}

func (vm *VM) CreateType(lib types.Lib) *types.Type {
	class := types.Type{}
	class.InitWithLib(lib)

	return &class
}
