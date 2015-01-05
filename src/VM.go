package main

import (
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

	// Set empty environment
	vm.Environment = make(map[string]Type)

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