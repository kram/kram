package main

import (
	"fmt"
	"log"
)

type Type interface {
	Init(string)
	Math(string, Type) Type
	Type() string
	toString() string
}

type VM struct {
	Environment map[string]Type
}

func (vm *VM) Run(tree Block) {

	vm.Environment = make(map[string]Type)

	for _, body := range tree.Body {
		vm.Operation(body)
	}
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

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationAssign(assign Assign) Type {

	value := vm.Operation(assign.Right)

	vm.Environment[assign.Name] = value

	return value
}

func (vm *VM) OperationMath(math Math) Type {

	left := vm.Operation(math.Left)
	right := vm.Operation(math.Right)

	fmt.Println(left, math.Method, right)

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
