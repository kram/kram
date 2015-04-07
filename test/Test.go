// Copyright (c) 2015 The Gus Project Developers. All rights reserved.
// See the LICENSE file at the top-level directory of this distribution. 
// This file may not be copied, modified, or distributed except according to those terms.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

func main() {
	all := true

	dirs := os.Args[1:]

	for _, dir := range dirs {
		files, _ := ioutil.ReadDir(dir)

		for _, file := range files {
			if !Test(dir + "/" + file.Name()) {
				all = false
			}
		}
	}

	if all {
		os.Exit(0)
	}

	os.Exit(1)
}

func Test(path string) bool {

	content, _ := ioutil.ReadFile(path)

	expect := ""

	re, _ := regexp.Compile(`(?m)^// (.*?)$`)
	for _, str := range re.FindAllString(string(content), -1) {
		expect += strings.Replace(str, "// ", "", -1) + "\n"
	}

	// Normalize newlines
	expect = strings.Replace(expect, "\r\n", "\n", -1)

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("Gus.exe", path)
	} else {
		cmd = exec.Command("./bin/Gus", path)
	}

	stdout, err := cmd.Output()

	if err != nil {
		println(path, err.Error())
		return false
	}

	if expect == string(stdout) {
		fmt.Printf("1: %s\n", path)
		return true
	}

	fmt.Printf("0: %s\n", path)
	fmt.Printf("Expected\n---\n'%s'---\ngot\n---\n'%s'---\n", expect, string(stdout))

	return false
}
