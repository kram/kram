// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package gus

import (
	"log"
	"strconv"
	"reflect"
	"fmt"
	"math"

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
	Classes []*types.Class

	Debug bool

	ShouldReturn []bool
}

func (vm *VM) Run(tree ins.Block) {

	// Set empty environment
	vm.env = &environment.Environment{}
	vm.env.Init()

	// Create empty lists
	vm.ShouldReturn = make([]bool, 0)
	vm.Classes = make([]*types.Class, 0)

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

	// Builtin
	libs = append(libs, &builtin.Bool{})
	libs = append(libs, &builtin.List{})
	libs = append(libs, &builtin.Map{})
	libs = append(libs, &builtin.Null{})
	libs = append(libs, &builtin.Number{})
	libs = append(libs, &builtin.String{})

	// Libraries
	libs = append(libs, &lib.Library_IO{})
	libs = append(libs, &lib.Library_String{})
	libs = append(libs, &lib.Library_File{})

	for _, li := range libs {

		instance, name := li.Instance()

		class := types.Class{}
		class.Init(name)
		class.Extension = instance

		vm.env.Set(name, &class)
	}
}

func (vm *VM) Operation(node ins.Node, on types.ON) types.Type {

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

func (vm *VM) OperationBlock(block ins.Block, on types.ON) types.Type {

	// Create new scope
	if on != types.ON_FOR_PART && block.Scope == true {
		vm.env = vm.env.Push()
	}

	if on == types.ON_METHOD_BODY {
		vm.ShouldReturn = append(vm.ShouldReturn, false)
	}

	var last types.Type

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

func (vm *VM) OperationAssign(assign ins.Assign) types.Type {

	var value types.Type

	// Assign.Right is already a Type{}
	// Used in a ForIn for example
	if t, ok := assign.Right.(types.Type); ok {
		value = t
	} else {
		value = vm.Operation(assign.Right, types.ON_NOTHING)
	}

	vm.env.Set(assign.Name, value)

	return value
}

func (vm *VM) OperationMath(math ins.Math) types.Type {

	left := vm.Operation(math.Left, types.ON_NOTHING)
	right := vm.Operation(math.Right, types.ON_NOTHING)

	if math.IsComparision {
		return vm.OperationMathCompare(left, right, math.Method)
	}

	return vm.OperationMathOperation(left, right, math.Method)
}

func (vm *VM) OperationMathCompare(left, right types.Type, method string) types.Type {

	if left_n, ok := left.(*types.LiteralNumber); ok {
		if right_n, ok := right.(*types.LiteralNumber); ok {
			return vm.OperationMathCompareNumbers(left_n, right_n, method)
		}
	}

	// Fallback to Class behaviour
	l := vm.GetAsClass(left)
	r := vm.GetAsClass(right)

	return l.Compare(vm, method, r)
}

func (vm *VM) OperationMathOperation(left, right types.Type, method string) types.Type {

	if left_n, ok := left.(*types.LiteralNumber); ok {
		if right_n, ok := right.(*types.LiteralNumber); ok {
			if res, ok := vm.OperationMathOperationNumbers(left_n, right_n, method); ok {
				return res
			}
		}
	}

	// Fallback to Class behaviour
	l := vm.GetAsClass(left)
	r := vm.GetAsClass(right)

	return l.Math(vm, method, r)
}

func (vm *VM) OperationMathCompareNumbers(left, right *types.LiteralNumber, method string) types.Type {

	b := false

	switch method {
	case ">":
		b = left.Number > right.Number
	case "<":
		b = left.Number < right.Number
	case ">=":
		b = left.Number >= right.Number
	case "<=":
		b = left.Number <= right.Number
	case "==":
		b = left.Number == right.Number
	case "!=":
		b = left.Number != right.Number
	default:
		log.Panicf("OperationMathCompareNumbers() is not implementing %s", method)
	}

	res := types.LiteralBool{}
	res.Bool = b

	return &res
}

func (vm *VM) OperationMathOperationNumbers(left, right *types.LiteralNumber, method string) (types.Type, bool) {
	val := float64(0)
	found := true

	switch method {
	case "+":
		val = left.Number + right.Number
	case "-":
		val = left.Number - right.Number
	case "*":
		val = left.Number * right.Number
	case "/":
		val = left.Number / right.Number
	case "%":
		val = math.Mod(left.Number, right.Number)
	case "**":
		val = math.Pow(left.Number, right.Number)
	case "++":
		left.Number++
		val = left.Number
	case "--":
		left.Number--
		val = left.Number
	default:
		found = false
	}

	res := types.LiteralNumber{
		Number: val,
	}

	return &res, found
}

func (vm *VM) OperationLiteral(literal ins.Literal) types.Type {

	if literal.Type == "number" {
		value, err := strconv.ParseFloat(literal.Value, 64)

		if err != nil {
			log.Panicf("Can not initialize Number as %s", literal.Value)
		}

		return &types.LiteralNumber{
			Number: value,
		}
	}

	if literal.Type == "string" {
		return &types.LiteralString{
			String: literal.Value,
		}
	}

	if literal.Type == "bool" {

		value := false

		if literal.Value == "true" {
			value = true
		}

		return &types.LiteralBool{
			Bool: value,
		}
	}

	if literal.Type == "null" {
		return &types.LiteralNull{}
	}

	log.Panicf("Not able to handle Literal %s", literal)

	return &types.LiteralNull{}
}

func (vm *VM) OperationVariable(variable ins.Variable) types.Type {

	if res, ok := vm.env.Get(variable.Name); ok {
		return res
	}

	log.Print("Undefined variable, " + variable.Name)

	return vm.CreateType(&builtin.Null{})
}

func (vm *VM) ClassOperationVariable(variable ins.Variable) types.Type {

	class := vm.Classes[len(vm.Classes)-1]

	if res, ok := class.Variables[variable.Name]; ok {
		return res
	}

	log.Print("Undefined variable, " + class.Type() + "." + variable.Name)

	return vm.CreateType(&builtin.Null{})
}

func (vm *VM) OperationSet(set ins.Set) types.Type {

	left, ok := vm.env.Get(set.Name)

	if !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right, types.ON_NOTHING)

	if vm.GetType(left) != vm.GetType(value) {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, vm.GetType(left), vm.GetAsClass(value).ToString(), vm.GetType(value))
	}

	vm.env.Set(set.Name, value)

	return value
}

