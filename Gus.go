//
// Copyright 2015 Gustav Westling <hej@gustav.tv. All rights reserved.
// Use of this source code is governed by a BSD-style license that
// can be found in the LICENSE file.
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/zegl/Gus/src"
)

func main() {

	debug := flag.Bool("debug", false, "Debuggning output")
	flag.Parse()

	files := os.Args[1:]

	for _, file := range files {
		content, err := ioutil.ReadFile(file)

		if err != nil {
			continue
		}

		var lexer = gus.Lexer{}
		lexer.Init(string(content))

		if *debug {
			fmt.Println("-------------------")
			fmt.Println("-       LEXER     -")
			fmt.Println("-------------------")

			b, _ := json.MarshalIndent(lexer.Tokens, "", "  ")
			fmt.Println(string(b))

			fmt.Println("-------------------")
			fmt.Println("-   INSTRUCTIONS  -")
			fmt.Println("-------------------")
		}

		var parse = gus.Parser{}
		parse.Debug = *debug

		tree := parse.Parse(lexer.Tokens)

		if *debug {
			b, _ := json.MarshalIndent(tree, "", "  ")
			fmt.Println(string(b))

			fmt.Println("-------------------")
			fmt.Println("-        VM       -")
			fmt.Println("-------------------")
		}

		var vm = gus.VM{}
		vm.Debug = *debug

		vm.Run(tree)

		/*if *debug {
			b, _ := json.MarshalIndent(vm.Environment().Env, "", "  ")
			fmt.Println("\n-------------------")
			fmt.Println("-   ENVIRONMENT   -")
			fmt.Println("-------------------")
			fmt.Println(string(b))
		}*/
	}
}
