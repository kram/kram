package main

import (
	"fmt"
	"log"
)

type Type interface {
	Init(string)
	Math(string, Type) Type
	Compare(string, Type) Type
	Type() string
	toString() string
}

type VM struct {
	// Contains variables
	Environment map[string]Type

	// The current stack of methods, used to know where to define a method
	Classes []*Class

	Debug bool
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

	if callClass, ok := node.(CallClass); ok {
		return vm.OperationCallClass(callClass)
	}

	if defineClass, ok := node.(DefineClass); ok {
		return vm.OperationDefineClass(defineClass)
	}

	if defineMethod, ok := node.(DefineMethod); ok {
		return vm.OperationDefineMethod(defineMethod)
	}

	if vm.Debug {
		fmt.Printf("Was not able to expecute %s\n", node)
	}

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

	// Default
	bl := Bool{}
	bl.Init("false")

	// Built in method
	if call.Left == "Println" {

		for _, param := range call.Parameters {
			fmt.Println(vm.Operation(param).toString())
		}

		bl.Init("true")
		return &bl
	}

	// Calling a method
	if len(vm.Classes) >= 0 {

		method := vm.Classes[len(vm.Classes)-1].Methods[call.Left]

		if len(method.Parameters) != len(call.Parameters) {
			fmt.Printf("Can not call %s.%s() (%d parameters) with %d parameters\n", vm.Classes[len(vm.Classes)-1].toString(), call.Left, len(method.Parameters), len(call.Parameters))

			return &bl
		}

		// Define variables
		for i, param := range method.Parameters {
			ass := Assign{}
			ass.Name = param.Name
			ass.Right = call.Parameters[i]

			vm.OperationAssign(ass)
		}

		return vm.OperationBlock(method.Body)
	}

	fmt.Printf("Call to undefined function %s\n", call.Left)

	return &bl
}

func (vm *VM) OperationDefineClass(def DefineClass) Type {

	class := Class{}
	class.Init(def.Name)

	// Push
	vm.Classes = append(vm.Classes, &class)

	vm.OperationBlock(def.Body)

	// Pop
	vm.Classes = vm.Classes[:len(vm.Classes)-1]

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
	method.Parameters = def.Parameters
	method.Body = def.Body
	method.IsStatic = def.IsStatic

	vm.Classes[len(vm.Classes)-1].AddMethod(def.Name, method)

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationCallClass(callClass CallClass) Type {

	if _, ok := vm.Environment[callClass.Left]; !ok {
		log.Panicf("No such class, %s", callClass.Left)
	}

	c := vm.Environment[callClass.Left]

	if class, ok := c.(*Class); !ok {
		log.Panicf("%s is not a class", callClass.Left)
	} else {

		// Push
		vm.Classes = append(vm.Classes, class)

		return vm.Operation(callClass.Method)

		// Pop
		vm.Classes = vm.Classes[:len(vm.Classes)-1]

	}

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) Echo(t Type) Type {
	return t
}
