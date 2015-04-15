// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package gus

import (
	"fmt"
	"log"
	"strings"

	ins "github.com/zegl/Gus/src/instructions"
)

// --------------- Symbols

type Symbol struct {
	Function   SymbolFunction
	Importance int
}

type SymbolFunction func(ON) ins.Node

// --------------- Symbols

type ON int

const (
	ON_DEFAULT     ON = 1 << iota // 1
	ON_CLASS_BODY                 // 2
	ON_PUSH_CLASS                 // 4
)

// --------------- Stack

type Stack struct {
	Items   *[]ins.Node
	Parents []*[]ins.Node
}

func (stack *Stack) Pop() {
	if len(stack.Parents) == 0 {
		items := make([]ins.Node, 0)
		stack.Items = &items
		return
	}

	stack.Items = stack.Parents[len(stack.Parents)-1]
	stack.Parents = stack.Parents[:len(stack.Parents)-1]
}

func (stack *Stack) Push() {
	stack.Parents = append(stack.Parents, stack.Items)

	items := make([]ins.Node, 0)
	stack.Items = &items
}

func (stack *Stack) Add(node ins.Node) {
	items := *stack.Items
	items = append(items, node)

	stack.Items = &items
}

func (stack *Stack) Reset() {
	stack.Empty()
	stack.Parents = make([]*[]ins.Node, 0)
}

func (stack *Stack) Empty() {
	items := make([]ins.Node, 0)
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

	Comparisions  map[string]bool
	StartOperators  map[string]bool
	LeftOnlyInfix map[string]bool
	RightOnlyInfix map[string]bool

	// The current stack (used by Expression)
	Stack Stack

	// Used for debugging
	Depth int
	Debug bool
}

func (p *Parser) Log(change int, str string, a ...interface{}) {

	if !p.Debug {
		return
	}

	if change > 0 {
		p.Depth += change
	}

	fmt.Print(strings.Repeat("--", p.Depth), str)
	fmt.Println(a)

	if change <= 0 {
		p.Depth += change
	}
}

func (p *Parser) Parse(tokens []Token) ins.Block {

	p.Log(1, "Parse()")

	p.Tokens = tokens
	p.Current = 0
	p.Symbols = make(map[string]Symbol)

	// Initialize Stack
	p.Stack.Reset()

	p.Symbol("var", p.Symbol_var, 0)
	p.Symbol("if", p.Symbol_if, 0)
	p.Symbol("class", p.Symbol_class, 0)
	//p.Symbol("static", p.Symbol_static, 0)
	p.Symbol("new", p.Symbol_new, 0)
	p.Symbol("return", p.Symbol_return, 0)
	p.Symbol("for", p.Symbol_for, 0)

	p.Symbol("[", p.Symbol_list, 5)
	p.Symbol("{", p.Symbol_map, 5)

	p.Symbol("name", p.Symbol_name, 2)

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

	// 123++
	p.LeftOnlyInfix = make(map[string]bool)
	p.LeftOnlyInfix["++"] = true
	p.LeftOnlyInfix["--"] = true

	// -123
	p.RightOnlyInfix = make(map[string]bool)
	p.RightOnlyInfix["-"] = true

	// List of all operators starting a new sub-expression
	// Starting off with a clonse of p.Comparisions
	p.StartOperators = make(map[string]bool)
	p.StartOperators["=="] = true
	p.StartOperators[">"] = true
	p.StartOperators[">="] = true
	p.StartOperators["<"] = true
	p.StartOperators["<="] = true
	p.StartOperators["&&"] = true
	p.StartOperators["||"] = true
	p.StartOperators["++"] = true
	p.StartOperators["--"] = true
	p.StartOperators["-"] = true
	p.StartOperators["+"] = true
	p.StartOperators["*"] = true
	p.StartOperators["/"] = true
	p.StartOperators["("] = true
	p.StartOperators["="] = true
	p.StartOperators[".."] = true
	p.StartOperators["..."] = true

	// Math
	p.Infix("+", 50)
	p.Infix("-", 50)
	p.Infix("*", 60)
	p.Infix("/", 60)

	// Builtins
	p.Infix("...", 70)
	p.Infix("..", 70)

	p.Infix(".", 80)
	p.Infix("(", 80)
	p.Infix("=", 80)
	p.Infix("++", 80)
	p.Infix("--", 80)

	file := p.ParseFile()

	p.Log(-1, "Parse()")

	return file
}

