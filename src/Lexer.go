// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package kram

import (
	"strings"
)

type Token struct {
	Type  string
	Value string
}

// The lexers responsibility is to split the input (eg. the sourcecode) and split it into different sections (aka tokens)
// This process is pretty stragit forward and there isn't any "magic" going on here.
//
// If the input is:
//		var my_var = 1 + 2
//
// The output becomes:
// 		[
// 		  { "Type": "keyword",   "Value": "var" },
// 		  { "Type": "name",      "Value": "my_var" },
// 		  { "Type": "operator",  "Value": "=" },
// 		  { "Type": "number",    "Value": "1" },
// 		  { "Type": "operator",  "Value": "+" },
// 		  { "Type": "number",    "Value": "2" },
// 		  { "Type": "EOL",       "Value": "" },
// 		  { "Type": "EOF",       "Value": "" }
// 		]
//
type Lexer struct {
	current string // The current character
	index   int    // Index of the current character
	length  int    // Length of the source
	source  []rune

	operators map[string]bool
	keywords  map[string]bool

	tokens []Token // The result goes here
}

// Initialize and run the lexer
func (lexer *Lexer) Init(source []byte) []Token {
	lexer.operators = make(map[string]bool)
	lexer.operators["+"] = true
	lexer.operators["-"] = true
	lexer.operators["*"] = true
	lexer.operators["/"] = true
	lexer.operators["%"] = true
	lexer.operators["**"] = true
	lexer.operators["="] = true
	lexer.operators["=="] = true
	lexer.operators[">"] = true
	lexer.operators[">="] = true
	lexer.operators["<"] = true
	lexer.operators["<="] = true
	lexer.operators["&&"] = true
	lexer.operators["||"] = true
	lexer.operators["..."] = true
	lexer.operators[".."] = true
	lexer.operators["."] = true
	lexer.operators["{"] = true
	lexer.operators["}"] = true
	lexer.operators[":"] = true
	lexer.operators[","] = true

	lexer.operators["++"] = true
	lexer.operators["--"] = true

	lexer.keywords = make(map[string]bool)
	lexer.keywords["if"] = true
	lexer.keywords["else"] = true
	lexer.keywords["var"] = true
	lexer.keywords["class"] = true
	lexer.keywords["static"] = true
	lexer.keywords["return"] = true
	lexer.keywords["for"] = true
	lexer.keywords["in"] = true
	lexer.keywords["fn"] = true

	lexer.source = []rune(string(source))
	lexer.length = len(lexer.source)

	lexer.parse()

	return lexer.tokens
}

// Loop over the input from the begining to the end
func (lexer *Lexer) parse() {
	last_type := ""

	for {
		t, v := lexer.parseNext()

		// Reached the end of the file, quit properly and return
		if t == "EOF" {
			lexer.push("EOL", "")
			lexer.push("EOF", "")
			return
		}

		lexer.index++

		// Nothing happened, just continue
		if t == "" && v == "" {
			continue
		}

		// Two numbers in a row (1 2 3) results in the final number 123
		// This is done by adding to the last number if we encounter two numbers in a row
		if last_type == "number" && t == "number" {
			lexer.append(v)
			continue
		}

		last_type = t
		lexer.push(t, v)
	}
}

func (lexer *Lexer) parseNext() (string, string) {
	// End of file
	if lexer.index >= lexer.length {
		return "EOF", ""
	}

	// Get current char
	lexer.current = lexer.charAtPos(lexer.index)

	// Line endings
	if lexer.current == "\n" || lexer.current == "\r" || lexer.current == "" {
		return "EOL", ""
	}

	// Ignore Whitespace
	if strings.TrimSpace(lexer.current) != lexer.current {
		return "", ""
	}

	// Comments
	if lexer.current == "/" && lexer.charAtPos(lexer.index+1) == "/" {
		return lexer.comment()
	}

	// Names
	// Begins with a char a-Z
	if (lexer.current >= "a" && lexer.current <= "z") || (lexer.current >= "A" && lexer.current <= "Z") {
		return lexer.name()
	}

	// Numbers
	if lexer.current >= "0" && lexer.current <= "9" {
		return lexer.number()
	}

	// Strings
	if lexer.current == "\"" {
		return lexer.string()
	}

	// operators
	if _, ok := lexer.operators[lexer.current]; ok {
		return lexer.operator()
	}

	return "operator", lexer.current
}

// Get the charater at a certain offset, used to look forward and backwards
func (lexer *Lexer) charAtPos(pos int) string {

	// End of file
	if pos >= lexer.length {
		return ""
	}

	return string(lexer.source[pos])
}

// Push to the lexers stack
func (lexer *Lexer) push(typ, value string) {
	lexer.tokens = append(lexer.tokens, Token{
		Type:  typ,
		Value: value,
	})
}

// Add to the previous item on the stack
func (lexer *Lexer) append(value string) {
	lexer.tokens[len(lexer.tokens)-1].Value += value
}

// Comments contine until the end of the file or a new row
func (lexer *Lexer) comment() (string, string) {
	for {
		lexer.index++

		lexer.current = lexer.charAtPos(lexer.index)

		if lexer.current == "\n" || lexer.current == "\r" || lexer.current == "" {
			return "EOL", ""
		}
	}

	return "", ""
}

// Parse names
func (lexer *Lexer) name() (string, string) {
	str := lexer.current

	for {
		c := lexer.charAtPos(lexer.index + 1)

		// After the beginning, a name can be a-Z0-9_
		if (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || (c >= "0" && c <= "9") || c == "_" {
			str += c
			lexer.index++
		} else {
			break
		}
	}

	if str == "true" || str == "false" {
		return "bool", str
	}

	if _, ok := lexer.keywords[str]; ok {
		return "keyword", str
	}

	return "name", str
}

// Parse numberss
func (lexer *Lexer) number() (string, string) {
	str := lexer.current

	// Look for more digits.
	for {
		c := lexer.charAtPos(lexer.index + 1)

		if (c < "0" || c > "9") && c != "." {
			break
		}

		// A dot needs to be followed by another digit to be valid
		if c == "." {
			cc := lexer.charAtPos(lexer.index + 2)
			if cc < "0" || cc > "9" {
				break
			}
		}

		lexer.index++
		str += c
	}

	// TODO Decimal
	// TODO Verify that it ends with a space?

	return "number", str
}

// Parse strings
func (lexer *Lexer) string() (string, string) {
	str := ""

	lexer.index++

	for {

		// End of string
		if lexer.charAtPos(lexer.index) == "\"" {
			break
		}

		// Escaping
		if lexer.charAtPos(lexer.index) == "\\" {
			lexer.index++
		}

		str += lexer.charAtPos(lexer.index)

		lexer.index++
	}

	return "string", str
}

// Parse operators
func (lexer *Lexer) operator() (string, string) {
	str := lexer.current

	for {

		next := lexer.charAtPos(lexer.index + 1)

		if next == "" {
			break
		}

		if _, ok := lexer.operators[str+next]; ok {
			lexer.index++
			str += lexer.charAtPos(lexer.index)
		} else {
			break
		}
	}

	return "operator", str
}
