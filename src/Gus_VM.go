package main

import (
	"fmt"
	"log"
)

type ON int

const (
	ON_NOTHING    ON = 1 << iota // 1
	ON_CLASS                     // 2
	ON_METHOD_BODY               // 4
	ON_FOR_PART                  // 8
)

type VM struct {
	// Contains variables
	Environment *Environment

	// The current stack of methods, used to know where to define a method
	Classes []*Class

	Debug bool

	ShouldReturn []bool
}

func (vm *VM) Run(tree Block) {

	// Set empty environment
	vm.Environment = &Environment{}
	vm.Environment.Env = make(map[string]Type)
	vm.ShouldReturn = make([]bool, 0)

	vm.Classes = make([]*Class, 0)

	vm.Libraries()

	vm.Operation(tree, ON_NOTHING)
}

func (vm *VM) Libraries() {

	libs := make([]Lib, 0)

	libs = append(libs, &Library_IO{})
	libs = append(libs, &Library_List{})
	libs = append(libs, &Library_String{})
	libs = append(libs, &Library_File{})
	libs = append(libs, &Library_Map{})

	for _, lib := range libs {

		instance, name := lib.Instance()

		class := Class{}
		class.Init(name)
		class.Extension = instance

		vm.Environment.Set(name, &class)
	}
}

func (vm *VM) Operation(node Node, on ON) Type {

	if assign, ok := node.(Assign); ok {
		return vm.OperationAssign(assign)
	}

	if math, ok := node.(Math); ok {
		return vm.OperationMath(math)
	}

	if literal, ok := node.(Literal); ok {
		return vm.OperationLiteral(literal)
	}

	if variable, ok := node.(Variable); ok {

		if on == ON_CLASS {
			return vm.ClassOperationVariable(variable)
		}

		return vm.OperationVariable(variable)
	}

	if set, ok := node.(Set); ok {

		if on == ON_CLASS {
			return vm.ClassOperationSet(set)
		}

		return vm.OperationSet(set)
	}

	if i, ok := node.(If); ok {
		return vm.OperationIf(i)
	}

	if block, ok := node.(Block); ok {
		return vm.OperationBlock(block, on)
	}

	if call, ok := node.(Call); ok {
		return vm.OperationCall(call)
	}

	if pushClass, ok := node.(PushClass); ok {
		return vm.OperationPushClass(pushClass)
	}

	if defineClass, ok := node.(DefineClass); ok {
		return vm.OperationDefineClass(defineClass)
	}

	if defineMethod, ok := node.(DefineMethod); ok {
		return vm.OperationDefineMethod(defineMethod)
	}

	if instance, ok := node.(Instance); ok {
		return vm.OperationInstance(instance)
	}

	if list, ok := node.(ListCreate); ok {
		return vm.OperationListCreate(list)
	}

	if access, ok := node.(AccessChildItem); ok {
		return vm.OperationAccessChildItem(access)
	}

	if m, ok := node.(MapCreate); ok {
		return vm.OperationMapCreate(m)
	}

	if ret, ok := node.(Return); ok {
		return vm.OperationReturn(ret)
	}

	if f, ok := node.(For); ok {
		return vm.OperationFor(f)
	}

	return &Null{}
}

func (vm *VM) OperationBlock(block Block, on ON) (last Type) {

	// Create new scope
	if on != ON_FOR_PART && block.Scope == true {
		vm.Environment = vm.Environment.Push()
	}

	if on == ON_METHOD_BODY {
		vm.ShouldReturn = append(vm.ShouldReturn, false)
	}

	for _, body := range block.Body {
		last = vm.Operation(body, ON_NOTHING)

		// Return statement
		if len(vm.ShouldReturn) > 0 && vm.ShouldReturn[len(vm.ShouldReturn)-1] {
			break
		}
	}

	// Pop
	if on == ON_METHOD_BODY {
		vm.ShouldReturn = vm.ShouldReturn[:len(vm.ShouldReturn)-1]
	}

	// Restore scope
	if on != ON_FOR_PART && block.Scope == true {
		vm.Environment = vm.Environment.Pop()
	}

	return last
}

func (vm *VM) OperationAssign(assign Assign) Type {

	var value Type

	// Assign.Right is already a Type{}
	// Used in a ForIn for example
	if t, ok := assign.Right.(Type); ok {
		value = t
	} else {
		value = vm.Operation(assign.Right, ON_NOTHING)	
	}

	vm.Environment.Set(assign.Name, value)

	return value
}

func (vm *VM) OperationMath(math Math) Type {

	left := vm.Operation(math.Left, ON_NOTHING)
	right := vm.Operation(math.Right, ON_NOTHING)

	if math.IsComparision {
		return left.Compare(math.Method, right)
	}

	return left.Math(math.Method, right)
}