// Add to the symbol table
func (p *Parser) Symbol(str string, function SymbolFunction, importance int) {
	p.Symbols[str] = Symbol{
		Function:   function,
		Importance: importance,
	}
}

//
// Shortcut for adding Infix's to the symbol table
//
func (p *Parser) Infix(str string, importance int) {
	p.Symbol(str, func(on ON) ins.Node {
		return p.ParseStatementPart(on)
	}, importance)
}

//
// Get the next token in p.Tokens
//
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

func (p *Parser) AdvanceAndExpect(t, v string) Token {
	next := p.Advance()

	if next.Type != t {
		log.Panicf("Expected %s %s got %s %s", t, v, next.Type, next.Value)
	}

	if v != "" && next.Value != v {
		log.Panicf("Expected %s %s got %s %s", t, v, next.Type, next.Value)	
	}

	return next
}

//
// Reverse progress made by p.Advance()
//
func (p *Parser) Reverse(times int) {
	p.Current -= times
}

//
// Take a sneek-peak et the next token
//
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

func (p *Parser) GetOperatorImportance(str string) int {
	if _, ok := p.Symbols[str]; ok {
		return p.Symbols[str].Importance
	}

	return 0
}

func (p *Parser) ParseNext(advance bool) ins.Node {
	return p.ParseNextWithON(advance, ON_DEFAULT)
}

func (p *Parser) ParseNextWithON(advance bool, on ON) ins.Node {
	if advance {
		p.Advance()
	}

	tok := p.Token

	p.Log(1, "ParseNext() (Start) ", on, tok)

	if tok.Type == "number" || tok.Type == "string" || tok.Type == "bool" || tok.Type == "name" {
		a := p.Symbols[tok.Type].Function(on)
		p.Log(-1, "ParseNext() (End) ", tok)
		return p.LookAheadWithON(a, on)
	}

	if _, ok := p.Symbols[tok.Value]; ok {
		a := p.Symbols[tok.Value].Function(on)
		p.Log(-1, "ParseNext() (End) ", tok)
		return p.LookAheadWithON(a, on)
	}

	p.Log(-1, "ParseNext() (Nil) ", tok)

	return &ins.Nil{}
}

func (p *Parser) ReadUntil(until []Token) ins.Node {
	return p.ReadUntilWithON(until, ON_DEFAULT)
}

func (p *Parser) ReadUntilWithON(until []Token, on ON) (res ins.Node) {
	p.Log(1, "ReadUntil() (Start)", on, until)

	res = &ins.Nil{}

	p.Stack.Push()
	first := true

	for {

		// Multiple statements can end at the same EOL
		if !first {
			for _, t := range until {
				if (t.Type == "EOL" || (t.Type == "operator" && t.Value == ";")) && p.Token.Type == t.Type {
					p.Log(-1, "ReadUntil() (Premature End)", until)
					p.Stack.Pop()
					return
				}
			}
		}

		first = false
		p.Advance()

		for _, t := range until {
			if p.Token.Type == t.Type && p.Token.Value == t.Value {
				p.Log(-1, "ReadUntil() (End)", until)
				p.Stack.Pop()
				return
			}
		}

		r := p.ParseNextWithON(false, on)

		if _, ok := r.(*ins.Nil); ok {
			p.Log(0, "ReadUntil()", "Was nil, not overwriting...")
			continue
		}

		res = r
		p.Stack.Add(r)
	}

	p.Stack.Pop()
	p.Log(-1, "ReadUntil() (End)", until)

	return
}

func (p *Parser) ParseBlock() ins.Block {
	return p.ParseBlockWithON(ON_DEFAULT)
}

func (p *Parser) ParseBlockWithON(on ON) ins.Block {
	p.Log(1, "ParseBlock()")

	block := ins.Block{}

	for {
		i := p.ReadUntilWithON([]Token{Token{"EOF", ""}, Token{"EOL", ""}, Token{"operator", "}"}}, on)

		if _, ok := i.(*ins.Nil); !ok {
			block.Body = append(block.Body, i)
		}

		if p.Token.Type == "operator" && p.Token.Value == "}" {
			p.Log(-1, "ParseBlock()")
			return block
		}

		if p.Token.Type == "EOF" {
			p.Log(-1, "ParseBlock() EOF")
			return block
		}
	}

	p.Log(-1, "ParseBlock()")
	return block
}

