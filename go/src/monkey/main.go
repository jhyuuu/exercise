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
    "monkey/object"
)


func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hello %s! This is the Monkey programming language!\n",
        user.Username)


    // DEBUG Code
    env := object.NewEnvironment()
    input := ""
    l := lexer.New(input)
    p := parser.New(l)
    program := p.ParseProgram()
    str := program.String()
    fmt.Println(str)
    evaluator.Eval(program, env)


    // repl.StartEvaluate(os.Stdin, os.Stdout)
    repl.StartCompile(os.Stdin, os.Stdout)
}
