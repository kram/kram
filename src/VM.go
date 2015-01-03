package main

import (
    "fmt"
    "log"
)

type Name struct {
    Name  string
    Value string
    Type  string
}

type VM struct {
    names         map[string]Name
    previous_name string

    tokens  []GusToken
    current int
}

func (v *VM) init(tokens []GusToken) {
    v.tokens = tokens
    v.current = 0
    v.names = make(map[string]Name)

    v.run()
}

func (v *VM) peek() GusToken {

    if v.current >= len(v.tokens) {
        return GusToken{Token: TOKEN_END_OF_PROGRAM}
    }

    return v.tokens[v.current]
}

func (v *VM) nextToken() GusToken {
    if v.current >= len(v.tokens) {
        return GusToken{Token: TOKEN_END_OF_PROGRAM}
    }

    token := v.tokens[v.current]

    v.advance()

    return token
}

func (v *VM) advance() {
    v.current++
}

func (v *VM) allocateVar() {

    // Variable name
    name := v.peek()

    if name.Token != TOKEN_NAME {
        log.Panic("Expected TOKEN_NAME, got ", name)
    }

    v.advance()

    // Equal sign
    eq := v.peek()

    if eq.Token != TOKEN_EQ {
        log.Panic("Expected TOKEN_EQ, got ", name)
    }

    v.advance()

    // Variable value
    value := v.peek()

    if value.Token != TOKEN_VALUE {
        log.Panic("Expected TOKEN_VALUE, got ", name)
    }

    v.names[name.Value] = Name{
        Name:  name.Value,
        Value: value.Value,
        Type:  "number",
    }

    v.advance()
}

func (v *VM) verifyName(str string) bool {
    if _, ok := v.names[str]; ok {
        v.previous_name = str
        return true
    }

    return false
}

func (v *VM) getName(str string) Name {
    return v.names[str]
}

func (v *VM) call(name Name, method string, argument GusToken) {
    if method == "+" {

        variable := v.names[name.Name]
        variable.Value = name.Value + argument.Value
        v.names[name.Name] = variable
    }
}

func (v *VM) run() {
    for {
        token := v.nextToken()

        fmt.Println(token)

        if token.Token == TOKEN_END_OF_PROGRAM {
            fmt.Println("END OF PROGRAM")
            break
        }

        switch token.Token {
        case TOKEN_VAR:
            v.allocateVar()

        case TOKEN_NAME:
            if !v.verifyName(token.Value) {
                log.Panicf("Unexpected TOKEN_NAME %s", token.Value)
            }

        case TOKEN_PLUS:
            name := v.getName(v.previous_name)
            v.call(name, "+", v.peek())
            v.advance()
        }
    }

    fmt.Println(v.names)
}
