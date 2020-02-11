package evaluator

import (
    "fmt"

    "monkey/object"
)

func BuiltinFuncLen(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *object.String:
            return &object.Integer{Value : int64(len(arg.Value))}
        case *object.ArrayObject:
            return &object.Integer{Value : int64(len(arg.Elements))}
        default:
            return newErrorObejct("argument to `len` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncPuts(args ...object.Object) object.Object {
    for _, arg := range args {
        fmt.Println(arg.Inspect())
    }
    return NULL
}

func BuiltinFuncFirst(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *object.ArrayObject:
            if len(arg.Elements) > 0 {
                return arg.Elements[0]
            } else {
                return NULL
            }
        default:
            return newErrorObejct("argument to `first` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncLast(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *object.ArrayObject:
            var length int = len(arg.Elements)
            if length > 0 {
                return arg.Elements[length - 1]
            } else {
                return NULL
            }
        default:
            return newErrorObejct("argument to `last` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncRest(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *object.ArrayObject:
            var length int = len(arg.Elements)
            if length > 0 {
                newElements := make([]object.Object, length - 1, length - 1)
                copy(newElements, arg.Elements[1:length])
                return &object.ArrayObject{Elements : newElements}
            } else {
                return NULL
            }
        default:
            return newErrorObejct("argument to `rest` not supported, got %s", arg.Type())
    }
}

func BuiltinFuncPush(args ...object.Object) object.Object {
    if len(args) != 2 {
        return newErrorObejct("wrong number of arguments. got=%d, want=2", len(args))
    }

    switch arg := args[0].(type) {
        case *object.ArrayObject:
            var length int = len(arg.Elements)
            newElements := make([]object.Object, length + 1, length + 1)
            copy(newElements, arg.Elements)
            newElements[length] = args[1]
            // if length > 0 {
            //     newElements := make([]object.Object, length - 1, length - 1)
            //     copy(newElements, arg.Elements[1:length])
            //     return &object.ArrayObject{Elements : newElements}
            // } else {
            //     return NULL
            // }
            return &object.ArrayObject{Elements : newElements}
        default:
            return newErrorObejct("argument to `push` not supported, got %s", arg.Type())
    }
}