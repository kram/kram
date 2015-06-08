// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

// This is a complete list of all the instructions that the VM supports
// All of these creates the Abstract Syntax Tree with a Block{} as the root object.
package instructions

type Node interface{}

type Nil struct {
	Nil bool
}

// A list of other operations
// Scope indicates if a new variable scope should be created, eg when in a for-loop or a method.
type Block struct {
	Block bool
	Body  []Node
	Scope bool
}

// Creates a new variable (named Name) and sets the value to Right
type Assign struct {
	Assign bool
	Name   string
	Right  Node
}

// Sets an existing variable (Name) to the valie of Right
type Set struct {
	Set   bool
	Name  string
	Right Node
}

// The value of a literal
// Supported Type's are "number" / "string" / "bool" / "null"
// and Value is a string representaion of that value (from the sourcecode)
type Literal struct {
	Literal bool
	Type    string
	Value   string
}

// Representation of a variable that is already existing
type Variable struct {
	Variable bool
	Name     string
}

// Mathemathical operations
// Both Left and Right are conditional, but at least one of them will always exist.
// Method could be + - ... and many other
type Math struct {
	Math          bool
	Method        string
	IsComparision bool
	Left          Node
	Right         Node
}

// The basics of If-cases, Condition is usually one or more (nested) Math-operations
type If struct {
	If        bool
	Condition Node
	True      Block
	False     Block
}

// Call (execute) the method named Left with Arguments as the arguments
type Call struct {
	Call      bool
	Left      Node
	Arguments []Argument
}

// Define a new class. Body is a block with DefineMethod's at the top level
type DefineClass struct {
	DefineClass bool
	Name        string
	Body        Block
}

// Defines a new method in the parent DefineClass
type DefineMethod struct {
	DefineMethod bool
	Name         string
	Body         Block
	IsStatic     bool
	Parameters   []Parameter
}

// Paraments are the receivers in a method
type Parameter struct {
	Name       string
	Default    Node
	HasDefault bool
}

// Arguments is the data thay you send to a method
type Argument struct {
	Argument bool
	IsNamed  bool
	Name     string
	Value    Node
}

// Instance creates a new instance of a class/type. This originates from the "new"-keyword in the sourcecode
type Instance struct {
	Instance  bool
	Left      string
	Arguments []Argument
}

// Create a new Map
type MapCreate struct {
	MapCreate bool
	Keys      []Node
	Values    []Node
}

// Creates a new List
type ListCreate struct {
	ListCreate bool
	Items      []Node
}

// Access a child item of a Map or a List
type AccessChildItem struct {
	AccessChildItem bool
	Item            Node
	Right           Node
}

// Abort method execution and return Statement to the parent
type Return struct {
	Return    bool
	Statement Node
}

// Pushes the class (Left) and continues to execute Right
type PushClass struct {
	PushClass bool
	Left      Node
	Right     Node
}

// There is two types of for loopes, IsForIn (for var a in range) and normal ones (for var a = 0; a < 5; a++)
// IsForIn uses Before (the Assign) and Each (the variable to iterate over). While the other type uses all three (from left to right)
type For struct {
	For       bool
	IsForIn   bool
	Before    Node
	Condition Node
	Each      Node
	Body      Block
}
