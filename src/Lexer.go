package main

import (
	"errors"
	"regexp"
	"strings"
)

type Lexer struct {
	source string
	pos    int
	length int
	Tokens []GusToken
}

func (l *Lexer) init(source string) {
	l.source = source
	l.pos = 0
	l.length = len(source)

	l.tokenize()
}

func (l *Lexer) peek() (string, error) {
	return l.peekSteps(1)
}

func (l *Lexer) peekSteps(steps int) (string, error) {

	if l.pos+steps >= l.length {
		return "", errors.New("End of file")
	}

	return string(l.source[l.pos : l.pos+steps]), nil
}

func (l *Lexer) emitToken(token Token) {
	tok := GusToken{Token: token}

	l.Tokens = append(l.Tokens, tok)
}

func (l *Lexer) emitTokenValue(token Token, value string) {
	tok := GusToken{
		Token: token,
		Value: value,
	}

	l.Tokens = append(l.Tokens, tok)
}

func (l *Lexer) next() (string, error) {

	if l.pos >= l.length {
		return "", errors.New("End of file")
	}

	// Get next char
	var char = l.source[l.pos]

	l.advance()

	return string(char), nil
}

func (l *Lexer) advance() {
	l.pos++
}

func (l *Lexer) advanceMulti(steps int) {
	l.pos += steps
}

func (l *Lexer) tokenize() {
	for {
		c, err := l.next()

		// End of file
		if err != nil {
			break
		}

		if l.isKeyword(c) {
			l.emitToken(l.readKeyword(c))
			continue
		}

		if l.isName(c) {
			name := l.readNameOrValue(c)
			l.emitTokenValue(TOKEN_NAME, name)
			continue
		}

		if l.isNumber(c) {
			value := l.readNameOrValue(c)
			l.emitTokenValue(TOKEN_NUMBER, value)
			continue
		}

		switch c {
		case "(":
			l.emitToken(TOKEN_LEFT_PAREN)
		case ")":
			l.emitToken(TOKEN_RIGHT_PAREN)

		case "=":
			peek, _ := l.peek()

			if peek == "=" {
				l.advance()
				l.emitToken(TOKEN_EQEQ)
			} else {
				l.emitToken(TOKEN_EQ)
			}

		case "+":

			peek, _ := l.peek()

			if peek == "=" {
				l.advance()
				l.emitToken(TOKEN_PLUSEQ)
			} else {
				l.emitToken(TOKEN_PLUS)
			}

		case "-":
			peek, _ := l.peek()

			if peek == "=" {
				l.advance()
				l.emitToken(TOKEN_MINUSEQ)
			} else {
				l.emitToken(TOKEN_MINUS)
			}
		}
	}
}

func (l *Lexer) isKeyword(c string) bool {

	keywords := make(map[string]Token)
	keywords["var"] = TOKEN_VAR
	keywords["if"] = TOKEN_IF
	keywords["else"] = TOKEN_ELSE

	str := l.readUntilWhitespace()

	// Append first char
	str = c + str

	// Test if str is a keyword
	if _, ok := keywords[str]; ok {
		return true
	}

	return false
}

func (l *Lexer) readKeyword(c string) Token {

	keywords := make(map[string]Token)
	keywords["var"] = TOKEN_VAR
	keywords["if"] = TOKEN_IF
	keywords["else"] = TOKEN_ELSE

	str := l.readUntilWhitespace()

	l.advanceMulti(len(str))

	// Append first char
	str = c + str

	// Test if str is a keyword
	if val, ok := keywords[str]; ok {
		return val
	}

	return TOKEN_LEXER_ERROR
}

func (l *Lexer) isName(c string) bool {
	str := l.readUntilWhitespace()

	match, _ := regexp.MatchString("^[a-zA-Z]+$", c+str)

	return match
}

func (l *Lexer) isNumber(c string) bool {
	str := l.readUntilWhitespace()

	match, _ := regexp.MatchString("^[0-9]+$", c+str)

	return match
}

func (l *Lexer) readNameOrValue(c string) string {
	str := l.readUntilWhitespace()

	l.advanceMulti(len(str))

	return c + str
}

func (l *Lexer) readUntilWhitespace() string {

	length := 0
	str := ""

	for {
		if l.pos+length > len(l.source) {
			return str
		}

		str = string(l.source[l.pos : l.pos+length])

		if strings.TrimSpace(str) != str {
			return strings.TrimSpace(str)
		}

		length++
	}

	return ""
}
