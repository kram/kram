// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package gus

import (
	"log"

	ins "github.com/zegl/Gus/src/instructions"
)

type Parser struct {
	tokens  []Token
	current int
	token   Token

	// Symbols, eg var + -...
	symbols map[string]Symbol

	comparisions   map[string]bool
	startOperators map[string]bool
	leftOnlyInfix  map[string]bool
	rightOnlyInfix map[string]bool

	// The current stack (used by Expression)
	stack Stack
}

func (parser *Parser) Parse(tokens []Token) ins.Block {

	parser.tokens = tokens
	parser.current = 0
	parser.symbols = make(map[string]Symbol)

	// Initialize Stack
	parser.stack.Reset()

	parser.symbol("var", parser.symbol_Assign, 0)
	parser.symbol("if", parser.symbol_If, 0)
	parser.symbol("class", parser.symbol_DefineClass, 0)
	parser.symbol("static", parser.symbol_DefineClassStatic, 0)
	parser.symbol("new", parser.symbol_New, 0)
	parser.symbol("return", parser.symbol_Return, 0)
	parser.symbol("for", parser.symbol_For, 0)

	parser.symbol("[", parser.symbol_List, 5)
	parser.symbol("{", parser.symbol_Map, 5)

	parser.symbol("name", parser.symbol_Name, 2)

	parser.infix("number", 0)
	parser.infix("string", 0)
	parser.infix("bool", 0)

	// Comparisions
	parser.infix("&&", 30)
	parser.infix("||", 30)
	parser.infix("==", 40)
	parser.infix("!=", 40)
	parser.infix("<", 40)
	parser.infix("<=", 40)
	parser.infix(">", 40)
	parser.infix(">=", 40)

	// Hashmap of comparisions
	parser.comparisions = make(map[string]bool)
	parser.comparisions["=="] = true
	parser.comparisions[">"] = true
	parser.comparisions[">="] = true
	parser.comparisions["<"] = true
	parser.comparisions["<="] = true
	parser.comparisions["&&"] = true
	parser.comparisions["||"] = true

	// 123++
	parser.leftOnlyInfix = make(map[string]bool)
	parser.leftOnlyInfix["++"] = true
	parser.leftOnlyInfix["--"] = true

	// -123
	parser.rightOnlyInfix = make(map[string]bool)
	parser.rightOnlyInfix["-"] = true

	// List of all operators starting a new sub-expression
	// Starting off with a clonse of parser.comparisions
	parser.startOperators = make(map[string]bool)
	parser.startOperators["=="] = true
	parser.startOperators[">"] = true
	parser.startOperators[">="] = true
	parser.startOperators["<"] = true
	parser.startOperators["<="] = true
	parser.startOperators["&&"] = true
	parser.startOperators["||"] = true
	parser.startOperators["++"] = true
	parser.startOperators["--"] = true
	parser.startOperators["-"] = true
	parser.startOperators["+"] = true
	parser.startOperators["*"] = true
	parser.startOperators["/"] = true
	parser.startOperators["("] = true
	parser.startOperators["="] = true
	parser.startOperators[".."] = true
	parser.startOperators["..."] = true

	// Math
	parser.infix("+", 50)
	parser.infix("-", 50)
	parser.infix("*", 60)
	parser.infix("/", 60)

	// Builtins
	parser.infix("...", 70)
	parser.infix("..", 70)

	parser.infix(".", 80)
	parser.infix("(", 80)
	parser.infix("=", 80)
	parser.infix("++", 80)
	parser.infix("--", 80)

	file := parser.parseFile()

	return file
}

// Add to the symbol table
func (parser *Parser) symbol(str string, function SymbolFunction, importance int) {
	parser.symbols[str] = Symbol{
		Function:   function,
		Importance: importance,
	}
}

// Shortcut for adding Infix's to the symbol table
func (parser *Parser) infix(str string, importance int) {
	parser.symbol(str, func(on ON) ins.Node {
		return parser.parseStatementPart(on)
	}, importance)
}

// Get the next token in parser.tokens
func (parser *Parser) advance() Token {
	if parser.current >= len(parser.tokens) {
		parser.token = Token{
			Type: "EOF",
		}

		return parser.token
	}

	token := parser.tokens[parser.current]
	parser.token = token
	parser.current++

	return token
}

