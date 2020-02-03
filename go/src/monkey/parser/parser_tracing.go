package parser

import (
    "fmt"
    "strings"
)

var traceLevel int = 0

const traceIndentPlaceholder string = "    "

func indentLevel() string {
    return strings.Repeat(traceIndentPlaceholder, traceLevel - 1)
}

func tracePrint(fs string) {
    fmt.Printf("%s%s\n", indentLevel(), fs)
}

func incIndent() {
    traceLevel += 1
}

func decIndent() {
    traceLevel -= 1
}

func trace(msg string) string {
    incIndent()
    tracePrint("BEGIN " + msg)
    return msg
}

func untrace(msg string) {
    tracePrint("END " + msg)
    decIndent()
}