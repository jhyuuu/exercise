package object

import (
    "bytes"
    "fmt"
    "strings"

    "monkey/ast"
)

type ObjectType string

const (
    NULL_OBJ  = "NULL"
    ERROR_OBJ = "ERROR"
    
    INTEGER_OBJ = "INTEGER"
    STRING_OBJ  = "STRING"
    BOOLEAN_OBJ = "BOOLEAN"
    
    RETURN_VALUE_OBJ = "RETURN_VALUE"
    
    FUNCTION_OBJ = "FUNCTION"
    BUILDTIN_OBJ = "BUILTIN"
)

type Object interface {
    Type() ObjectType
    Inspect() string
}

type Integer struct {
    Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

type String struct {
    Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string { return s.Value }

type Boolean struct {
    Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

type Null struct {}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string { return "null" }

type ReturnObject struct {
    Value Object
}

func (ro *ReturnObject) Type() ObjectType { return RETURN_VALUE_OBJ }
func (ro *ReturnObject) Inspect() string { return ro.Value.Inspect() }

type ErrorObject struct {
    Message string
}

func (eo *ErrorObject) Type() ObjectType { return ERROR_OBJ }
func (eo *ErrorObject) Inspect() string { return "ERROR: " + eo.Message }

type FunctionObject struct {
    Parameters []*ast.Identifier
    Body       *ast.BlockStatement
    Env        *Environment
}

func (fo *FunctionObject) Type() ObjectType { return FUNCTION_OBJ }
func (fo *FunctionObject) Inspect() string {
    var out bytes.Buffer

    params := []string{}
    for _, p := range fo.Parameters {
        params = append(params, p.String())
    }

    out.WriteString("fn")
    out.WriteString("(")
    out.WriteString(strings.Join(params, ", "))
    out.WriteString(") {\n")
    out.WriteString(fo.Body.String())
    out.WriteString("\n")

    return out.String()
}

type BuiltinFunction func(args ...Object) Object

type BuiltinObject struct {
    Fn BuiltinFunction
}

func (bo *BuiltinObject) Type() ObjectType { return BUILDTIN_OBJ }
func (bo *BuiltinObject) Inspect() string { return "builtin function" }