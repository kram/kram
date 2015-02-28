package main

import (
	"log"
)

// --------------- Symbols

type Symbol struct {
	Function     SymbolReturn
	CaseFunction SymbolCaseReturn
	Importance   int
	IsStatement  bool
}

type SymbolReturn func(Expecting) Node
type SymbolCaseReturn func(Expecting) Symbol

// --------------- Symbols

// --------------- Constants

type Expecting int

const (
	EXPECTING_NOTHING     Expecting = 1 << iota // 1
	EXPECTING_CLASS_BODY                        // 2
	EXPECTING_IF_BODY                           // 4
	EXPECTING_METHOD_BODY                       // 8
	EXPECTING_EXPRESSION                        // 16
	EXPECTING_FOR_PART                          // 32
)

// --------------- Constants

// --------------- Stack

type Stack struct {
	Items   *[]Node
	Parents []*[]Node
}

func (stack *Stack) Pop() {
	if len(stack.Parents) == 0 {
		items := make([]Node, 0)
		stack.Items = &items
		return
	}

	stack.Items = stack.Parents[len(stack.Parents)-1]
	stack.Parents = stack.Parents[:len(stack.Parents)-1]
}

func (stack *Stack) Push() {
	stack.Parents = append(stack.Parents, stack.Items)

	items := make([]Node, 0)
	stack.Items = &items
}

func (stack *Stack) Add(node Node) {
	items := *stack.Items
	items = append(items, node)

	stack.Items = &items
}

func (stack *Stack) Reset() {
	stack.Empty()
	stack.Parents = make([]*[]Node, 0)
}

func (stack *Stack) Empty() {
	items := make([]Node, 0)
	stack.Items = &items
}

// --------------- Stack

// --------------- Parser

type Parser struct {
	Tokens  []Token
	Current int
	Token   Token

	// Symbols, eg var + -...
	Symbols map[string]Symbol

	Comparisions map[string]bool

	// The current stack (used by Expression)
	Stack Stack
}