func (vm *VM) OperationLiteral(literal Literal) Type {

	if literal.Type == "number" {
		number := Number{}
		number.Init(literal.Value)
		return &number
	}

	if literal.Type == "string" {
		str := String{}
		str.Init(literal.Value)
		return &str
	}

	if literal.Type == "bool" {
		bl := Bool{}
		bl.Init(literal.Value)
		return &bl
	}

	if literal.Type == "null" {
		null := Null{}
		return &null
	}

	log.Panicf("Not able to handle Literal %s", literal)

	return &Null{}
}

func (vm *VM) OperationVariable(variable Variable) Type {

	if res, ok := vm.Environment.Get(variable.Name); ok {
		return res
	}

	log.Print("Undefined variable, " + variable.Name)

	return &Null{}
}

func (vm *VM) ClassOperationVariable(variable Variable) Type {

	class := vm.Classes[len(vm.Classes)-1]

	if res, ok := class.Variables[variable.Name]; ok {
		return res
	}

	log.Print("Undefined variable, " + class.Type() + "." + variable.Name)

	return &Null{}
}

func (vm *VM) OperationSet(set Set) Type {

	l, ok := vm.Environment.Get(set.Name)

	if !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right, ON_NOTHING)

	if l.Type() != value.Type() {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, l.Type(), value.ToString(), value.Type())
	}

	vm.Environment.Set(set.Name, value)

	return value
}

func (vm *VM) ClassOperationSet(set Set) Type {

	class := vm.Classes[len(vm.Classes)-1]

	l, ok := class.Variables[set.Name]

	if !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right, ON_NOTHING)

	if l.Type() != "Null" && l.Type() != value.Type() {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, l.Type(), value.ToString(), value.Type())
	}

	class.SetVariable(set.Name, value)

	return value
}

func (vm *VM) OperationIf(i If) Type {

	con := vm.Operation(i.Condition, ON_NOTHING)

	if con.Type() != "Bool" {
		log.Panicf("Expecing bool in condition, %s (%s)", con.ToString(), con.Type())
	}

	if con.ToString() == "true" {
		return vm.Operation(i.True, ON_NOTHING)
	}

	return vm.Operation(i.False, ON_NOTHING)
}

func (vm *VM) OperationCall(call Call) Type {

	// Default
	bl := Bool{}
	bl.Init("false")

	// Calling a method
	if len(vm.Classes) >= 0 {
		return vm.Classes[len(vm.Classes)-1].Invoke(vm, vm.Operation(call.Left, ON_NOTHING).ToString(), call.Parameters)
	}

	fmt.Printf("Call to undefined function %s\n", call.Left)

	return &bl
}

func (vm *VM) OperationDefineClass(def DefineClass) Type {

	class := Class{}
	class.Init(def.Name)

	// Push
	vm.Classes = append(vm.Classes, &class)
	vm.Environment = vm.Environment.Push()

	for _, body := range def.Body.Body {

		/*if assign, ok := body.(Assign); ok {
			class.SetVariable(assign.Name, vm.Operation(assign.Right, ON_NOTHING))
			continue
		}*/

		// Fallback
		vm.Operation(body, ON_NOTHING)
	}

	// Pop
	vm.Classes = vm.Classes[:len(vm.Classes)-1]
	vm.Environment = vm.Environment.Pop()

	vm.Environment.Set(def.Name, &class)

	return &Null{}
}

func (vm *VM) OperationDefineMethod(def DefineMethod) Type {

	if len(vm.Classes) == 0 {
		log.Panic("Unable to define method, not in class")
	}

	method := Method{}
	method.Parameters = def.Parameters
	method.Body = def.Body
	method.IsStatic = def.IsStatic
	method.IsPublic = def.IsPublic

	vm.Classes[len(vm.Classes)-1].AddMethod(def.Name, method)

	return &Null{}
}

func (vm *VM) OperationPushClass(pushClass PushClass) Type {

	name := vm.Operation(pushClass.Left, ON_NOTHING).ToString()

	c, ok := vm.Environment.Get(name)

	if name == "self" {
		c = vm.Classes[len(vm.Classes)-1]
		ok = true
	}

	if !ok {
		log.Panicf("No such class, %s", name)
	}

	if class, ok := c.(*Class); !ok {
		log.Panicf("%s is not a class", name)
	} else {

		// Push
		vm.Classes = append(vm.Classes, class)

		res := vm.Operation(pushClass.Right, ON_CLASS)

		// Pop
		vm.Classes = vm.Classes[:len(vm.Classes)-1]

		return res
	}

	return &Null{}
}

func (vm *VM) OperationMapCreate(m MapCreate) Type {
	mapinstance := vm.OperationInstance(Instance{
		Left: "Map",
	})

	class, ok := mapinstance.(*Class)

	if !ok {
		log.Panicf("Expected Map, got something else.")
	}

	params := make([]Node, 0)

	for i, key := range m.Keys {
		params = append(params, key)
		params = append(params, m.Values[i])
	}

	class.Invoke(vm, "Init", params)

	return mapinstance
}

