package main

import (
	"encoding/json"
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
		tree := parse.Parse(lexer.Tokens)

		var vm = VM{}
		vm.Run(tree)

		b, _ := json.MarshalIndent(vm.Environment, "", "    ")
		fmt.Println(string(b))
	}
}