func (parser *Parser) advanceAndExpect(t, v string) Token {
	next := parser.advance()

	if next.Type != t {
		log.Panicf("Expected %s %s got %s %s", t, v, next.Type, next.Value)
	}

	if v != "" && next.Value != v {
		log.Panicf("Expected %s %s got %s %s", t, v, next.Type, next.Value)
	}

	return next
}

// Reverse progress made by parser.advance()
func (parser *Parser) reverse(times int) {
	parser.current -= times
}

// Take a sneek-peak et the next token
func (parser *Parser) nextToken(i int) Token {

	// End or beginning of parser.tokens (i can be negative)
	if parser.current+i >= len(parser.tokens) || parser.current+i < 0 {
		parser.token = Token{
			Type: "EOF",
		}

		return parser.token
	}

	return parser.tokens[parser.current+i]
}

func (parser *Parser) getOperatorImportance(str string) int {
	if _, ok := parser.symbols[str]; ok {
		return parser.symbols[str].Importance
	}

	return 0
}

func (parser *Parser) readUntil(until []Token) ins.Node {
	return parser.readUntilWithON(until, ON_DEFAULT)
}

func (parser *Parser) readUntilWithON(until []Token, on ON) (res ins.Node) {
	res = &ins.Nil{}

	parser.stack.Push()
	first := true

	for {
		// Multiple statements can end at the same EOL
		if !first {
			for _, t := range until {
				if (t.Type == "EOL" || (t.Type == "operator" && t.Value == ";")) && parser.token.Type == t.Type {

					parser.stack.Pop()
					return
				}
			}
		}

		first = false
		parser.advance()

		for _, t := range until {
			if parser.token.Type == t.Type && parser.token.Value == t.Value {

				parser.stack.Pop()
				return
			}
		}

		r := parser.parseNextWithON(false, on)

		if _, ok := r.(*ins.Nil); ok {

			continue
		}

		res = r
		parser.stack.Add(r)
	}

	parser.stack.Pop()

	return
}

func (parser *Parser) parseNext(advance bool) ins.Node {
	return parser.parseNextWithON(advance, ON_DEFAULT)
}

func (parser *Parser) parseNextWithON(advance bool, on ON) ins.Node {
	if advance {
		parser.advance()
	}

	tok := parser.token

	if _, ok := parser.symbols[tok.Value]; ok {
		a := parser.symbols[tok.Value].Function(on)

		return parser.lookAheadWithON(a, on)
	}

	if tok.Type == "number" || tok.Type == "string" || tok.Type == "bool" || tok.Type == "name" {
		a := parser.symbols[tok.Type].Function(on)

		return parser.lookAheadWithON(a, on)
	}

	return &ins.Nil{}
}

func (parser *Parser) parseBlock() ins.Block {
	return parser.parseBlockWithON(ON_DEFAULT)
}

func (parser *Parser) parseBlockWithON(on ON) ins.Block {
	block := ins.Block{}

	for {
		i := parser.readUntilWithON([]Token{Token{"EOF", ""}, Token{"EOL", ""}, Token{"operator", "}"}}, on)

		if _, ok := i.(*ins.Nil); !ok {
			block.Body = append(block.Body, i)
		}

		if parser.token.Type == "operator" && parser.token.Value == "}" {

			return block
		}

		if parser.token.Type == "EOF" {

			return block
		}
	}

	return block
}

func (parser *Parser) parseFile() ins.Block {
	block := ins.Block{}

	for {
		if next := parser.nextToken(1); next.Type == "EOF" {
			break
		}

		block.Body = append(block.Body, parser.parseBlock())
	}

	return block
}

func (parser *Parser) parseStatementPart(on ON) ins.Node {

	previous := parser.topOfStack()
	current := parser.token

	// Number or string
	if current.Type == "number" || current.Type == "string" || current.Type == "bool" {
		literal := ins.Literal{
			Type:  current.Type,
			Value: current.Value,
		}

		return parser.lookAhead(literal)
	}

	// Variables
	if current.Type == "name" {
		variable := ins.Variable{}
		variable.Name = current.Value

		return parser.lookAheadWithON(variable, on)
	}

	// Operator overloading
	if current.Type == "operator" && on == ON_CLASS_BODY {
		return parser.symbol_MethodWithName(current.Value)
	}

	// Math exceptions
	if current.Type == "operator" && current.Value == "-" {
		if _, ok := parser.rightOnlyInfix[current.Value]; ok {
			parser.reverse(1)
			return parser.symbol_Math(ins.Nil{})
		}
	}

	return parser.lookAhead(previous)
}