func (vm *VM) ClassOperationSet(set ins.Set) types.Type {

	class := vm.Classes[len(vm.Classes)-1]

	left, ok := class.Variables[set.Name]

	if !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right, types.ON_NOTHING)

	if vm.GetType(left) != vm.GetType(value) {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, vm.GetType(left), vm.GetAsClass(value).ToString(), vm.GetType(value))
	}

	class.SetVariable(set.Name, value)

	return value
}

func (vm *VM) OperationIf(i ins.If) types.Type {

	con := vm.Operation(i.Condition, types.ON_NOTHING)

	// Literal bool
	if bl, ok := con.(*types.LiteralBool); ok {
		if bl.Bool {
			return vm.Operation(i.True, types.ON_NOTHING)
		} else {
			return vm.Operation(i.False, types.ON_NOTHING)
		}
	}

	// Value of the class Bool
	if vm.GetType(con) != "Bool" {
		log.Panicf("Expecing bool in condition, %s (%s)", vm.GetAsClass(con).ToString(), vm.GetType(con))
	}

	if vm.GetAsClass(con).ToString() == "true" {
		return vm.Operation(i.True, types.ON_NOTHING)
	}

	return vm.Operation(i.False, types.ON_NOTHING)
}

func (vm *VM) OperationCall(call ins.Call) types.Type {

	params := make([]types.Type, len(call.Parameters))

	for i, param := range call.Parameters {
		params[i] = vm.Operation(param, types.ON_NOTHING)
	}

	left := vm.Operation(call.Left, types.ON_NOTHING)

	var method string

	// Optimized string
	if str, ok := left.(*types.LiteralString); ok {
		method = str.String

	// Fallbacked string behaviour
	} else {
		method = vm.GetAsClass(left).ToString()
	}

	return vm.Classes[len(vm.Classes)-1].Invoke(vm, method, params)
}

