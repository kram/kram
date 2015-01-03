package main

import (
    "os"
    "fmt"
    "io/ioutil"
    "strings"
    "regexp"
)

type Branch struct {
    Type string // "compare", "call", "define", "do", "word", "value"
    Value string // "123"
    Args []*Branch
}

func main() {
    files := os.Args[1:]
    fmt.Println(files)

    for _, file := range files {
        content, _ := ioutil.ReadFile(file)
        parser(string(content))
    }
}

func getNextPiece(pieces []string) (string, []string) {

    if len(pieces) == 0 {
        empty := []string{}
        return "", empty
    }

    // This is a whitespace
    if strings.TrimSpace(pieces[0]) == "" {
        return getNextPiece(pieces[1:])
    }

    return pieces[0], pieces[1:]
}

func parser(content string) Branch {

    var tree = Branch{
        Type: "do",
    }

    fmt.Println(tree)

    // Split input by whitespace
    pieces := strings.FieldsFunc(content, func (r rune) bool {
        return r == ' ' || r == '\n' || r == '\r' || r == '\t'
    })

    iter := 0

    piece := ""

    for {
        piece, pieces = getNextPiece(pieces)

        // Done!
        if piece == "" {
            return tree
        }

        // Remove whitespace from beginning
        content = strings.TrimSpace(content)

        // fmt.Println("pieces", pieces)
        fmt.Println("piece", piece)

        // var
        if piece == "var" {

            args := make([]*Branch, 2)

            var wordPiece string
            wordPiece, pieces = getNextPiece(pieces)
            word := parser(wordPiece)

            args[0] = &word


//            args[0], pieces = getNextPiece(pieces)
  //          args[1], pieces = getNextPiece(pieces)

            var branch = Branch{
                Type: "assign",
                Args: args,
            }

            tree.Args = append(tree.Args, &branch)
        }

        // Number
        matchNumber, _ := regexp.MatchString("[0-9]", piece)
        if matchNumber {
            // Numbers are returned as they can't live on their own
            return Branch{
                Type: "value",
                Value: piece,
            }
        }

        // Word
        matchWord, _ := regexp.MatchString("[a-Z]", piece)
        if matchWord {
            // Numbers are returned as they can't live on their own
            return Branch{
                Type: "word",
                Value: piece,
            }
        }


        iter++
    }

    fmt.Println(tree)

    return tree
}