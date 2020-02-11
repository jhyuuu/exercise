package object

import (
    "bytes"
    "fmt"
    "strings"
    "hash/fnv"

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

    ARRAY_OBJ = "ARRAY"
    HASH_OBJ  = "HASH"
)

type HashKey struct {
    Type  ObjectType
    Value uint64
}

type Hashable interface {
    GenHashKey() HashKey
}

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

type ArrayObject struct {
    Elements []Object
}

func (ao *ArrayObject) Type() ObjectType { return ARRAY_OBJ }
func (ao *ArrayObject) Inspect() string {
    var out bytes.Buffer

    elements := []string{}
    for _, e := range ao.Elements {
        elements = append(elements, e.Inspect())
    }

    out.WriteString("[")
    out.WriteString(strings.Join(elements, ", "))
    out.WriteString("]")

    return out.String()
}

func (b *Boolean) GenHashKey() HashKey {
    var value uint64

    if b.Value {
        value = 1
    } else {
        value = 0
    }

    return HashKey{Type : b.Type(), Value : value}
}

func (i *Integer) GenHashKey() HashKey {
    return HashKey{Type : i.Type(), Value : uint64(i.Value)}
}

func (s *String) GenHashKey() HashKey {
    h := fnv.New64a()
    h.Write([]byte(s.Value))

    return HashKey{Type : s.Type(), Value : h.Sum64()}
}

type HashPair struct {
    Key   Object
    Value Object
}

type HashObject struct {
    Pairs map[HashKey]HashPair
}

func (ho *HashObject) Type() ObjectType { return HASH_OBJ }
func (ho *HashObject) Inspect() string {
    var out bytes.Buffer

    pairs := []string{}
    for _, pair := range ho.Pairs {
        pairs = append(pairs, pair.Key.Inspect() + ": " + pair.Value.Inspect())
    }

    out.WriteString("{")
    out.WriteString(strings.Join(pairs, ", "))
    out.WriteString("}")

    return out.String()
}