package main

import (
	"fmt"
	"log"
	"encoding/json"
)

type Type interface {
	Init(string)
	Math(string, Type) Type
	Compare(string, Type) Type
	Type() string
	toString() string
}

type ON int

const (
	ON_NOTHING ON = 1 << iota // 1
	ON_CLASS                  // 2
)

type VM struct {
	// Contains variables
	Environment *Environment

	// The current stack of methods, used to know where to define a method
	Classes []*Class

	Debug bool
}

func (vm *VM) Run(tree Block) {

	// Set empty environment
	vm.Environment = &Environment{}
	vm.Environment.Env = make(map[string]Type)

	vm.Classes = make([]*Class, 0)

	vm.Operation(tree, ON_NOTHING)
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

	if instance, ok := node.(Instance); ok {
		return vm.OperationInstance(instance)
	}

	if vm.Debug {
		fmt.Printf("Was not able to execute %s\n", node)
	}

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationBlock(block Block) (last Type) {

	vm.Environment = vm.Environment.Push()

	for _, body := range block.Body {
		last = vm.Operation(body, ON_NOTHING)
	}

	vm.Environment = vm.Environment.Pop()

	return last
}

func (vm *VM) OperationAssign(assign Assign) Type {

	value := vm.Operation(assign.Right, ON_NOTHING)

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

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationVariable(variable Variable) Type {

	if res, ok := vm.Environment.Get(variable.Name); ok {
		return res
	}

	log.Print("Undefined variable, " + variable.Name)

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) ClassOperationVariable(variable Variable) Type {

	class := vm.Classes[len(vm.Classes) - 1]

	if res, ok := class.Variables[variable.Name]; ok {
		return res
	}

	log.Print("Undefined variable, " + class.Type() + "." + variable.Name)

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationSet(set Set) Type {

	l, ok := vm.Environment.Get(set.Name)

	if !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right, ON_NOTHING)

	if l.Type() != value.Type() {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, l.Type(), value.toString(), value.Type())
	}

	vm.Environment.Set(set.Name, value)

	return value
}

func (vm *VM) ClassOperationSet(set Set) Type {

	class := vm.Classes[len(vm.Classes) - 1]

	l, ok := class.Variables[set.Name]

	if !ok {
		log.Panicf("Can not set %s, %s is undefined", set.Name, set.Name)
	}

	value := vm.Operation(set.Right, ON_NOTHING)

	if l.Type() != "Null" && l.Type() != value.Type() {
		log.Panicf("Can not set %s (type %s), to %s (type %s)", set.Name, l.Type(), value.toString(), value.Type())
	}

	class.SetVariable(set.Name, value)

	return value
}

func (vm *VM) OperationIf(i If) Type {

	con := vm.Operation(i.Condition, ON_NOTHING)

	if con.Type() != "Bool" {
		log.Panicf("Expecing bool in condition, %s (%s)", con.toString(), con.Type())
	}

	if con.toString() == "true" {
		return vm.Operation(i.True, ON_NOTHING)
	}

	return vm.Operation(i.False, ON_NOTHING)
}

func (vm *VM) OperationCall(call Call) Type {

	// Default
	bl := Bool{}
	bl.Init("false")

	// Built in method
	if call.Left == "Println" {

		for _, param := range call.Parameters {
			fmt.Println(vm.Operation(param, ON_NOTHING).toString())
		}

		bl.Init("true")
		return &bl
	}

	if call.Left == "Dump" {

		b, _ := json.MarshalIndent(vm.Environment, "", "  ")
		fmt.Println(string(b))

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

		vm.Environment = vm.Environment.Push()

		// Define variables
		for i, param := range method.Parameters {
			ass := Assign{}
			ass.Name = param.Name
			ass.Right = call.Parameters[i]

			vm.OperationAssign(ass)
		}

		body := vm.OperationBlock(method.Body)

		vm.Environment = vm.Environment.Push()

		return body
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

		if assign, ok := body.(Assign); ok {
			class.SetVariable(assign.Name, vm.Operation(assign.Right, ON_NOTHING))
			continue
		}

		// Fallback
		vm.Operation(body, ON_NOTHING)
	}

	// Pop
	vm.Classes = vm.Classes[:len(vm.Classes)-1]
	vm.Environment = vm.Environment.Pop()

	vm.Environment.Set(def.Name, &class)

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
	method.IsPublic = def.IsPublic

	vm.Classes[len(vm.Classes)-1].AddMethod(def.Name, method)

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationCallClass(callClass CallClass) Type {

	c, ok := vm.Environment.Get(callClass.Left)

	if callClass.Left == "self" {
		c = vm.Classes[len(vm.Classes) - 1]
		ok = true
	}

	if !ok {
		log.Panicf("No such class, %s", callClass.Left)
	}

	if class, ok := c.(*Class); !ok {
		log.Panicf("%s is not a class", callClass.Left)
	} else {

		// Push
		vm.Classes = append(vm.Classes, class)

		return vm.Operation(callClass.Method, ON_CLASS)

		// Pop
		vm.Classes = vm.Classes[:len(vm.Classes)-1]

	}

	// Default
	bl := Bool{}
	bl.Init("false")

	return &bl
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

func (vm *VM) Clone(in Type) (out Type) {

	if class, ok := in.(*Class); ok {
		res := Class{}
		res.Methods = class.Methods
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


	out.Init(in.toString())
	return out
}