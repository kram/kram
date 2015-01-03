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
    var lexer = Lexer{}
    content, _ := ioutil.ReadFile(file)
    lexer.init(string(content))
    fmt.Println(lexer.Tokens)

    var vm = VM{}
    vm.init(lexer.Tokens)
    }
}