func (parser *Parser) parseArguments() []ins.Argument {
	params := make([]ins.Argument, 0)

	for {
		next := parser.nextToken(-1)

		// We're done here
		if (next.Type == "operator" && next.Value == ")") || next.Type == "EOL" || next.Type == "EOF" {
			break
		}

		param := parser.readUntilWithON([]Token{Token{"operator", ")"}, Token{"operator", ","}, Token{"EOF", ""}, Token{"EOL", ""}}, ON_ARGUMENTS)

		if math, ok := param.(ins.Math); ok {
			if math.Method == "=" {

				name, found_name := math.Left.(ins.Variable)

				if !found_name {
					log.Panic("Named argument, could not find valud name")
				}

				params = append(params, ins.Argument{
					Name:    name.Name,
					IsNamed: true,
					Value:   math.Right,
				})

				continue
			}
		}

		if _, ok := param.(*ins.Nil); !ok {
			params = append(params, ins.Argument{
				Value: param,
			})

			continue
		}
	}

	return params
}

func (parser *Parser) topOfStack() ins.Node {
	if len(*parser.stack.Items) > 0 {
		items := *parser.stack.Items
		return items[len(items)-1]
	}

	return ins.Nil{}
}

func (parser *Parser) lookAhead(in ins.Node) ins.Node {
	return parser.lookAheadWithON(in, ON_DEFAULT)
}

func (parser *Parser) lookAheadWithON(in ins.Node, on ON) ins.Node {
	next := parser.nextToken(0)

	// PushClass
	// IO.Println("123")
	//   ^
	if next.Type == "operator" && next.Value == "." {
		return parser.symbol_PushClass(in)
	}

	// Call
	// IO.Println("123")
	//           ^
	if next.Type == "operator" && next.Value == "(" {
		return parser.symbol_Call(in, on)
	}

	// We encountered an operator, check the type of the previous expression
	if next.Type == "operator" {
		if _, ok := parser.startOperators[next.Value]; ok {
			return parser.symbol_Math(in)
		}

		return in
	}

	// Default is to do nothing
	return in
}

func (parser *Parser) symbol_PushClass(in ins.Node) ins.Node {
	parser.advanceAndExpect("operator", ".")

	push := ins.PushClass{}
	push.Left = in

	// Convert Variable to literal
	if v, ok := push.Left.(ins.Variable); ok {
		push.Left = ins.Literal{
			Type:  "string",
			Value: v.Name,
		}
	}

	push.Right = parser.parseNextWithON(true, ON_PUSH_CLASS)

	return parser.lookAhead(push)
}

func (parser *Parser) symbol_Call(in ins.Node, on ON) ins.Node {
	parser.advanceAndExpect("operator", "(")

	// Method definitions
	if on == ON_CLASS_BODY {
		if variable, ok := in.(ins.Variable); ok {
			return parser.symbol_MethodWithName(variable.Name)
		}

		log.Panic("Encountered unknown in ON_CLASS_BODY")
	}

	// Calling a method
	if on == ON_PUSH_CLASS {
		call := ins.Call{}
		call.Arguments = parser.parseArguments()

		// Convert Variable to literal
		if v, ok := in.(ins.Variable); ok {
			call.Left = ins.Literal{
				Type:  "string",
				Value: v.Name,
			}
		}

		return parser.lookAhead(call)
	}

	call := ins.Call{}
	call.Arguments = parser.parseArguments()
	call.Left = in

	return parser.lookAhead(call)
}

func (parser *Parser) symbol_Math(previous ins.Node) ins.Node {
	current := parser.nextToken(0)

	parser.advance()

	math := ins.Math{}
	math.Method = current.Value // + - * /

	// Differentiate between comparisions and arithmetic operators
	if _, ok := parser.comparisions[math.Method]; ok {
		math.IsComparision = true
	} else {
		math.IsComparision = false
	}

	if prev, ok := previous.(ins.Math); ok {
		if parser.getOperatorImportance(prev.Method) < parser.getOperatorImportance(math.Method) {
			math.Left = prev.Left
			math.Method = prev.Method
			math.Right = ins.Math{
				Method: current.Value,
				Left:   prev.Right,
				Right:  parser.parseNext(true),
			}
		} else {
			math.Left = previous
			math.Right = parser.parseNext(true)
		}

		return parser.lookAhead(math)
	}

	_, isLeftOnly := parser.leftOnlyInfix[math.Method]
	_, isRightOnly := parser.rightOnlyInfix[math.Method]

	if _, ok := previous.(ins.Literal); ok {
		math.Left = previous

		if !isLeftOnly {
			math.Right = parser.parseNext(true)
		}

		return parser.lookAhead(math)
	}

	if _, ok := previous.(ins.Variable); ok {
		math.Left = previous

		if !isLeftOnly {
			math.Right = parser.parseNext(true)
		}

		return parser.lookAhead(math)
	}

	if isRightOnly {
		math.Left = parser.parseNext(true)

		return parser.lookAhead(math)
	}

	math.Left = previous
	math.Right = parser.parseNext(true)

	return parser.lookAhead(math)
}

