package main

import (
    "fmt"
    "os"
    "os/user"

    "monkey/repl"
    "monkey/lexer"
    "monkey/parser"
    // "monkey/ast"
    "monkey/evaluator"
)


func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hello %s! This is the Monkey programming language!\n",
        user.Username)


    // DEBUG Code
    input := ""
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    str := program.String()
    fmt.Println(str)
    evaluator.Eval(program)

    repl.Start(os.Stdin, os.Stdout)
}