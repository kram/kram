package main

import (
	"./Instructions"
	"./Types"
	"encoding/json"
	"fmt"
	"log"
)

type ON int

const (
	ON_NOTHING ON = 1 << iota // 1
	ON_CLASS                  // 2
)

type VM struct {
	// Contains variables
	Environment *Environment

	// The current stack of methods, used to know where to define a method
	Classes []*types.Class

	Debug bool
}

func (vm *VM) Run(tree instructions.Block) {

	// Set empty environment
	vm.Environment = &Environment{}
	vm.Environment.Env = make(map[string]types.Type)

	vm.Classes = make([]*types.Class, 0)

	vm.Operation(tree, ON_NOTHING)
}

func (vm *VM) Operation(node instructions.Node, on ON) types.Type {

	if assign, ok := node.(instructions.Assign); ok {
		return vm.OperationAssign(assign)
	}

	if math, ok := node.(instructions.Math); ok {
		return vm.OperationMath(math)
	}

	if literal, ok := node.(instructions.Literal); ok {
		return vm.OperationLiteral(literal)
	}

	if variable, ok := node.(instructions.Variable); ok {

		if on == ON_CLASS {
			return vm.ClassOperationVariable(variable)
		}

		return vm.OperationVariable(variable)
	}

	if set, ok := node.(instructions.Set); ok {

		if on == ON_CLASS {
			return vm.ClassOperationSet(set)
		}

		return vm.OperationSet(set)
	}

	if i, ok := node.(instructions.If); ok {
		return vm.OperationIf(i)
	}

	if block, ok := node.(instructions.Block); ok {
		return vm.OperationBlock(block)
	}

	if call, ok := node.(instructions.Call); ok {
		return vm.OperationCall(call)
	}

	if callClass, ok := node.(instructions.CallClass); ok {
		return vm.OperationCallClass(callClass)
	}

	if defineClass, ok := node.(instructions.DefineClass); ok {
		return vm.OperationDefineClass(defineClass)
	}

	if defineMethod, ok := node.(instructions.DefineMethod); ok {
		return vm.OperationDefineMethod(defineMethod)
	}

	if instance, ok := node.(instructions.Instance); ok {
		return vm.OperationInstance(instance)
	}

	if vm.Debug {
		fmt.Printf("Was not able to execute %s\n", node)
	}

	// Default
	bl := types.Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationBlock(block instructions.Block) (last types.Type) {

	vm.Environment = vm.Environment.Push()

	for _, body := range block.Body {
		last = vm.Operation(body, ON_NOTHING)
	}

	vm.Environment = vm.Environment.Pop()

	return last
}

func (vm *VM) OperationAssign(assign instructions.Assign) types.Type {

	value := vm.Operation(assign.Right, ON_NOTHING)

	vm.Environment.Set(assign.Name, value)

	return value
}

func (vm *VM) OperationMath(math instructions.Math) types.Type {

	left := vm.Operation(math.Left, ON_NOTHING)
	right := vm.Operation(math.Right, ON_NOTHING)

	if math.IsComparision {
		return left.Compare(math.Method, right)
	}

	return left.Math(math.Method, right)
}

func (vm *VM) OperationLiteral(literal instructions.Literal) types.Type {

	if literal.Type == "number" {
		number := types.Number{}
		number.Init(literal.Value)
		return &number
	}

	if literal.Type == "string" {
		str := types.String{}
		str.Init(literal.Value)
		return &str
	}

	if literal.Type == "bool" {
		bl := types.Bool{}
		bl.Init(literal.Value)
		return &bl
	}

	if literal.Type == "null" {
		null := types.Null{}
		return &null
	}

	log.Panicf("Not able to handle Literal %s", literal)

	// Default
	bl := types.Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationVariable(variable instructions.Variable) types.Type {

	if res, ok := vm.Environment.Get(variable.Name); ok {
		return res
	}

	log.Print("Undefined variable, " + variable.Name)

	// Default
	bl := types.Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) ClassOperationVariable(variable instructions.Variable) types.Type {

	class := vm.Classes[len(vm.Classes)-1]

	if res, ok := class.Variables[variable.Name]; ok {
		return res
	}

	log.Print("Undefined variable, " + class.Type() + "." + variable.Name)

	// Default
	bl := types.Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationSet(set instructions.Set) types.Type {

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

func (vm *VM) ClassOperationSet(set instructions.Set) types.Type {

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

func (vm *VM) OperationIf(i instructions.If) types.Type {

	con := vm.Operation(i.Condition, ON_NOTHING)

	if con.Type() != "Bool" {
		log.Panicf("Expecing bool in condition, %s (%s)", con.ToString(), con.Type())
	}

	if con.ToString() == "true" {
		return vm.Operation(i.True, ON_NOTHING)
	}

	return vm.Operation(i.False, ON_NOTHING)
}

func (vm *VM) OperationCall(call instructions.Call) types.Type {

	// Default
	bl := types.Bool{}
	bl.Init("false")

	// Built in method
	if call.Left == "Println" {

		for _, param := range call.Parameters {
			fmt.Println(vm.Operation(param, ON_NOTHING).ToString())
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
			fmt.Printf("Can not call %s.%s() (%d parameters) with %d parameters\n", vm.Classes[len(vm.Classes)-1].ToString(), call.Left, len(method.Parameters), len(call.Parameters))

			return &bl
		}

		vm.Environment = vm.Environment.Push()

		// Define variables
		for i, param := range method.Parameters {
			ass := instructions.Assign{}
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

func (vm *VM) OperationDefineClass(def instructions.DefineClass) types.Type {

	class := types.Class{}
	class.Init(def.Name)

	// Push
	vm.Classes = append(vm.Classes, &class)
	vm.Environment = vm.Environment.Push()

	for _, body := range def.Body.Body {

		if assign, ok := body.(instructions.Assign); ok {
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
	bl := types.Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationDefineMethod(def instructions.DefineMethod) types.Type {

	if len(vm.Classes) == 0 {
		log.Panic("Unable to define method, not in class")
	}

	method := types.Method{}
	method.Parameters = def.Parameters
	method.Body = def.Body
	method.IsStatic = def.IsStatic
	method.IsPublic = def.IsPublic

	vm.Classes[len(vm.Classes)-1].AddMethod(def.Name, method)

	// Default
	bl := types.Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationCallClass(callClass instructions.CallClass) types.Type {

	c, ok := vm.Environment.Get(callClass.Left)

	if callClass.Left == "self" {
		c = vm.Classes[len(vm.Classes)-1]
		ok = true
	}

	if !ok {
		log.Panicf("No such class, %s", callClass.Left)
	}

	if class, ok := c.(*types.Class); !ok {
		log.Panicf("%s is not a class", callClass.Left)
	} else {

		// Push
		vm.Classes = append(vm.Classes, class)

		return vm.Operation(callClass.Method, ON_CLASS)

		// Pop
		vm.Classes = vm.Classes[:len(vm.Classes)-1]

	}

	// Default
	bl := types.Bool{}
	bl.Init("false")

	return &bl
}

func (vm *VM) OperationInstance(instance instructions.Instance) types.Type {

	in, ok := vm.Environment.Get(instance.Left)

	if !ok {
		log.Panicf("No such class, %s", instance.Left)
	}

	class, ok := in.(*types.Class)

	if !ok {
		log.Panicf("%s is not a class", instance.Left)
	}

	return vm.Clone(class)
}

func (vm *VM) Clone(in types.Type) (out types.Type) {

	if class, ok := in.(*types.Class); ok {
		res := types.Class{}
		res.Methods = class.Methods
		res.Variables = make(map[string]types.Type)

		for name, def := range class.Variables {
			res.Variables[name] = vm.Clone(def)
		}

		return &res
	}

	if _, ok := in.(*types.Number); ok {
		out = &types.Number{}
	}

	if _, ok := in.(*types.Null); ok {
		out = &types.Null{}
	}

	if _, ok := in.(*types.Bool); ok {
		out = &types.Bool{}
	}

	if _, ok := in.(*types.String); ok {
		out = &types.String{}
	}

	out.Init(in.ToString())
	return out
}
