package main

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

	Tokens []Token // The result goes here
}

func (l *Lexer) Init(source string) {

	l.Operators = make(map[string]bool)
	l.Operators["="] = true
	/*l.Operators["=="] = true
	l.Operators[">"] = true
	l.Operators[">="] = true
	l.Operators["<"] = true
	l.Operators["<="] = true
	l.Operators["!"] = true
	l.Operators["!!"] = true
	l.Operators["&&"] = true
	l.Operators["||"] = true*/
	l.Operators["+"] = true
	l.Operators["-"] = true
	l.Operators["*"] = true
	l.Operators["/"] = true

	l.Length = len(source)
	l.Source = source

	l.Parse()
}

func (l *Lexer) Parse() {

	for {

		// End
		if l.I >= l.Length {
			l.Push("EOL", "")
			l.Push("EOF", "")
			break
		}

		// Get current char
		l.C = string(l.Source[l.I])

		// Line endings
		if l.C == "\n" || l.C == "\r" || l.C == "" {
			l.I++
			l.Push("EOL", "")
		}

		// Ignore Whitespace
		if strings.TrimSpace(l.C) != l.C {
			l.I++
			continue
		}

		// Comments
		if l.C == "/" && l.CharAtPos(l.I+1) == "/" {
			l.I++

			// Comments contine until the end of the file or a new row
			for {
				l.C = l.CharAtPos(l.I)

				if l.C == "\n" || l.C == "\r" || l.C == "" {
					break
				}

				l.I++
			}
			continue
		}

		// Names
		// Begins with a char a-Z
		if (l.C >= "a" && l.C <= "z") || (l.C >= "A" && l.C <= "Z") {
			str := l.C
			l.I++

			for {
				l.C = l.CharAtPos(l.I)

				// After the beginning, a name can be a-Z0-9_
				if (l.C >= "a" && l.C <= "z") || (l.C >= "A" && l.C <= "Z") || (l.C >= "0" && l.C <= "9") || l.C == "_" {
					str += l.C
					l.I++
				} else {
					break
				}
			}

			l.Push("name", str)
			continue
		}

		// Numbers
		if l.C >= "0" && l.C <= "9" {
			str := l.C
			l.I++

			// Look for more digits.
			for {
				l.C = l.CharAtPos(l.I)

				if l.C < "0" || l.C > "9" {
					break
				}

				l.I++
				str += l.C
			}

			// TODO Decimal
			// TODO Verify that it ends with a space?

			l.Push("number", str)
			continue
		}

		// Strings
		if l.C == "\"" {
			str := ""
			l.I++

			// TODO escaping

			for {
				if l.CharAtPos(l.I) == "\"" {
					l.I++
					break
				}

				l.C = l.CharAtPos(l.I)
				str += l.C
				l.I++
			}

			l.Push("string", str)
			continue
		}

		// Operators
		if _, ok := l.Operators[l.C]; ok {
			l.I++
			str := l.C

			for {
				if _, ok := l.Operators[str+l.CharAtPos(l.I)]; ok {
					l.C = l.CharAtPos(l.I)
					l.I++
					str += l.C
				} else {
					break
				}
			}

			l.Push("operator", str)
			continue
		}

		l.I++

		l.Push("operator", l.C)
	}
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