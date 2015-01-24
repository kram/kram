package main

type Node interface{}

type Nil struct{}

type Block struct {
	Block bool
	Body  []Node
}

type Assign struct {
	Assign bool
	Name   string
	Right  Node
}

type Set struct {
	Set   bool
	Name  string
	Right Node
}

type Literal struct {
	Literal bool
	Type    string
	Value   string
}

type Variable struct {
	Variable bool
	Name     string
}

type Math struct {
	Math          bool
	Method        string
	IsComparision bool
	Left          Node
	Right         Node
}

type If struct {
	If        bool
	Condition Node
	True      Block
	False     Block
}

type Condition struct {
	Condition string // && || > < >= <=
	Left      Node
	Right     Node
}

type Call struct {
	Call       bool
	Left       Node
	Parameters []Node
}

type DefineClass struct {
	DefineClass bool
	Name        string
	Body        Block
}

type DefineMethod struct {
	DefineMethod bool
	Name         string
	Body         Block
	IsStatic     bool
	IsPublic     bool
	Parameters   []Parameter
}

type Parameter struct {
	Name string
}

type Instance struct {
	Instance bool
	Left     string
}

type CreateList struct {
	CreateList bool
	Items      []Node
}

type Return struct {
	Return    bool
	Statement Node
}

type PushClass struct {
	PushClass bool
	Left      Node
	Right     Node
}