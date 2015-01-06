package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	filepath.Walk("test/tests", Test)
}

func Test(path string, file os.FileInfo, err error) error {

	// No not test the dir :)
	if file.IsDir() {
		return nil
	}

	content, _ := ioutil.ReadFile(path)

	expect := ""
	
	re, _ := regexp.Compile(`(?m)^// (.*?)$`)
	for _, str := range re.FindAllString(string(content), -1) {
		expect += strings.Replace(str, "// ", "", -1) + "\n"
	}

	// Normalize newlines
	expect = strings.Replace(expect, "\r\n", "\n", -1)

	cmd := exec.Command("Gus.exe", path)
	stdout, err := cmd.Output()

	if err != nil {
		println(err.Error())
		return nil
	}

	if expect == string(stdout) {
		fmt.Printf("1: %s\n", path)
	} else {
		fmt.Printf("0: %s\n", path)
		fmt.Printf("Expected\n---\n'%s'---\ngot\n---\n'%s'---\n", expect, string(stdout))
	}

	return nil
}