func (vm *VM) OperationListCreate(list ListCreate) Type {
	l := vm.OperationInstance(Instance{
		Left: "List",
	})

	class, ok := l.(*Class)

	if !ok {
		log.Panicf("Expected List, got something else.")
	}

	class.Invoke(vm, "Init", list.Items)

	return l
}

func (vm *VM) OperationAccessChildItem(access AccessChildItem) Type {

	// Extract the List or Map
	item := vm.Operation(access.Item, ON_NOTHING)

	// Is Map
	if item.Type() == "Map" {
		return vm.OperationAccessChildItemMap(access, item)
	}

	if (item.Type() != "List") {
		log.Panicf("Expected List or Map in [], got %s", item.Type())
	}

	class, ok := item.(*Class)

	if !ok {
		log.Panic("Expected object to be of type *Class")
	}

	library, ok := class.Extension.(*Library_List)

	if !ok {
		log.Panic("Expected class to be of type *Library_List")
	}

	// Get position to access from the list
	position := vm.Operation(access.Right, ON_NOTHING)

	return library.ItemAt([]Type{position})
}

func (vm *VM) OperationAccessChildItemMap(access AccessChildItem, item Type) Type {
	class, ok := item.(*Class)

	if !ok {
		log.Panic("Expected object to be of type *Class")
	}

	library, ok := class.Extension.(*Library_Map)

	if !ok {
		log.Panic("Expected class to be of type *Library_List")
	}

	// Get position to access from the list
	position := vm.Operation(access.Right, ON_NOTHING)

	return library.Get([]Type{position})
}

func (vm *VM) OperationReturn(ret Return) Type {

	vm.ShouldReturn[len(vm.ShouldReturn)-1] = true

	return vm.Operation(ret.Statement, ON_NOTHING)
}

func (vm *VM) OperationInstance(instance Instance) Type {

	in, ok := vm.Environment.Get(instance.Left)

	if !ok {
		log.Panicf("No such class, %s", instance.Left)
	}

	class, ok := in.(*Class)

	if !ok {
		log.Panicf("%s is not a class", instance.Left)
	}

	return vm.Clone(class)
}

//
// for before; condition; each { body }
//
func (vm *VM) OperationFor(f For) Type {

	if f.IsForIn {
		return vm.OperationForIn(f)
	}

	// Create variable scope
	vm.Environment = vm.Environment.Push()

	// Execute before part
	vm.Operation(f.Before, ON_FOR_PART)

	for {

		// Test condition
		res := vm.Operation(f.Condition, ON_FOR_PART)

		condition, is_bool := res.(*Bool)

		if !is_bool {
			log.Panicf("Expected bool in for, got %s", res.Type())
		}

		if !condition.Value {
			break
		}

		// Execute body
		vm.Operation(f.Body, ON_NOTHING)

		// Execute part after each run
		vm.Operation(f.Each, ON_FOR_PART)
	}

	// Restore scope
	vm.Environment = vm.Environment.Pop()

	return &Null{}
}

// for var item in 1..2
// for var item in ["first", "second"]
// for var item in list
func (vm *VM) OperationForIn(f For) Type {

	// Create variable scope
	vm.Environment = vm.Environment.Push()
 
	// Convert Before to an assign object
	assign, assign_ok := f.Before.(Assign)

	// Get iterator object
	each := vm.Operation(f.Each, ON_NOTHING)

	if (each.Type() != "List") {
		log.Panic("Expected List in for ... in, got %s", each.Type())
	}

	class, ok := each.(*Class)

	if !ok {
		log.Panic("Expected object to be of type *Class")
	}

	list, ok := class.Extension.(*Library_List)

	if !ok {
		log.Panic("Expected class to be of type *Library_List")
	}

	length := list.Length()

	for key := 0; key < length; key++ {

		// Update variable
		if assign_ok {
			item := list.ItemAtPosition(key)
			assign.Right = item
			vm.Operation(assign, ON_NOTHING)
		}

		// Run body
		vm.Operation(f.Body, ON_NOTHING)
	}

	// Restore scope
	vm.Environment = vm.Environment.Pop()

	return &Null{}
}

//
// Clones a type, returns the new one
//
func (vm *VM) Clone(in Type) (out Type) {

	if class, ok := in.(*Class); ok {
		res := Class{}
		res.Class = class.Class
		res.Methods = class.Methods
		res.Extension, _ = class.Extension.Instance()
		res.Variables = make(map[string]Type)

		for name, def := range class.Variables {
			res.Variables[name] = vm.Clone(def)
		}

		return &res
	}

	if _, ok := in.(*Number); ok {
		out = &Number{}
	}

	if _, ok := in.(*Null); ok {
		out = &Null{}
	}

	if _, ok := in.(*Bool); ok {
		out = &Bool{}
	}

	if _, ok := in.(*String); ok {
		out = &String{}
	}

	out.Init(in.ToString())
	return out
}
