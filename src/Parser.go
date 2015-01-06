package main

import (
	"log"
	"fmt"
)

// --------------- Symbols

type Symbol struct {
	Function     SymbolReturn
	CaseFunction SymbolCaseReturn
	Importance   int
	IsStatement  bool
}

type SymbolReturn func() Node
type SymbolCaseReturn func() Symbol

// --------------- Symbols

// --------------- Stack

type Stack struct {
	Items  *[]Node
	Parents []*[]Node
}

func (stack *Stack) Pop() {
	if len(stack.Parents) == 0 {
		items := make([]Node, 0)
		stack.Items = &items
		return
	}
	
	stack.Items = stack.Parents[len(stack.Parents) - 1]
	stack.Parents = stack.Parents[:len(stack.Parents) - 1]
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
	p.Symbol("var", func() Node {
		n := Assign{}

		name := p.Advance()

		if name.Type != "name" {
			log.Panicf("var, expected name, got %s", name.Type)
		}

		n.Name = name.Value

		eq := p.Advance()

		if eq.Type != "operator" && eq.Value == "=" {
			log.Panicf("var, expected =, got %s %s", eq.Type, eq.Value)
		}

		stat, ok := p.Statement()

		if ok {
			n.Right = stat
		} else {
			log.Panicf("Found no statement to Assign to %s", n.Name)
		}

		return n
	}, 0, true)

	p.SymbolCase("variable", func() Symbol {

		sym := Symbol{}
		sym.Importance = 0
		sym.IsStatement = false

		// The basic Infix function
		sym.Function = func() Node {
			return p.Expression(false)
		}

		// Var as assignment
		if len(*p.Stack.Items) == 0 {
			sym.IsStatement = true
			sym.Function = func() Node {
				
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

					// Put Nil on the stack
					p.Stack.Add(&Nil{})

					stat, ok := p.Statement()

					if ok {
						set.Right = stat
					} else {
						log.Panic("Found no statement to Assign")
					}

					return set
				}

				// Calls
				// IO.Println("123")
				// ^^
				if tok.Type == "operator" && tok.Value == "." {
					class := CallClass{}

					class.Left = name.Value
					method, _ := p.Statement()
					class.Method = method

					return class
				}

				// Calls
				// IO.Println("123")
				//    ^^^^^^^
				if tok.Type == "operator" && tok.Value == "(" {
					method := Call{}

					method.Left = name.Value
					method.Parameters = make([]Node, 0)

					for {
						stat, _ := p.Statement()
						method.Parameters = append(method.Parameters, stat)

						if p.Token.Type == "operator" && p.Token.Value == "," {
							continue
						}

						break
					}

					return method
				}

				p.Reverse()
				return p.Expression(false)
			}
		}

		return sym
	})

	// var
	p.Symbol("if", func() Node {
		i := If{}

		// Put Nil on the stack
		p.Stack.Add(&Nil{})

		stat, ok := p.Statement()

		if ok {
			i.Condition = stat
		} else {
			log.Panic("Found no statement to If")
		}

		i.True = p.Statements()

		p.Advance()

		if p.Token.Type == "keyword" && p.Token.Value == "else" {
			p.Advance()
			i.False = p.Statements()
		}

		return i
	}, 0, true)

	// Define a class
	p.Symbol("class", func() Node {

		class := DefineClass{}

		name := p.Advance()

		if name.Type != "name" {
			log.Panicf("Expected name after class, got %s (%s)", name.Type, name.Value)
		}

		class.Name = name.Value
		class.Body = p.Statements()

		return class

	}, 0, true)

	// Define a static method
	p.Symbol("static", func() Node {

		p.Advance()

		method := p.Method()
		method.IsStatic = true

		return method
	}, 0, true)

	p.Infix("number", 0)
	p.Infix("string", 0)
	p.Infix("bool", 0)

	p.Infix("&&", 30);
	p.Infix("||", 30);

	p.Infix("==", 40);
	p.Infix("!=", 40);
	p.Infix("<", 40);
	p.Infix("<=", 40);
	p.Infix(">", 40);
	p.Infix(">=", 40);

	p.Infix("+", 50)
	p.Infix("-", 50)
	p.Infix("*", 60)
	p.Infix("/", 60)

	top := p.Statements()

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
	p.Symbol(str, func() Node {
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

func (p *Parser) Reverse() {
	p.Current--
	p.Current--
}

func (p *Parser) GetOperatorImportance(str string) int {

	if _, ok := p.Symbols[str]; ok {
		return p.Symbols[str].Importance
	}

	return 0
}

func (p *Parser) Previous() Node {

	if len(*p.Stack.Items) > 0 {
		items := *p.Stack.Items
		return items[len(items)-1]
	}

	return Nil{}
}

func (p *Parser) Expression(advance bool) Node {

	if advance {
		p.Advance()
	}

	previous := p.Previous()
	current := p.Token

	// Number or string
	if current.Type == "number" || current.Type == "string" || current.Type == "bool" {
		literal := Literal{
			Type:  current.Type,
			Value: current.Value,
		}

		p.Stack.Add(literal)

		return literal
	}

	// Variables
	if current.Type == "name" {
		variable := Variable{}
		variable.Name = current.Value

		p.Stack.Add(variable)

		return variable
	}

	// We encountered an operator, check the type of the previous expression
	if current.Type == "operator" {

		math := Math{}
		math.Method = current.Value // + - * /

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

		p.Stack.Empty()
		p.Stack.Add(math)

		return math
	}

	return Nil{}
}

func (p *Parser) Method() DefineMethod {

	method := DefineMethod{}
	method.Parameters = make([]Parameter, 0)

	if p.Token.Type != "name" {
		log.Panicf("Expecting method name, got %s (%s)", p.Token.Type, p.Token.Value)
	}

	method.Name = p.Token.Value
	method.Body = p.Statements()

	fmt.Println(method)

	return method
}

func (p *Parser) Statement() (Node, bool) {

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

		if _, ok := p.Symbols[tok.Value]; ok {
			statement = p.Symbols[tok.Value].Function()
			hasContent = true

			if p.Symbols[tok.Value].IsStatement {
				break
			}

			continue
		}

		if tok.Type == "number" || tok.Type == "string" || tok.Type == "bool" {
			statement = p.Symbols[tok.Type].Function()
			hasContent = true
			continue
		}

		if tok.Type == "name" {
			sym := p.Symbols["variable"].CaseFunction()
			statement = sym.Function()
			hasContent = true
			continue
		}

		if tok.Type == "operator" && tok.Value == "}" {
			hasContent = true
			break
		}

		// log.Panicf("How do I handle %s %s?\n", tok.Type, tok.Value)
	}

	p.Stack.Pop()

	return statement, hasContent
}

func (p *Parser) Statements() Block {
	n := Block{}

	for {

		p.Stack.Push()

		statement, ok := p.Statement()

		if ok && statement != nil {
			n.Body = append(n.Body, statement)
		}

		p.Stack.Pop()

		if (p.Token.Type == "operator" && p.Token.Value == "}") || p.Token.Type == "EOF" {

			// To force a new statement
			p.Token.Type = "ForceStatement"
			break
		}
	}

	return n
}