func (parser *Parser) symbol_Assign(on ON) ins.Node {
	n := ins.Assign{}

	name := parser.advance()

	if name.Type != "name" {
		log.Panicf("var, expected name, got %s", name.Type)
	}

	n.Name = name.Value

	next := parser.nextToken(0)

	// As in:
	// for var name in
	if next.Type == "keyword" && next.Value == "in" {
		return n
	}

	if next.Type == "operator" && next.Value == "=" {
		parser.advance()
	} else {
		log.Panicf("var, expected = got %s, %s", next.Type, next.Value)
	}

	n.Right = parser.readUntil([]Token{Token{"EOL", ""}, Token{"EOF", ""}, Token{"operator", "}"}, Token{"operator", ";"}})

	return n
}

func (parser *Parser) symbol_Name(on ON) ins.Node {
	// Var as assignment
	if len(*parser.stack.Items) == 0 {
		name := parser.token

		if name.Type != "name" {
			log.Panicf("var, expected name, got %s", name.Type)
		}

		next := parser.nextToken(0)

		// Set
		// abc = 123
		if next.Type == "operator" && next.Value == "=" {

			// Hijack and return early when dealing with method parameters
			if on == ON_METHOD_PARAMETERS || on == ON_ARGUMENTS {
				return parser.parseStatementPart(on)
			}

			set := ins.Set{}
			set.Name = name.Value

			parser.advance()

			set.Right = parser.readUntil([]Token{Token{"EOL", ""}})

			return set
		}

		return parser.parseStatementPart(on)
	}

	return parser.parseStatementPart(on)
}

func (parser *Parser) symbol_If(on ON) ins.Node {
	i := ins.If{}

	i.Condition = parser.readUntil([]Token{Token{"operator", "{"}})

	i.True = parser.parseBlock()
	i.True.Scope = true // Create new scope

	next := parser.nextToken(0)

	if next.Type == "keyword" && next.Value == "else" {
		parser.advanceAndExpect("keyword", "else")
		parser.advanceAndExpect("operator", "{")
		i.False = parser.parseBlock()
		i.False.Scope = true // Create new scope
	}

	return i
}

func (parser *Parser) symbol_DefineClass(on ON) ins.Node {
	class := ins.DefineClass{}

	name := parser.advance()

	if name.Type != "name" {
		log.Panicf("Expected name after class, got %s (%s)", name.Type, name.Value)
	}

	block_start := parser.advance()

	if block_start.Type != "operator" || block_start.Value != "{" {
		log.Panicf("Expected { after class name, got %s (%s)", name.Type, name.Value)
	}

	class.Name = name.Value
	class.Body = parser.parseBlockWithON(ON_CLASS_BODY)
	class.Body.Scope = true

	return class
}

func (parser *Parser) symbol_DefineClassStatic(on ON) ins.Node {
	parser.advance()

	method := parser.symbol_Method()
	method.IsStatic = true

	return method
}

func (parser *Parser) symbol_New(on ON) ins.Node {

	// "new" is also the name of constructors
	// This is added so that the lowercase version of the method name also works just fine
	if on == ON_CLASS_BODY {
		return parser.symbol_MethodWithName("New")
	}

	inst := ins.Instance{}

	name := parser.advance()

	if name.Type != "name" {
		log.Panicf("Expected name after new, got %s (%s)", name.Type, name.Value)
	}

	inst.Left = name.Value

	next := parser.advance()

	if next.Type != "operator" && next.Value != "(" {
		log.Panicf("Expected ( after new, got %s (%s)", name.Type, name.Value)
	}

	inst.Arguments = parser.parseArguments()

	return inst
}