func (p *Parser) ParseFile() ins.Block {

	p.Log(1, "ParseFile()")

	block := ins.Block{}

	for {
		if next := p.NextToken(1); next.Type == "EOF" {
			break
		}

		block.Body = append(block.Body, p.ParseBlock())
	}

	p.Log(-1, "ParseFile()")
	return block
}

func (p *Parser) TopOfStack() ins.Node {
	if len(*p.Stack.Items) > 0 {
		items := *p.Stack.Items
		return items[len(items)-1]
	}

	return ins.Nil{}
}

func (p *Parser) LookAhead(in ins.Node) ins.Node {
	return p.LookAheadWithON(in, ON_DEFAULT)
}

func (p *Parser) LookAheadWithON(in ins.Node, on ON) ins.Node {
	next := p.NextToken(0)

	// PushClass
	// IO.Println("123")
	//   ^
	if next.Type == "operator" && next.Value == "." {
		return p.SymbolPushClass(in)
	} 

	// Call
	// IO.Println("123")
	//           ^
	if next.Type == "operator" && next.Value == "(" {
		return p.SymbolCall(in, on)
	}

	// We encountered an operator, check the type of the previous expression
	if next.Type == "operator" {
		if _, ok := p.StartOperators[next.Value]; ok  {
			return p.SymbolMath(in)
		}

		return in
	}

	// Default is to do nothing
	return in
}

func (p *Parser) SymbolPushClass(in ins.Node) ins.Node {
	p.Log(1, "SymbolPushClass()")

	p.AdvanceAndExpect("operator", ".")

	push := ins.PushClass{}
	push.Left = in

	// Convert Variable to literal
	if v, ok := push.Left.(ins.Variable); ok {
		push.Left = ins.Literal{
			Type:  "string",
			Value: v.Name,
		}
	}

	push.Right = p.ParseNextWithON(true, ON_PUSH_CLASS)

	p.Log(-1, "SymbolPushClass()")

	return p.LookAhead(push)
}

func (p *Parser) SymbolCall(in ins.Node, on ON) ins.Node {
	p.Log(1, "SymbolCall()")

	p.AdvanceAndExpect("operator", "(")

	// Method definitions
	if on == ON_CLASS_BODY {
		if variable, ok := in.(ins.Variable); ok {
			return p.Symbol_MethodWithName(variable.Name)
		}

		log.Panic("Encountered unknown in ON_CLASS_BODY")
	}

	// Calling a method
	if on == ON_PUSH_CLASS {
		call := ins.Call{}
		call.Parameters = p.ParseParameters()

		// Convert Variable to literal
		if v, ok := in.(ins.Variable); ok {
			call.Left = ins.Literal{
				Type:  "string",
				Value: v.Name,
			}
		}

		return p.LookAhead(call)
	}

	log.Panic("SymbolCall() on", on)

	return in
}

func (p *Parser) SymbolMath(previous ins.Node) ins.Node {
	p.Log(1, "SymbolMath()")

	current := p.NextToken(0)

	p.Advance()

	math := ins.Math{}
	math.Method = current.Value // + - * /

	// Differentiate between comparisions and arithmetic operators
	if _, ok := p.Comparisions[math.Method]; ok {
		math.IsComparision = true
	} else {
		math.IsComparision = false
	}

	if prev, ok := previous.(ins.Math); ok {
		if p.GetOperatorImportance(prev.Method) < p.GetOperatorImportance(math.Method) {
			math.Left = prev.Left
			math.Method = prev.Method
			math.Right = ins.Math{
				Method: current.Value,
				Left:   prev.Right,
				Right:  p.ParseNext(true),
			}
		} else {
			math.Left = previous
			math.Right = p.ParseNext(true)
		}

		return p.LookAhead(math)
	}

	_, isLeftOnly := p.LeftOnlyInfix[math.Method]
	_, isRightOnly := p.RightOnlyInfix[math.Method]

	if _, ok := previous.(ins.Literal); ok {
		math.Left = previous

		if !isLeftOnly {
			math.Right = p.ParseNext(true)
		}

		return p.LookAhead(math)
	}

	if _, ok := previous.(ins.Variable); ok {
		math.Left = previous

		if !isLeftOnly {
			math.Right = p.ParseNext(true)
		}

		return p.LookAhead(math)
	}

	if isRightOnly {
		math.Left = p.ParseNext(true)

		return p.LookAhead(math)
	}

	math.Left = previous
	math.Right = p.ParseNext(true)

	p.Log(-1, "SymbolMath()")

	return p.LookAhead(math)
}

