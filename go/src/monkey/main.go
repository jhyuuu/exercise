package main

import (
    "fmt"
    "os"
    "os/user"

    "monkey/repl"
    "monkey/lexer"
    "monkey/parser"
    // "monkey/ast"
)


func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hello %s! This is the Monkey programming language!\n",
        user.Username)


    input := "5 + 3 * 2"
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    str := program.String()
    fmt.Println(str)

    repl.Start(os.Stdin, os.Stdout)
}