func (p *Parser) Parse(tokens []Token) Block {
	p.Tokens = tokens
	p.Current = 0
	p.Symbols = make(map[string]Symbol)

	// Initialize Stack
	p.Stack.Reset()

	// var
	p.Symbol("var", func(expecting Expecting) Node {
		n := Assign{}

		name := p.Advance()

		if name.Type != "name" {
			log.Panicf("var, expected name, got %s", name.Type)
		}

		n.Name = name.Value

		eq := p.Advance()

		p.Stack.Add(&Nil{})

		// for var a in 1..2
		// for var a in ["first", "second"]
		// for var a in list
		if expecting == EXPECTING_FOR_PART && eq.Type == "keyword" && eq.Value == "in" {

			// Define an iterator object with the name that we already have
			iter := Iterate{}
			iter.Object, _ = p.Statement(EXPECTING_EXPRESSION)
			iter.Name = n.Name

			return iter
		}

		if !(eq.Type == "operator" && eq.Value == "=") {
			log.Panicf("var, expected =, got %s %s", eq.Type, eq.Value)
		}

		n.Right = p.Expressions()

		return n
	}, 0, true)

	p.SymbolCase("variable", func(expecting Expecting) Symbol {

		sym := Symbol{}
		sym.Importance = 0
		sym.IsStatement = false

		// The basic Infix function
		sym.Function = func(expecting Expecting) Node {
			return p.Expression(false)
		}

		if expecting == EXPECTING_CLASS_BODY {
			sym.Function = func(expecting Expecting) Node {
				return p.Method()
			}

			return sym
		}

		// Var as assignment
		if len(*p.Stack.Items) == 0 {
			sym.IsStatement = true
			sym.Function = func(expecting Expecting) Node {

				name := p.Token

				if name.Type != "name" {
					log.Panicf("var, expected name, got %s", name.Type)
				}

				tok := p.Advance()

				// Set
				// abc = 123
				if tok.Type == "operator" && tok.Value == "=" {
					set := Set{}
					set.Name = name.Value
					set.Right = p.Expressions()

					return set
				}

				if tok.Type == "EOL" || tok.Type == "EOF" {
					return &Nil{}
				}

				p.Reverse(2)
				return p.Expressions()
			}
		}

		return sym
	})

	// var
	p.Symbol("if", func(expecting Expecting) Node {
		i := If{}

		i.Condition = p.Expressions()
		i.True = p.Statements(EXPECTING_IF_BODY)

		p.Advance()

		if p.Token.Type == "keyword" && p.Token.Value == "else" {
			p.Advance()
			i.False = p.Statements(EXPECTING_IF_BODY)
		}

		return i
	}, 0, true)

	// Define a class
	p.Symbol("class", func(expecting Expecting) Node {

		class := DefineClass{}

		name := p.Advance()

		if name.Type != "name" {
			log.Panicf("Expected name after class, got %s (%s)", name.Type, name.Value)
		}

		p.Advance()

		p.Stack.Add(&class)

		class.Name = name.Value
		class.Body = p.Statements(EXPECTING_CLASS_BODY)

		return class

	}, 0, true)

	// Define a static method
	p.Symbol("static", func(expecting Expecting) Node {

		p.Advance()

		method := p.Method()
		method.IsStatic = true

		return method
	}, 0, true)

	// Create class instance
	p.Symbol("new", func(expecting Expecting) Node {
		inst := Instance{}

		name := p.Advance()

		if name.Type != "name" {
			log.Panicf("Expected name after new, got %s (%s)", name.Type, name.Value)
		}

		inst.Left = name.Value

		next := p.Advance()

		if next.Type != "operator" && next.Value != "(" {
			log.Panicf("Expected ( after new, got %s (%s)", name.Type, name.Value)
		}

		next = p.Advance()

		if next.Type != "operator" && next.Value != ")" {
			log.Panicf("Expected ) after new, got %s (%s)", name.Type, name.Value)
		}

		return inst
	}, 0, true)

	p.Symbol("[", func(expecting Expecting) Node {

		list := CreateList{}
		list.Items = make([]Node, 0)

		for {
			if i, ok := p.Statement(EXPECTING_NOTHING); ok {
				list.Items = append(list.Items, i)
			} else {
				break
			}
		}

		return list
	}, 0, true)

	p.Symbol("return", func(expecting Expecting) Node {

		res := Return{}

		if i, ok := p.Statement(EXPECTING_NOTHING); ok {
			res.Statement = i
		} else {
			res.Statement = Literal{Type: "null"}
		}

		return res

	}, 0, true)

	p.Symbol("for", func(expecting Expecting) Node {

		f := For{}

		// Before
		f.Before = p.Statements(EXPECTING_FOR_PART)

		// Test if we got an iterator, if that is the case we should skip to the body part directly
		if _, ok := f.Before.Body[0].(Iterate); ok {
			f.IsForIn = true
			f.Body = p.Statements(EXPECTING_NOTHING)
			return f
		}

		// Condition
		p.Advance()

		f.Condition = p.Expressions()
		p.Advance()

		// After
		f.Each = p.Statements(EXPECTING_FOR_PART)

		// For body
		f.Body = p.Statements(EXPECTING_NOTHING)

		return f

	}, 0, true)

	p.Infix("number", 0)
	p.Infix("string", 0)
	p.Infix("bool", 0)

	// Comparisions
	p.Infix("&&", 30)
	p.Infix("||", 30)
	p.Infix("==", 40)
	p.Infix("!=", 40)
	p.Infix("<", 40)
	p.Infix("<=", 40)
	p.Infix(">", 40)
	p.Infix(">=", 40)

	// Hashmap of comparisions
	p.Comparisions = make(map[string]bool)
	p.Comparisions["=="] = true
	p.Comparisions[">"] = true
	p.Comparisions[">="] = true
	p.Comparisions["<"] = true
	p.Comparisions["<="] = true
	p.Comparisions["&&"] = true
	p.Comparisions["||"] = true

	// Math
	p.Infix("+", 50)
	p.Infix("-", 50)
	p.Infix("*", 60)
	p.Infix("/", 60)

	// Builtins
	p.Infix("...", 70)
	p.Infix("..", 70)

	top := p.Statements(EXPECTING_NOTHING)

	return top
}

