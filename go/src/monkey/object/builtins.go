package object

import (
    "fmt"
)

var Builtins = []struct {
    Name    string
    Builtin *BuiltinObject
} {
    { Name :"len",   Builtin: &BuiltinObject {Fn : BuiltinFuncLen}},
    { Name :"puts",  Builtin: &BuiltinObject {Fn : BuiltinFuncPuts}},
    { Name :"first", Builtin: &BuiltinObject {Fn : BuiltinFuncFirst}},
    { Name :"last",  Builtin: &BuiltinObject {Fn : BuiltinFuncLast}},
    { Name :"rest",  Builtin: &BuiltinObject {Fn : BuiltinFuncRest}},
    { Name :"push",  Builtin: &BuiltinObject {Fn : BuiltinFuncPush}},
}

func GetBuiltinByName(name string) *BuiltinObject {
    for _, def := range Builtins {
        if def.Name == name {
            return def.Builtin
        }
    }

    return nil
}

func BuiltinFuncLen(args ...Object) Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *String:
            return &Integer{Value : int64(len(arg.Value))}
        case *ArrayObject:
            return &Integer{Value : int64(len(arg.Elements))}
        default:
            return newErrorObejct("argument to `len` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncPuts(args ...Object) Object {
    for _, arg := range args {
        fmt.Println(arg.Inspect())
    }
    return nil
}

func BuiltinFuncFirst(args ...Object) Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *ArrayObject:
            if len(arg.Elements) > 0 {
                return arg.Elements[0]
            } else {
                return nil
            }
        default:
            return newErrorObejct("argument to `first` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncLast(args ...Object) Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *ArrayObject:
            var length int = len(arg.Elements)
            if length > 0 {
                return arg.Elements[length - 1]
            } else {
                return nil
            }
        default:
            return newErrorObejct("argument to `last` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncRest(args ...Object) Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *ArrayObject:
            var length int = len(arg.Elements)
            if length > 0 {
                newElements := make([]Object, length - 1, length - 1)
                copy(newElements, arg.Elements[1:length])
                return &ArrayObject{Elements : newElements}
            } else {
                return nil
            }
        default:
            return newErrorObejct("argument to `rest` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncPush(args ...Object) Object {
    if len(args) != 2 {
        return newErrorObejct("wrong number of arguments. got=%d, want=2", len(args))
    }

    switch arg := args[0].(type) {
        case *ArrayObject:
            var length int = len(arg.Elements)
            newElements := make([]Object, length + 1, length + 1)
            copy(newElements, arg.Elements)
            newElements[length] = args[1]

            return &ArrayObject{Elements : newElements}
        default:
            return newErrorObejct("argument to `push` not supported, got %s", arg.Type())
    }
}

func newErrorObejct(format string, a ...interface{}) *ErrorObject {
    return &ErrorObject{Message : fmt.Sprintf(format, a...)}
}