func (vm *VM) OperationDefineClass(def ins.DefineClass) types.Type {

	class := types.Class{}
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

func (vm *VM) OperationDefineMethod(def ins.DefineMethod) types.Type {

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

func (vm *VM) OperationPushClass(pushClass ins.PushClass) types.Type {

	left := vm.Operation(pushClass.Left, types.ON_NOTHING)

	var name string

	// Optimize for strings
	if str, ok := left.(*types.LiteralString); ok {
		name = str.String
	} else {
		name = vm.GetAsClass(left).ToString()
	}

	// Do not change the current pushed class
	// Continue imediately
	if name == "self" {
		return vm.Operation(pushClass.Right, types.ON_CLASS)
	}

	var class *types.Class

	value, ok := vm.env.Get(name)

	// There is no such class, use left
	if !ok {
		class = vm.GetAsClass(left)
	} else {
		class = vm.GetAsClass(value)
	}

	// Push
	vm.Classes = append(vm.Classes, class)

	res := vm.Operation(pushClass.Right, types.ON_CLASS)

	// Pop
	vm.Classes = vm.Classes[:len(vm.Classes)-1]

	return res
}

func (vm *VM) OperationMapCreate(m ins.MapCreate) types.Type {
	ma := builtin.Map{}

	params := make([]*types.Class, 0)

	for i, key := range m.Keys {
		params = append(params, vm.GetAsClass(vm.Operation(key, types.ON_NOTHING)))
		params = append(params, vm.GetAsClass(vm.Operation(m.Values[i], types.ON_NOTHING)))
	}

	ma.InitWithParams(params)

	return vm.CreateType(&ma)
}

func (vm *VM) OperationListCreate(list ins.ListCreate) types.Type {
	l := builtin.List{}

	params := make([]*types.Class, len(list.Items))

	for i, item := range list.Items {
		params[i] = vm.GetAsClass(vm.Operation(item, types.ON_NOTHING))
	}

	l.InitWithParams(params)

	return vm.CreateType(&l)
}

func (vm *VM) OperationAccessChildItem(access ins.AccessChildItem) types.Type {

	// Extract the List or Map
	item := vm.Operation(access.Item, types.ON_NOTHING)

	// Is Map
	if vm.GetType(item) == "Map" {
		return vm.OperationAccessChildItemMap(access, item)
	}

	if vm.GetType(item) != "List" {
		log.Panicf("Expected List or Map in [], got %s", vm.GetType(item))
	}

	library, ok := vm.GetAsClass(item).Extension.(*builtin.List)

	if !ok {
		log.Panic("Expected class to be of types.Type *builtin.List")
	}

	// Get position to access from the list
	position := vm.GetAsClass(vm.Operation(access.Right, types.ON_NOTHING))

	return library.ItemAt([]*types.Class{position})
}

func (vm *VM) OperationAccessChildItemMap(access ins.AccessChildItem, item types.Type) types.Type {
	library, ok := vm.GetAsClass(item).Extension.(*builtin.Map)

	if !ok {
		log.Panic("Expected class to be of types.Type *builtin.Map")
	}

	// Get position to access from the list
	position := vm.GetAsClass(vm.Operation(access.Right, types.ON_NOTHING))

	return library.Get([]*types.Class{position})
}

func (vm *VM) OperationReturn(ret ins.Return) types.Type {

	vm.ShouldReturn[len(vm.ShouldReturn)-1] = true

	return vm.Operation(ret.Statement, types.ON_NOTHING)
}

func (vm *VM) OperationInstance(instance ins.Instance) types.Type {

	in, ok := vm.env.Get(instance.Left)

	if !ok {
		log.Panicf("No such class, %s", instance.Left)
	}

	inst := vm.Clone(in)

	if len(instance.Parameters) > 0 {
		params := make([]*types.Class, 0)

		for _, node := range instance.Parameters {
			params = append(params, vm.GetAsClass(vm.Operation(node, types.ON_NOTHING)))
		}

		vm.GetAsClass(inst).Extension.InitWithParams(params)
	}

	return inst
}

//
// for before; condition; each { body }
//
func (vm *VM) OperationFor(f ins.For) types.Type {

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

		condition, is_bool := vm.GetAsClass(res).Extension.(*builtin.Bool)

		if !is_bool {
			log.Panicf("Expected bool in for, got %s", vm.GetType(res))
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
func (vm *VM) OperationForIn(f ins.For) types.Type {

	// Create variable scope
	vm.env = vm.env.Push()

	// Convert Before to an assign object
	assign, assign_ok := f.Before.(ins.Assign)

	// Get iterator object
	each := vm.Operation(f.Each, types.ON_NOTHING)

	if vm.GetType(each) != "List" {
		log.Panic("Expected List in for ... in, got %s", vm.GetType(each))
	}

	list, ok := vm.GetAsClass(each).Extension.(*builtin.List)

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
func (vm *VM) Clone(input types.Type) (out types.Type) {

	in := vm.GetAsClass(input)

	res := types.Class{}
	res.Class = in.Class
	res.Methods = in.Methods
	res.Extension, _ = in.Extension.Instance()
	res.Variables = make(map[string]types.Type)

	for name, def := range in.Variables {
		res.Variables[name] = vm.Clone(def)
	}

	return &res
}

func (vm *VM) CreateType(lib types.Lib) types.Type {
	class := types.Class{}
	class.InitWithLib(lib)

	return &class
}

func (vm VM) GetAsClass(in types.Type) *types.Class {

	if class, ok := in.(*types.Class); ok {
		return class
	}

	if lit, ok := in.(*types.LiteralNumber); ok {
		number := builtin.Number{}
		number.Value = lit.Number
		return vm.CreateType(&number).(*types.Class)
	}

	if lit, ok := in.(*types.LiteralString); ok {
		str := builtin.String{}
		str.Value = lit.String

		return vm.CreateType(&str).(*types.Class)
	}

	if lit, ok := in.(*types.LiteralBool); ok {
		bl := builtin.Bool{}
		bl.Value = lit.Bool

		return vm.CreateType(&bl).(*types.Class)
	}

	if _, ok := in.(*types.LiteralNull); ok {
		return vm.CreateType(&builtin.Null{}).(*types.Class)
	}

	log.Println("GetAsClass() defaulted to null")
	log.Println(in)
	fmt.Println(reflect.ValueOf(in).Type().String())

	return vm.CreateType(&builtin.Null{}).(*types.Class)
}

func (vm VM) GetType(in types.Type) string {
	
	if class, ok := in.(*types.Class); ok {
		return class.Type()
	}

	if _, ok := in.(*types.LiteralNumber); ok {
		return "Number"
	}

	if _, ok := in.(*types.LiteralString); ok {
		return "String"
	}

	if _, ok := in.(*types.LiteralBool); ok {
		return "Bool"
	}

	if _, ok := in.(*types.LiteralNull); ok {
		return "Null"
	}

	log.Println("GetType() could not detect type")

	return "Null"
}