func (p *Parser) ParseStatementPart(on ON) ins.Node {

	previous := p.TopOfStack()
	current := p.Token

	p.Log(1, "ParseStatementPart()", current, previous)

	// Number or string
	if current.Type == "number" || current.Type == "string" || current.Type == "bool" {
		literal := ins.Literal{
			Type:  current.Type,
			Value: current.Value,
		}

		p.Log(-1, "ParseStatementPart()")

		return p.LookAhead(literal)
	}

	// Variables
	if current.Type == "name" {
		variable := ins.Variable{}
		variable.Name = current.Value

		p.Log(-1, "ParseStatementPart()")

		return p.LookAheadWithON(variable, on)
	}

	p.Log(-1, "ParseStatementPart()")

	return p.LookAhead(previous)

	return ins.Nil{}
}

func (p *Parser) ParseParameters() []ins.Node {

	p.Log(1, "ParseParameters()")

	params := make([]ins.Node, 0)

	for {
		next := p.NextToken(-1)

		// We're done here
		if (next.Type == "operator" && next.Value == ")") || next.Type == "EOL" || next.Type == "EOF" {
			break
		}

		param := p.ReadUntil([]Token{Token{"operator", ")"}, Token{"operator", ","}, Token{"EOF", ""}, Token{"EOL", ""}})

		if _, ok := param.(*ins.Nil); !ok {
			params = append(params, param)
		}
	}

	p.Log(-1, "ParseStatementPart()")

	return params
}

func (p *Parser) Symbol_var(on ON) ins.Node {
	n := ins.Assign{}

	name := p.Advance()

	if name.Type != "name" {
		log.Panicf("var, expected name, got %s", name.Type)
	}

	n.Name = name.Value

	next := p.NextToken(0)

	// As in:
	// for var name in
	if next.Type == "keyword" && next.Value == "in" {
		return n
	}

	if next.Type == "operator" && next.Value == "=" {
		p.Advance()
	} else {
		log.Panicf("var, expected = got %s, %s", next.Type, next.Value)
	}

	n.Right = p.ReadUntil([]Token{Token{"EOL", ""}, Token{"EOF", ""}, Token{"operator", "}"}, Token{"operator", ";"}})

	return n
}

func (p *Parser) Symbol_name(on ON) ins.Node {
	// Var as assignment
	if len(*p.Stack.Items) == 0 {
		name := p.Token

		if name.Type != "name" {
			log.Panicf("var, expected name, got %s", name.Type)
		}

		next := p.NextToken(0)

		// Set
		// abc = 123
		if next.Type == "operator" && next.Value == "=" {
			set := ins.Set{}
			set.Name = name.Value

			p.Advance()

			set.Right = p.ReadUntil([]Token{Token{"EOL", ""}})

			return set
		}

		/*if next.Type == "EOL" || next.Type == "EOF" {
			fmt.Println("Should we really end up here? 81238nadouas8u")
			return &ins.Nil{}
		}*/

		return p.ParseStatementPart(on)
	}

	return p.ParseStatementPart(on)
}

func (p *Parser) Symbol_if(on ON) ins.Node {
	i := ins.If{}

	i.Condition = p.ReadUntil([]Token{Token{"operator", "{"}})

	i.True = p.ParseBlock()
	i.True.Scope = true // Create new scope

	next := p.NextToken(0)

	if next.Type == "keyword" && next.Value == "else" {
		p.AdvanceAndExpect("keyword", "else")
		p.AdvanceAndExpect("operator", "{")
		i.False = p.ParseBlock()
		i.False.Scope = true // Create new scope
	}

	return i
}

func (p *Parser) Symbol_class(on ON) ins.Node {
	class := ins.DefineClass{}

	name := p.Advance()

	if name.Type != "name" {
		log.Panicf("Expected name after class, got %s (%s)", name.Type, name.Value)
	}

	block_start := p.Advance()

	if block_start.Type != "operator" || block_start.Value != "{" {
		log.Panicf("Expected { after class name, got %s (%s)", name.Type, name.Value)
	}

	class.Name = name.Value
	class.Body = p.ParseBlockWithON(ON_CLASS_BODY)
	class.Body.Scope = true

	return class
}

/*
func (p *Parser) Symbol_static() ins.Node {
	p.Advance()

	method := p.Symbol_method()
	method.IsStatic = true

	return method
}
*/