// Add to the symbol table
func (p *Parser) Symbol(str string, function SymbolReturn, importance int, isStatement bool) {
	p.Symbols[str] = Symbol{
		Function:    function,
		Importance:  importance,
		IsStatement: isStatement,
	}
}

func (p *Parser) SymbolCase(str string, function SymbolCaseReturn) {
	p.Symbols[str] = Symbol{
		CaseFunction: function,
	}
}

// Shortcut for adding Infix's to the symbol table
func (p *Parser) Infix(str string, importance int) {
	p.Symbol(str, func(expecting Expecting) Node {
		return p.Expression(false)
	}, importance, false)
}

func (p *Parser) Advance() Token {

	if p.Current >= len(p.Tokens) {
		p.Token = Token{
			Type: "EOF",
		}

		return p.Token
	}

	token := p.Tokens[p.Current]
	p.Token = token
	p.Current++

	return token
}

func (p *Parser) Reverse(times int) {
	p.Current -= times
}

func (p *Parser) NextToken(i int) Token {

	// End or beginning of p.Tokens (i can be negative)
	if p.Current+i >= len(p.Tokens) || p.Current+i < 0 {
		p.Token = Token{
			Type: "EOF",
		}

		return p.Token
	}

	return p.Tokens[p.Current+i]
}

func (p *Parser) Previous() Node {

	if len(*p.Stack.Items) > 0 {
		items := *p.Stack.Items
		return items[len(items)-1]
	}

	return Nil{}
}

func (p *Parser) GetOperatorImportance(str string) int {

	if _, ok := p.Symbols[str]; ok {
		return p.Symbols[str].Importance
	}

	return 0
}

func (p *Parser) Expression(advance bool) Node {

	if advance {
		p.Advance()
	}

	previous := p.Previous()
	current := p.Token

	if current.Type == "operator" && (current.Value == "}" || current.Value == "{" || current.Value == ")" || current.Value == "," || current.Value == ";") {
		return Nil{}
	}

	// Number or string
	if current.Type == "number" || current.Type == "string" || current.Type == "bool" {
		literal := Literal{
			Type:  current.Type,
			Value: current.Value,
		}

		return literal
	}

	// Variables
	if current.Type == "name" {
		variable := Variable{}
		variable.Name = current.Value

		return variable
	}

	// PushClass
	// IO.Println("123")
	//   ^
	if current.Type == "operator" && current.Value == "." {
		push := PushClass{}
		push.Left = previous

		// Convert Variable to literal
		if v, ok := push.Left.(Variable); ok {
			push.Left = Literal{
				Type:  "string",
				Value: v.Name,
			}
		}

		push.Right = p.Expressions()

		return push
	}

	// Call
	// IO.Println("123")
	//    ^^^^^^^
	if current.Type == "operator" && current.Value == "(" {

		method := Call{}
		method.Left = previous

		// Convert Variable to literal
		if v, ok := method.Left.(Variable); ok {
			method.Left = Literal{
				Type:  "string",
				Value: v.Name,
			}
		}

		method.Parameters = make([]Node, 0)

		for {
			param := p.Expressions()

			p.Advance()

			if _, ok := param.(Nil); ok {
				break
			}

			method.Parameters = append(method.Parameters, param)
		}

		return method
	}

	// We encountered an operator, check the type of the previous expression
	if current.Type == "operator" {

		math := Math{}
		math.Method = current.Value // + - * /

		// Differentiate between comparisions and arithmetic operators
		if _, ok := p.Comparisions[math.Method]; ok {
			math.IsComparision = true
		} else {
			math.IsComparision = false
		}

		prev, ok := previous.(Math)

		if ok {
			if p.GetOperatorImportance(prev.Method) < p.GetOperatorImportance(math.Method) {
				math.Left = prev.Left
				math.Method = prev.Method
				math.Right = Math{
					Method: current.Value,
					Left:   prev.Right,
					Right:  p.Expression(true),
				}
			} else {
				math.Left = previous
				math.Right = p.Expression(true)
			}
		}

		_, ok = previous.(Literal)
		if ok {
			math.Left = previous
			math.Right = p.Expression(true)
		}

		_, ok = previous.(Variable)
		if ok {
			math.Left = previous
			math.Right = p.Expression(true)
		}

		p.Stack.Add(math)

		return math
	}

	return Nil{}
}

