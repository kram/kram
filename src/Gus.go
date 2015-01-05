package main

import (
    "fmt"
    "io/ioutil"
    "os"
)

func main() {
    files := os.Args[1:]
    fmt.Println(files)

    for _, file := range files {
        content, _ := ioutil.ReadFile(file)

        var lexer = Lexer{}
        lexer.Init(string(content))

        var parse = Parser{}
        parse.Parse(lexer.Tokens)
    }
}