package main

import (
	"log"
)

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
	Math   bool
	Method string
	Left   Node
	Right  Node
}

type If struct {
	If        bool
	Condition Condition
	True      Block
	False     Block
}

type Condition struct {
	Condition string // && || > < >= <=
	Left      Node
	Right     Node
}

type Symbol struct {
	Function     SymbolReturn
	CaseFunction SymbolCaseReturn
	Importance   int
	IsStatement  bool
}

type SymbolReturn func() Node
type SymbolCaseReturn func() Symbol

type Parser struct {
	Tokens  []Token
	Current int
	Token   Token

	// Symbols, eg var + -...
	Symbols map[string]Symbol

	// The current stack (used by Expression)
	Stack []Node

	// Current Statement()
	Stat        map[int]Node
	CurrentStat int
}

func (p *Parser) Parse(tokens []Token) Block {
	p.Tokens = tokens
	p.Current = 0
	p.Symbols = make(map[string]Symbol)
	p.Stat = make(map[int]Node)

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
		if len(p.Stack) == 0 {
			sym.IsStatement = true
			sym.Function = func() Node {
				n := Set{}

				name := p.Token

				if name.Type != "name" {
					log.Panicf("var, expected name, got %s", name.Type)
				}

				n.Name = name.Value

				eq := p.Advance()

				if eq.Type != "operator" && eq.Value == "=" {
					log.Panicf("var, expected =, got %s %s", eq.Type, eq.Value)
				}

				// Put Nil on the stack
				p.Stack = append(p.Stack, &Nil{})

				stat, ok := p.Statement()

				if ok {
					n.Right = stat
				} else {
					log.Panic("Found no statement to Assign")
				}

				return n
			}
		}

		return sym
	})

	p.Infix("number", 0)
	p.Infix("string", 0)
	p.Infix("bool", 0)

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

func (p *Parser) GetOperatorImportance(str string) int {

	if _, ok := p.Symbols[str]; ok {
		return p.Symbols[str].Importance
	}

	return 0
}

func (p *Parser) Previous() Node {
	if len(p.Stack) > 0 {
		return p.Stack[len(p.Stack)-1]
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

		p.Stack = append(p.Stack, literal)

		return literal
	}

	// Variables
	if current.Type == "name" {
		variable := Variable{}
		variable.Name = current.Value

		p.Stack = append(p.Stack, variable)

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

		p.Stack = make([]Node, 0)
		p.Stack = append(p.Stack, math)

		return math
	}

	return Nil{}
}

func (p *Parser) Statement() (Node, bool) {

	p.CurrentStat++
	current := p.CurrentStat

	hasContent := false

	for {
		tok := p.Advance()

		if tok.Type == "EOF" || tok.Type == "EOL" {
			break
		}

		if _, ok := p.Symbols[tok.Value]; ok {
			p.Stat[current] = p.Symbols[tok.Value].Function()
			hasContent = true

			if p.Symbols[tok.Value].IsStatement {
				break
			}

			continue
		}

		if tok.Type == "number" || tok.Type == "string" || tok.Type == "bool" {
			p.Stat[current] = p.Symbols[tok.Type].Function()
			hasContent = true
			continue
		}

		if tok.Type == "name" {
			sym := p.Symbols["variable"].CaseFunction()
			p.Stat[current] = sym.Function()
			hasContent = true
			continue
		}
	}

	p.CurrentStat--

	return p.Stat[current], hasContent
}

func (p *Parser) Statements() Block {
	n := Block{}

	for {
		p.Stack = make([]Node, 0)

		if (p.Token.Value == "}") || p.Token.Type == "EOF" {
			break
		}

		statement, ok := p.Statement()

		if ok {
			n.Body = append(n.Body, statement)
		}
	}

	return n
}