func (p *Parser) Expressions() Node {

	p.Stack.Push()

	for {
		expression := p.Expression(true)

		if _, ok := expression.(Nil); ok {
			p.Reverse(1)
			return p.Previous()
		}

		p.Stack.Add(expression)
	}

	p.Stack.Pop()

	return Nil{}
}

func (p *Parser) Method() DefineMethod {

	method := DefineMethod{}
	method.Parameters = make([]Parameter, 0)

	if p.Token.Type != "name" {
		log.Panicf("Expecting method name, got %s (%s)", p.Token.Type, p.Token.Value)
	}

	method.Name = p.Token.Value

	// IsPublic
	if string(method.Name[0]) >= "A" && string(method.Name[0]) <= "Z" {
		method.IsPublic = true
	}

	method.Parameters = make([]Parameter, 0)

	next := p.NextToken(0)

	if next.Type == "operator" && next.Value == "(" && next.Type == "operator" && next.Value == ")" {
		method.Body = p.Statements(EXPECTING_METHOD_BODY)
	} else {
		for {

			tok := p.Advance()

			if tok.Type == "operator" && tok.Value == ")" {
				break
			}

			if tok.Type == "name" {
				param := Parameter{}
				param.Name = tok.Value
				method.Parameters = append(method.Parameters, param)
			}
		}

		method.Body = p.Statements(EXPECTING_METHOD_BODY)
	}

	return method
}

func (p *Parser) Statement(expecting Expecting) (Node, bool) {

	p.Stack.Push()

	var statement Node

	hasContent := false

	for {
		tok := p.Advance()

		if tok.Type == "EOF" || tok.Type == "EOL" {
			break
		}

		// IO.Println("first", "second")
		if tok.Type == "operator" && tok.Value == "," {
			break
		}

		if tok.Type == "operator" && tok.Value == ")" {
			break
		}

		// for var a = 0; a < 10; a++ {}
		if expecting == EXPECTING_FOR_PART && tok.Type == "operator" && tok.Value == ";" {
			hasContent = true
			break
		}

		// for var a in abc {}
		if expecting == EXPECTING_FOR_PART && tok.Type == "keyword" && tok.Value == "in" {
			hasContent = true
			break
		}

		if _, ok := p.Symbols[tok.Value]; ok {

			statement = p.Symbols[tok.Value].Function(expecting)

			p.Stack.Add(statement)

			hasContent = true

			if p.Symbols[tok.Value].IsStatement {
				break
			}

			continue
		}

		if tok.Type == "number" || tok.Type == "string" || tok.Type == "bool" {
			statement = p.Symbols[tok.Type].Function(expecting)
			p.Stack.Add(statement)
			hasContent = true
			continue
		}

		if tok.Type == "name" {
			sym := p.Symbols["variable"].CaseFunction(expecting)
			statement = sym.Function(expecting)
			hasContent = true
			break
		}

		if tok.Type == "operator" && tok.Value == "}" {
			hasContent = true
			break
		}
	}

	p.Stack.Pop()

	return statement, hasContent
}

func (p *Parser) Statements(expecting Expecting) Block {
	n := Block{}

	for {

		statement, ok := p.Statement(expecting)

		if ok && statement != nil {
			n.Body = append(n.Body, statement)
		}

		if (p.Token.Type == "operator" && p.Token.Value == "}") || p.Token.Type == "EOF" {

			// To force a new statement
			p.Token.Type = "ForceStatement"
			break
		}

		if expecting == EXPECTING_FOR_PART && (p.Token.Type == "operator" && p.Token.Value == ";" || p.Token.Type == "EOL" || p.Token.Type == "keyword" && p.Token.Value == "in") {
			p.Token.Type = "ForceStatement"
			break
		}
	}

	return n
}
