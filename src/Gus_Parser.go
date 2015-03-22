package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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

	Depth int
}

func (p *Parser) Log(change int, str string, a ...interface{}) {

	if change > 0 {
		p.Depth += change
	}

	fmt.Print(strings.Repeat("--", p.Depth), str)
	fmt.Println(a)

	if change <= 0 {
		p.Depth += change
	}
}

func (p *Parser) Parse(tokens []Token) Block {

	p.Log(1, "Parse()")

	p.Tokens = tokens
	p.Current = 0
	p.Symbols = make(map[string]Symbol)

	// Initialize Stack
	p.Stack.Reset()

	p.Symbol("var", p.Symbol_var, 0, true)
	p.Symbol("if", p.Symbol_if, 0, true)
	p.Symbol("class", p.Symbol_class, 0, true)
	//p.Symbol("static", p.Symbol_static, 0, true)
	p.Symbol("new", p.Symbol_new, 0, true)
	//p.Symbol("[", p.Symbol_list, 0, true)
	//p.Symbol("return", p.Symbol_return, 0, true)
	p.Symbol("for", p.Symbol_for, 0, true)

	p.SymbolCase("variable", p.Symbol_variable)

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

//
// Shortcut for adding Infix's to the symbol table
//
func (p *Parser) Infix(str string, importance int) {
	p.Symbol(str, func(expecting Expecting) Node {
		return p.ParseStatementPart()
	}, importance, false)
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

func (p *Parser) ParseNext(advance bool) Node {

	if advance {
		p.Advance()
	}

	tok := p.Token

	p.Log(1, "ParseNext() (Start) ", tok)

	expecting := EXPECTING_NOTHING

	if _, ok := p.Symbols[tok.Value]; ok {
		a := p.Symbols[tok.Value].Function(expecting)
		p.Log(-1, "ParseNext() (End) ", tok)
		return a
	}

	if tok.Type == "number" || tok.Type == "string" || tok.Type == "bool" {
		a := p.Symbols[tok.Type].Function(expecting)
		p.Log(-1, "ParseNext() (End) ", tok)
		return a
	}

	if tok.Type == "name" {
		sym := p.Symbols["variable"].CaseFunction(expecting)
		a := sym.Function(expecting)
		p.Log(-1, "ParseNext() (End) ", tok)
		return a
	}

	p.Log(-1, "ParseNext() (Nil) ", tok)

	return &Nil{}
}

func (p *Parser) ReadUntil(until []Token) (res Node) {
	p.Log(1, "ReadUntil() (Start)", until)

	res = &Nil{}

	p.Stack.Push()

	first := true

	for {

		if !first {
			for _, t := range until {
				if p.Token.Type == t.Type && p.Token.Value == t.Value {
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

		fmt.Println()
		r := p.ParseNext(false)

		if _, ok := r.(Nil); ok {
			fmt.Println("Was nil, not overwriting...")
			p.Log(0, "ReadUntil()", "Was nil, not overwriting...")
			continue
		}

		if _, ok := r.(*Nil); ok {
			fmt.Println("Was nil, not overwriting...")
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

func (p *Parser) ParseBlock() Block {

	p.Log(1, "ParseBlock()")

	block := Block{}

	for {
		i := p.ReadUntil([]Token{Token{"EOF", ""}, Token{"EOL", ""}, Token{"operator", "}"}})

		b, _ := json.MarshalIndent(i, "", "  ")
		fmt.Println(string(b))

		if _, ok := i.(Nil); ok {
			p.Log(-1, "ParseBlock() Was nil")
			return block
		}

		block.Body = append(block.Body, i)

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

func (p *Parser) ParseFile() Block {

	p.Log(1, "ParseFile()")

	block := Block{}

	for {
		if next := p.NextToken(1); next.Type == "EOF" {
			break
		}

		block.Body = append(block.Body, p.ParseBlock())
	}

	p.Log(-1, "ParseFile()")
	return block
}

func (p *Parser) TopOfStack() Node {
	if len(*p.Stack.Items) > 0 {
		items := *p.Stack.Items
		return items[len(items)-1]
	}

	return Nil{}
}

func (p *Parser) ParseStatementPart() Node {

	previous := p.TopOfStack()
	current := p.Token

	p.Log(1, "ParseStatementPart()", current, previous)

	// Number or string
	if current.Type == "number" || current.Type == "string" || current.Type == "bool" {
		literal := Literal{
			Type:  current.Type,
			Value: current.Value,
		}

		p.Log(-1, "ParseStatementPart()")

		return literal
	}

	// Variables
	if current.Type == "name" {
		variable := Variable{}
		variable.Name = current.Value

		p.Log(-1, "ParseStatementPart()")

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

		push.Right = p.ParseNext(true)

		p.Log(-1, "ParseStatementPart()")

		return push
	}

	// Assignment
	if current.Type == "operator" && current.Value == "=" {
		if assignment, ok := previous.(Assign); ok {
			assignment.Right = p.ParseNext(true)
			return assignment
		}

		log.Panicf("Expected previous to be an assignment, it wasn't")
	}

	// Call
	// IO.Println("123")
	//           ^
	if current.Type == "operator" && current.Value == "(" {

		call := Call{}
		call.Parameters = make([]Node, 0)

		// Get parameters
		for {
			next := p.NextToken(0)

			if next.Type == "operator" && next.Value == ")" {
				break
			}

			call.Parameters = append(call.Parameters, p.ParseNext(true))
		}

		// Put Call{} into the a previous PushClass if neeccesary
		if push, ok := previous.(PushClass); ok {
			call.Left = push.Right

			// Convert Variable to literal
			if v, ok := call.Left.(Variable); ok {
				call.Left = Literal{
					Type:  "string",
					Value: v.Name,
				}
			}

			push.Right = call

			p.Log(-1, "ParseStatementPart()")

			return push
		}

		// When the previous was a name
		// This is now a method definition
		if variable, ok := previous.(Variable); ok {
			return p.Symbol_MethodWithName(variable.Name)
		}


		// Leave this to see if it actually can happen
		call.Left = previous
		fmt.Println("This happened, 918238yyhaUSHDHASD")
		os.Exit(1)

		p.Log(-1, "ParseStatementPart()")

		return call
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
					Right:  p.ParseNext(true),
				}
			} else {
				math.Left = previous
				math.Right = p.ParseNext(true)
			}
		}

		_, ok = previous.(Literal)
		if ok {
			math.Left = previous
			math.Right = p.ParseNext(true)
		}

		_, ok = previous.(Variable)
		if ok {
			math.Left = previous
			math.Right = p.ParseNext(true)
		}

		p.Log(-1, "ParseStatementPart()")

		return math
	}

	p.Log(-1, "ParseStatementPart()")

	return Nil{}
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

func (p *Parser) Symbol_var(expecting Expecting) Node {
	n := Assign{}

	name := p.Advance()

	if name.Type != "name" {
		log.Panicf("var, expected name, got %s", name.Type)
	}

	n.Name = name.Value

	next := p.NextToken(0)

	// eq := p.Advance()

	// p.Stack.Add(&Nil{})

	// for var a in 1..2
	// for var a in ["first", "second"]
	// for var a in list
	if expecting == EXPECTING_FOR_PART && next.Type == "keyword" && next.Value == "in" {

		fmt.Println("Got keyword in, 1290802180280")

		// Define an iterator object with the name that we already have
		iter := Iterate{}
		//iter.Object, _ = p.Statement(EXPECTING_EXPRESSION)
		//iter.Name = n.Name

		return iter
	}

	b, _ := json.MarshalIndent(n, "", "  ")
	fmt.Println(string(b))

	return n

	//if !(eq.Type == "operator" && eq.Value == "=") {
	//	log.Panicf("var, expected =, got %s %s", eq.Type, eq.Value)
	//}

	// todo
	// n.Right = p.Expressions()

	return n
}

func (p *Parser) Symbol_variable(expecting Expecting) Symbol {
	sym := Symbol{}
	sym.Importance = 0
	sym.IsStatement = false

	// The basic Infix function
	sym.Function = func(expecting Expecting) Node {
		return p.ParseStatementPart()
	}

	if expecting == EXPECTING_CLASS_BODY {
		sym.Function = func(expecting Expecting) Node {
			return Nil{}
			// todo
			//return p.Symbol_method()
		}

		return sym
	}

	// Var as assignment
	if len(*p.Stack.Items) == 0 {

		fmt.Println("VAR AS ASSIGNMENT")

		sym.IsStatement = true
		sym.Function = func(expecting Expecting) Node {

			name := p.Token

			if name.Type != "name" {
				log.Panicf("var, expected name, got %s", name.Type)
			}

			next := p.NextToken(0)

			// Set
			// abc = 123
			if next.Type == "operator" && next.Value == "=" {
				set := Set{}
				set.Name = name.Value

				p.Advance()

				set.Right = p.ReadUntil([]Token{Token{"EOL", ""}})

				b, _ := json.MarshalIndent(set, "", "  ")
				fmt.Println(string(b))

				return set
			}

			if next.Type == "EOL" || next.Type == "EOF" {
				fmt.Println("Should we really end up here? 81238nadouas8u")
				return &Nil{}
			}

			return p.ParseStatementPart()
		}
	}

	return sym
}

func (p *Parser) Symbol_if(expecting Expecting) Node {

	fmt.Println("Symbol_if()")

	i := If{}

	i.Condition = p.ReadUntil([]Token{Token{"operator", "{"}})

	i.True = p.ParseBlock()

	next := p.NextToken(0)

	if next.Type == "keyword" && next.Value == "else" {
		i.False = p.ParseBlock()
	}

	return i
}

func (p *Parser) Symbol_class(expecting Expecting) Node {
	class := DefineClass{}

	name := p.Advance()

	if name.Type != "name" {
		log.Panicf("Expected name after class, got %s (%s)", name.Type, name.Value)
	}

	class.Name = name.Value
	class.Body = p.ParseBlock()

	return class
}

/*
func (p *Parser) Symbol_static(expecting Expecting) Node {
	p.Advance()

	method := p.Symbol_method()
	method.IsStatic = true

	return method
}
*/

func (p *Parser) Symbol_new(expecting Expecting) Node {
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
}

/*
func (p *Parser) Symbol_list(expecting Expecting) Node {
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
}
*/

/*
func (p *Parser) Symbol_return(expecting Expecting) Node {
	res := Return{}

	if i, ok := p.Statement(EXPECTING_NOTHING); ok {
		res.Statement = i
	} else {
		res.Statement = Literal{Type: "null"}
	}

	return res
}
*/

func (p *Parser) Symbol_for(expecting Expecting) Node {
	f := For{}

	f.Before = p.ReadUntil([]Token{Token{"operator", ";"}})
	f.Condition = p.ReadUntil([]Token{Token{"operator", ";"}})
	f.Each = p.ReadUntil([]Token{Token{"operator", "{"}})

	f.Body = p.ParseBlock()

	return f

	// Test if we got an iterator, if that is the case we should skip to the body part directly
	//if _, ok := f.Before.Body[0].(Iterate); ok {
	f.IsForIn = true
	//f.Body = p.Statements(EXPECTING_NOTHING)
	return f
	//}

	return f
}

func (p *Parser) Symbol_MethodWithName(name string) DefineMethod {
	// Initialize
	method := DefineMethod{}
	method.Parameters = make([]Parameter, 0)

	method.Name = name

	// IsPublic
	if string(method.Name[0]) >= "A" && string(method.Name[0]) <= "Z" {
		method.IsPublic = true
	}

	for {
		next := p.NextToken(0)

		// We're done where when the next char is a )
		if next.Type == "operator" && next.Value == ")" {
			break
		}

		param := p.ReadUntil([]Token{Token{"operator", ")"}, Token{"operator", ","}})

		fmt.Println("TODO: Convert parameters correctly (123sjsjsjass)")
		fmt.Println(param)
		os.Exit(1)

		//method.Parameters = append(method.Parameters, param)
	}

	method.Body = p.ParseBlock()

	/*next := p.NextToken(0)

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
	}*/

	return method
}

/*
func (p *Parser) Symbol_method() DefineMethod {
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
*/
