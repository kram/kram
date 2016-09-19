// Copyright (c) 2015 The kram Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution.
// This file may not be copied, modified, or distributed except according to those terms.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kram/kram/src"
	"github.com/kram/kram/src/compiler"
)

// The main function is very simple.
//
// 1. Read the input file
// 2. Run the lexer
// 3. Run the parser on the result from the lexer and create an AST
// 4. Run the VM with the AST as instructions
func main() {

	debug := flag.Bool("debug", false, "Debuggning output")
	flag.Parse()

	files := os.Args[1:]

	for _, file := range files {
		content, err := ioutil.ReadFile(file)

		if err != nil {
			continue
		}

		var lexer = kram.Lexer{}
		tokens := lexer.Init(content)

		if *debug {
			fmt.Println("-------------------")
			fmt.Println("-       LEXER     -")
			fmt.Println("-------------------")

			b, _ := json.MarshalIndent(tokens, "", "  ")
			fmt.Println(string(b))

			fmt.Println("-------------------")
			fmt.Println("-   INSTRUCTIONS  -")
			fmt.Println("-------------------")
		}

		var parse = kram.Parser{}

		tree := parse.Parse(tokens)

		if *debug {
			b, _ := json.MarshalIndent(tree, "", "  ")
			fmt.Println(string(b))

			fmt.Println("-------------------")
			fmt.Println("-        VM       -")
			fmt.Println("-------------------")
		}

		/*var vm = kram.VM{}
		vm.Debug = *debug
		vm.Run(tree)*/

		compiler.Run(tree)
	}
}