func (p *Parser) Symbol_new(on ON) ins.Node {
	inst := ins.Instance{}

	name := p.Advance()

	if name.Type != "name" {
		log.Panicf("Expected name after new, got %s (%s)", name.Type, name.Value)
	}

	inst.Left = name.Value

	next := p.Advance()

	if next.Type != "operator" && next.Value != "(" {
		log.Panicf("Expected ( after new, got %s (%s)", name.Type, name.Value)
	}

	inst.Parameters = p.ParseParameters()

	return inst
}

func (p *Parser) Symbol_list(on ON) ins.Node {
	if len(*p.Stack.Items) == 0 {
		return p.Symbol_ListCreate()
	}

	return p.Symbol_ListAccess()
}

func (p *Parser) Symbol_ListCreate() ins.Node {
	list := ins.ListCreate{}
	list.Items = make([]ins.Node, 0)

	for {
		next := p.NextToken(-1)

		if next.Type == "operator" && next.Value == "]" {
			break
		}

		list.Items = append(list.Items, p.ReadUntil([]Token{Token{"operator", ","}, Token{"operator", "]"}}))
	}

	return list
}

func (p *Parser) Symbol_ListAccess() ins.Node {
	access := ins.AccessChildItem{}
	access.Item = p.TopOfStack()
	access.Right = p.ReadUntil([]Token{Token{"operator", "]"}, Token{"EOF", ""}})

	return access
}

func (p *Parser) Symbol_return(on ON) ins.Node {
	res := ins.Return{}
	res.Statement = p.ReadUntil([]Token{Token{"EOL", ""}, Token{"EOF", ""}, Token{"operator", "}"}})

	return res
}

func (p *Parser) Symbol_map(on ON) ins.Node {
	m := ins.MapCreate{}
	m.Keys = make([]ins.Node, 0)
	m.Values = make([]ins.Node, 0)

	is_key := true

	for {
		next := p.NextToken(-1)

		if next.Type == "operator" && next.Value == "}" {
			break
		}

		read := p.ReadUntil([]Token{Token{"operator", ","}, Token{"operator", ":"}, Token{"operator", "}"}})

		if _, ok := read.(*ins.Nil); ok {
			return m
		}

		if is_key {
			m.Keys = append(m.Keys, read)
		} else {
			m.Values = append(m.Values, read)
		}

		is_key = !is_key
	}

	return m
}

func (p *Parser) Symbol_for(on ON) ins.Node {
	f := ins.For{}

	f.Before = p.ReadUntil([]Token{Token{"operator", ";"}, Token{"keyword", "in"}})

	next := p.NextToken(-1)

	if next.Type == "keyword" && next.Value == "in" {
		return p.Symbol_for_in(f)
	}

	return p.Symbol_for_normal(f)
}

func (p *Parser) Symbol_for_normal(f ins.For) ins.For {
	f.Condition = p.ReadUntil([]Token{Token{"operator", ";"}})
	f.Each = p.ReadUntil([]Token{Token{"operator", "{"}})

	f.Body = p.ParseBlock()
	f.Body.Scope = true

	return f
}

func (p *Parser) Symbol_for_in(f ins.For) ins.For {
	f.IsForIn = true
	f.Each = p.ReadUntil([]Token{Token{"operator", "{"}})

	f.Body = p.ParseBlock()
	f.Body.Scope = true

	return f
}

func (p *Parser) Symbol_MethodWithName(name string) ins.DefineMethod {
	method := ins.DefineMethod{}
	method.Parameters = make([]ins.Parameter, 0)

	method.Name = name

	// IsPublic
	if string(method.Name[0]) >= "A" && string(method.Name[0]) <= "Z" {
		method.IsPublic = true
	}

	for {
		next := p.NextToken(-1)

		// We're done where when the next char is a )
		if next.Type == "operator" && next.Value == ")" {
			break
		}

		param := p.ReadUntil([]Token{Token{"operator", ")"}, Token{"operator", ","}, Token{"operator", "{"}, Token{"EOF", ""}})

		// Convert Variable{} to a Parameter{}
		// They are basically the same, but not really
		if v, ok := param.(ins.Variable); ok {
			par := ins.Parameter{}
			par.Name = v.Name

			method.Parameters = append(method.Parameters, par)
		}
	}

	block_start := p.Advance()

	if block_start.Type != "operator" || block_start.Value != "{" {
		log.Panicf("Expected { after method name, got %s (%s)", block_start.Type, block_start.Value)
	}

	method.Body = p.ParseBlock()
	method.Body.Scope = true

	return method
}
