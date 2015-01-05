package main

import (
	"fmt"
	"reflect"
)

type Type interface {
	Init(string)
	Math(string, Type) Type
	toString() string
}

type VM struct {
	Environment map[string]Type
}

func (vm *VM) Run(tree Block) {

	vm.Environment = make(map[string]Type)

	for _, body := range tree.Body {
		b := body
		vm.Operation(b)
	}
}

func (vm *VM) Operation(node Node) Type {
	fmt.Println("Operation()")
	fmt.Println(reflect.TypeOf(node).String())

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

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationAssign(assign Assign) Type {
	fmt.Println("OperationAssign()")

	value := vm.Operation(assign.Right)

	vm.Environment[assign.Name] = value

	return value
}

func (vm *VM) OperationMath(math Math) Type {
	fmt.Println("OperationMath()")

	left := vm.Operation(math.Left)
	right := vm.Operation(math.Right)

	fmt.Println(left, math.Method, right)

	return left.Math(math.Method, right)
}

func (vm *VM) OperationLiteral(literal Literal) Type {
	fmt.Println("OperationLiteral()")

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
	fmt.Println("OperationVariable()")

	if _, ok := vm.Environment[variable.Name]; ok {
		return vm.Environment[variable.Name]
	}

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}
