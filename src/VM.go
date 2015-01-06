package main

import (
	"log"
	"fmt"
)

type Type interface {
	Init(string)
	Math(string, Type) Type
	Type() string
	toString() string
}

type VM struct {
	// Contains variables
	Environment map[string]Type

	// The current stack of methods, used to know where to define a method
	Classes []*Class
}

func (vm *VM) Run(tree Block) {

	// Set empty environment
	vm.Environment = make(map[string]Type)
	vm.Classes = make([]*Class, 0)

	vm.Operation(tree)
}

func (vm *VM) Operation(node Node) Type {

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
		return vm.OperationVariable(variable)
	}

	if set, ok := node.(Set); ok {
		return vm.OperationSet(set)
	}

	if i, ok := node.(If); ok {
		return vm.OperationIf(i)
	}

	if block, ok := node.(Block); ok {
		return vm.OperationBlock(block)
	}

	if call, ok := node.(Call); ok {
		return vm.OperationCall(call)
	}

	if defineClass, ok := node.(DefineClass); ok {
		return vm.OperationDefineClass(defineClass)
	}

	if defineMethod, ok := node.(DefineMethod); ok {
		return vm.OperationDefineMethod(defineMethod)
	}

	log.Panicf("Was not able to expecute %s", node)

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationBlock(block Block) (last Type) {

	for _, body := range block.Body {
		last = vm.Operation(body)
	}

	return last
}

func (vm *VM) OperationAssign(assign Assign) Type {

	value := vm.Operation(assign.Right)

	vm.Environment[assign.Name] = value

	return value
}

func (vm *VM) OperationMath(math Math) Type {

	left := vm.Operation(math.Left)
	right := vm.Operation(math.Right)

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

	log.Panicf("Not able to handle Literal %s", literal)

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationVariable(variable Variable) Type {

	if _, ok := vm.Environment[variable.Name]; ok {
		return vm.Environment[variable.Name]
	}

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationSet(set Set) Type {

	if _, ok := vm.Environment[set.Name]; !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right)

	if vm.Environment[set.Name].Type() != value.Type() {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, vm.Environment[set.Name].Type(), value.toString(), value.Type())
	}

	vm.Environment[set.Name] = value

	return vm.Environment[set.Name]
}

func (vm *VM) OperationIf(i If) Type {

	con := vm.Operation(i.Condition)

	if con.Type() != "Bool" {
		log.Panicf("Expecing bool in condition, %s (%s)", con.toString(), con.Type())
	}

	if con.toString() == "true" {
		return vm.Operation(i.True)
	}
	
	return vm.Operation(i.False)
}

func (vm *VM) OperationCall(call Call) Type {

	params := make([]Type, 0)

	for _, param := range call.Parameters {
		params = append(params, vm.Operation(param))
	}


	// Built in method
	if call.Left == "Println" {
		for _, p := range params {
			fmt.Println(p.toString())
		}

		bl := Bool{}
		bl.Init("true")
		return &bl
	}

	fmt.Println("Call to undefined function %s", call.Left)

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationDefineClass(def DefineClass) Type {

	class := Class{}
	class.Init(def.Name)

	// Push
	vm.Classes = append(vm.Classes, &class)

	vm.OperationBlock(def.Body)

	// Pop
	vm.Classes = vm.Classes[:len(vm.Classes) - 1]

	vm.Environment[def.Name] = &class

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationDefineMethod(def DefineMethod) Type {

	if len(vm.Classes) == 0 {
		log.Panic("Unable to define method, not in class")
	}

	method := Method{}
	method.Parameters = make([]string, 0)
	method.Body = def.Body

	vm.Classes[len(vm.Classes) - 1].AddMethod(def.Name, method)

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}