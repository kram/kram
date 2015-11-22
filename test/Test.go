// Copyright (c) 2015 The kram Project Developers. All rights reserved.
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

	// ./Test /path/to/bin /path/to/tests
	bindir := os.Args[1]
	testsdir := os.Args[2]

	files, _ := ioutil.ReadDir(testsdir)

	for _, file := range files {
		if !Test(bindir, testsdir + "/" + file.Name()) {
			all = false
		}
	}

	if all {
		os.Exit(0)
	}

	os.Exit(1)
}

func Test(bindir, path string) bool {

	content, _ := ioutil.ReadFile(path)

	expect := ""

	re, _ := regexp.Compile(`(?m)// (.*?)$`)
	for _, str := range re.FindAllString(string(content), -1) {
		expect += strings.Replace(str, "// ", "", -1) + "\n"
	}

	// Normalize newlines
	expect = strings.Replace(expect, "\r\n", "\n", -1)

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command(bindir + "/kram.exe", path)
	} else {
		cmd = exec.Command(bindir + "/kram", path)
	}

	stdout, err := cmd.Output()

	if err != nil {
		if err.Error() != "exit status 1" {
			println(path, err.Error())
			return false
		}
	}

	if expect == string(stdout) {
		fmt.Printf("1: %s\n", path)
		return true
	}

	fmt.Printf("0: %s\n", path)
	fmt.Printf("Expected\n---\n'%s'---\ngot\n---\n'%s'---\n", expect, string(stdout))

	return false
}