func (parser *Parser) symbol_List(on ON) ins.Node {
	if len(*parser.stack.Items) == 0 {
		return parser.symbol_ListCreate()
	}

	return parser.symbol_ListAccess()
}

func (parser *Parser) symbol_ListCreate() ins.Node {
	list := ins.ListCreate{}
	list.Items = make([]ins.Node, 0)

	for {
		next := parser.nextToken(-1)

		if next.Type == "operator" && next.Value == "]" {
			break
		}

		list.Items = append(list.Items, parser.readUntil([]Token{Token{"operator", ","}, Token{"operator", "]"}}))
	}

	return list
}

func (parser *Parser) symbol_ListAccess() ins.Node {
	access := ins.AccessChildItem{}
	access.Item = parser.topOfStack()
	access.Right = parser.readUntil([]Token{Token{"operator", "]"}, Token{"EOF", ""}})

	return access
}

func (parser *Parser) symbol_Return(on ON) ins.Node {
	res := ins.Return{}
	res.Statement = parser.readUntil([]Token{Token{"EOL", ""}, Token{"EOF", ""}, Token{"operator", "}"}})

	return res
}

func (parser *Parser) symbol_Map(on ON) ins.Node {
	m := ins.MapCreate{}
	m.Keys = make([]ins.Node, 0)
	m.Values = make([]ins.Node, 0)

	is_key := true

	for {
		next := parser.nextToken(-1)

		if next.Type == "operator" && next.Value == "}" {
			break
		}

		read := parser.readUntil([]Token{Token{"operator", ","}, Token{"operator", ":"}, Token{"operator", "}"}})

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

func (parser *Parser) symbol_For(on ON) ins.Node {
	f := ins.For{}

	f.Before = parser.readUntil([]Token{Token{"operator", ";"}, Token{"keyword", "in"}})

	next := parser.nextToken(-1)

	if next.Type == "keyword" && next.Value == "in" {
		return parser.symbol_For_in(f)
	}

	return parser.symbol_For_normal(f)
}

func (parser *Parser) symbol_For_normal(f ins.For) ins.For {
	f.Condition = parser.readUntil([]Token{Token{"operator", ";"}})
	f.Each = parser.readUntil([]Token{Token{"operator", "{"}})

	f.Body = parser.parseBlock()
	f.Body.Scope = true

	return f
}

func (parser *Parser) symbol_For_in(f ins.For) ins.For {
	f.IsForIn = true
	f.Each = parser.readUntil([]Token{Token{"operator", "{"}})

	f.Body = parser.parseBlock()
	f.Body.Scope = true

	return f
}

func (parser *Parser) symbol_Method() ins.DefineMethod {
	name := parser.nextToken(-1)

	if name.Type != "name" {
		log.Panicf("Expeced name after method, got %s", name.Type)
	}

	return parser.symbol_MethodWithName(name.Value)
}

func (parser *Parser) symbol_MethodWithName(name string) ins.DefineMethod {
	method := ins.DefineMethod{}
	method.Parameters = make([]ins.Parameter, 0)

	method.Name = name

	for {
		next := parser.nextToken(-1)

		// We're done where when the next char is a )
		if next.Type == "operator" && next.Value == ")" {
			break
		}

		param := parser.readUntilWithON([]Token{Token{"operator", "="}, Token{"operator", ")"}, Token{"operator", ","}, Token{"operator", "{"}, Token{"EOF", ""}}, ON_METHOD_PARAMETERS)

		// Test if has default value
		// Will be returned as a Math{a = b}
		if math, ok := param.(ins.Math); ok {
			par := ins.Parameter{}

			if l, ok := math.Left.(ins.Variable); ok {
				par.Name = l.Name
			} else {
				log.Panic("Parameters with default value, could not find valid name")
			}

			par.Default = math.Right
			par.HasDefault = true
			method.Parameters = append(method.Parameters, par)
			continue
		}

		// Convert Variable{} to a Parameter{}
		// They are basically the same, but not really
		if v, ok := param.(ins.Variable); ok {
			par := ins.Parameter{}
			par.Name = v.Name
			method.Parameters = append(method.Parameters, par)
			continue
		}
	}

	block_start := parser.advance()

	if block_start.Type != "operator" || block_start.Value != "{" {
		log.Panicf("Expected { after method name, got %s (%s)", block_start.Type, block_start.Value)
	}

	method.Body = parser.parseBlock()
	method.Body.Scope = true

	return method
}
