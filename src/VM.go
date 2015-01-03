package main

import (
    "fmt"
    "log"
)

type Name struct {
    Name  string
    Value Type
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
        log.Panicf("Expected TOKEN_NAME, got %s", name)
    }

    v.advance()

    // Equal sign
    eq := v.peek()

    if eq.Token != TOKEN_EQ {
        log.Panicf("Expected TOKEN_EQ, got %s", name)
    }

    v.advance()

    // Variable value
    value := v.peek()

    if value.Token == TOKEN_NUMBER {

        number := new(Number)
        number.Init(value.Value)

        v.names[name.Value] = Name{
            Name:  name.Value,
            Value: number,
            Type:  "Number",
        }

        v.advance()

        return
    }

    if value.Token == TOKEN_STRING {

        str := new(String)
        str.Init(value.Value)

        v.names[name.Value] = Name{
            Name:  name.Value,
            Value: str,
            Type:  "String",
        }

        v.advance()

        return
    }

    log.Panicf("Expected TOKEN_NUMBER or TOKEN_STRING, got %f", name)
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
    /*variable := v.names[name.Name]

    name_value, _ := strconv.Atoi(name.Value)
    argument_value, _ := strconv.Atoi(argument.Value)

    if method == "+" {

        if variable.Type != "Number" {
            log.Fatalf("You can only add to Numbers, %s given", variable.Type)
            return
        }

        result := name_value + argument_value
        variable.Value = strconv.Itoa(result)
        v.names[name.Name] = variable
    }

    if method == "-" {

        if variable.Type != "Number" {
            log.Fatalf("You can only substract from Numbers, %s given", variable.Type)
            return
        }

        result := name_value - argument_value
        variable.Value = strconv.Itoa(result)
        v.names[name.Name] = variable
    }*/
}

func (v *VM) set(name Name, value GusToken) {

    // TOKEN_NUMBER
    // TOKEN_NAME
    // TOKEN_STRING

    // a = 123
    if value.Token == TOKEN_NUMBER {

        if name.Type != "Number" {
            log.Fatalf("Can only set variables of type Number to a number")
        }

        number := new(Number)
        number.Init(value.Value)

        v.names[name.Name] = Name{
            Name:  name.Name,
            Value: number,
            Type:  "Number",
        }

        v.advance()
        return
    }

    // a = "Hello World!"
    if value.Token == TOKEN_STRING {

        if name.Type != "String" {
            log.Fatalf("Can only set variables of type String to a string")
        }

        str := new(String)
        str.Init(value.Value)

        v.names[name.Name] = Name{
            Name:  name.Name,
            Value: str,
            Type:  "String",
        }

        v.advance()
        return
    }


    // a = b
    if value.Token == TOKEN_NAME {

        // Get reference
        ref := v.getName(value.Value)

        if ref.Type == "Number" {

            number := new(Number)
            number.Init(ref.Value.toString())

            v.names[name.Name] = Name{
                Name:  name.Name,
                Value: number,
                Type:  ref.Type,
            }

            v.advance()
            return
        }

        if ref.Type == "String" {

            str := new(String)
            str.Init(ref.Value.toString())

            v.names[name.Name] = Name{
                Name:  name.Name,
                Value: str,
                Type:  ref.Type,
            }

            v.advance()
            return
        }
    }

    log.Fatalf("Was unable to set %s to %s", name.Type, value.Token)
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

        case TOKEN_EQ:
            name := v.getName(v.previous_name)
            v.set(name, v.peek())

        // Increase variable
        case TOKEN_PLUSEQ:
            name := v.getName(v.previous_name)
            v.call(name, "+", v.peek())
            v.advance()

        // Decrease variable
        case TOKEN_MINUSEQ:
            name := v.getName(v.previous_name)
            v.call(name, "+", v.peek())
            v.advance()
        }
    }
}
