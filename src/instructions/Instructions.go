package instructions

type Node interface{}

type Nil struct {
	Nil bool
}

type Block struct {
	Block bool
	Body  []Node
	Scope bool
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

type MapCreate struct {
	MapCreate bool
	Keys      []Node
	Values    []Node
}

type ListCreate struct {
	ListCreate bool
	Items      []Node
}

type AccessChildItem struct {
	AccessChildItem bool
	Item            Node
	Right           Node
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

type For struct {
	For       bool
	IsForIn   bool
	Before    Node
	Condition Node
	Each      Node
	Body      Block
}

type Iterate struct {
	Iterate bool
	Name    string
	Object  Node
}
