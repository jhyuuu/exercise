package evaluator

import (
    "monkey/object"
)

func BuiltinFuncLen(args ...object.Object) object.Object {
    if len(args) != 1 {
        return newErrorObejct("wrong number of arguments. got=%d, want=1", len(args))
    }

    switch arg := args[0].(type) {
        case *object.String:
            return &object.Integer{Value : int64(len(arg.Value))}
        default:
            return newErrorObejct("argument to `len` not supported, got %s", arg.Type())
    }
}