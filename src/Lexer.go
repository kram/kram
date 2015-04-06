package gus

import (
	"strings"
)

type Token struct {
	Type  string
	Value string
}

type Lexer struct {
	C      string // The current character
	I      int    // Index of the current character
	Length int    // Length of the source
	Source string

	Operators map[string]bool
	Keywords  map[string]bool

	Tokens []Token // The result goes here
}

func (l *Lexer) Init(source string) {

	l.Operators = make(map[string]bool)
	l.Operators["+"] = true
	l.Operators["-"] = true
	l.Operators["*"] = true
	l.Operators["/"] = true
	l.Operators["%"] = true
	l.Operators["**"] = true
	l.Operators["="] = true
	l.Operators["=="] = true
	l.Operators[">"] = true
	l.Operators[">="] = true
	l.Operators["<"] = true
	l.Operators["<="] = true
	l.Operators["&&"] = true
	l.Operators["||"] = true
	l.Operators["..."] = true
	l.Operators[".."] = true
	l.Operators["."] = true
	l.Operators["{"] = true
	l.Operators["}"] = true
	l.Operators[":"] = true
	l.Operators[","] = true

	l.Operators["++"] = true
	l.Operators["--"] = true

	l.Keywords = make(map[string]bool)
	l.Keywords["if"] = true
	l.Keywords["else"] = true
	l.Keywords["var"] = true
	l.Keywords["class"] = true
	l.Keywords["static"] = true
	l.Keywords["return"] = true
	l.Keywords["for"] = true
	l.Keywords["in"] = true

	l.Length = len(source)
	l.Source = source

	l.Parse()
}

func (l *Lexer) Parse() {

	last_type := ""

	for {
		t, v := l.ParseNext()

		if t == "EOF" {

			// Push both EOL and EOF before quitting
			l.Push("EOL", "")
			l.Push("EOF", "")

			return
		}

		l.I++

		if t == "" && v == "" {
			continue
		}

		if last_type == "number" && t == "number" {
			l.Append(v)
			continue
		}

		last_type = t

		l.Push(t, v)
	}
}

func (l *Lexer) ParseNext() (string, string) {
	// End
	if l.I >= l.Length {
		return "EOF", ""
	}

	// Get current char
	l.C = l.CharAtPos(l.I)

	// Line endings
	if l.C == "\n" || l.C == "\r" || l.C == "" {
		return "EOL", ""
	}

	// Ignore Whitespace
	if strings.TrimSpace(l.C) != l.C {
		return "", ""
	}

	// Comments
	if l.C == "/" && l.CharAtPos(l.I+1) == "/" {

		// Comments contine until the end of the file or a new row
		for {
			l.I++

			l.C = l.CharAtPos(l.I)

			if l.C == "\n" || l.C == "\r" || l.C == "" {
				return "", ""
			}
		}

		return "", ""
	}

	// Names
	// Begins with a char a-Z
	if (l.C >= "a" && l.C <= "z") || (l.C >= "A" && l.C <= "Z") {
		str := l.C

		for {
			c := l.CharAtPos(l.I + 1)

			// After the beginning, a name can be a-Z0-9_
			if (c >= "a" && c <= "z") || (c >= "A" && c <= "Z") || (c >= "0" && c <= "9") || c == "_" {
				str += c
				l.I++
			} else {
				break
			}
		}

		if str == "true" || str == "false" {
			return "bool", str
		}

		if _, ok := l.Keywords[str]; ok {
			return "keyword", str
		}

		return "name", str
	}

	// Numbers
	if l.C >= "0" && l.C <= "9" {
		str := l.C

		// Look for more digits.
		for {
			c := l.CharAtPos(l.I + 1)

			if c < "0" || c > "9" {
				break
			}

			l.I++
			str += c
		}

		// TODO Decimal
		// TODO Verify that it ends with a space?

		return "number", str
	}

	// Strings
	if l.C == "\"" {
		str := ""

		// TODO escaping

		for {
			if l.CharAtPos(l.I+1) == "\"" {
				l.I++
				break
			}

			l.I++

			str += l.CharAtPos(l.I)
		}

		return "string", str
	}

	// Operators
	if _, ok := l.Operators[l.C]; ok {
		str := l.C

		for {

			next := l.CharAtPos(l.I+1)

			if next == "" {
				break
			}

			if _, ok := l.Operators[str+next]; ok {
				l.I++
				str += l.CharAtPos(l.I)
			} else {
				break
			}
		}

		return "operator", str
	}

	return "operator", l.C
}

func (l *Lexer) CharAtPos(pos int) string {
	if pos >= l.Length {
		return ""
	}

	return string(l.Source[pos])
}

func (l *Lexer) Push(typ, value string) {
	l.Tokens = append(l.Tokens, Token{
		Type:  typ,
		Value: value,
	})
}

func (l *Lexer) Append(value string) {
	l.Tokens[len(l.Tokens)-1].Value += value
}
