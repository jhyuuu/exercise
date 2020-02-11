package evaluator

import (
    "fmt"

    "monkey/ast"
    "monkey/object"
)

var (
    NULL  = &object.Null{}
    TRUE  = &object.Boolean{Value : true}
    FALSE = &object.Boolean{Value : false}
)

var builtins = map[string]*object.BuiltinObject{
    "len"    : &object.BuiltinObject {Fn : BuiltinFuncLen},
    "puts"   : &object.BuiltinObject {Fn : BuiltinFuncPuts},
    "first"  : &object.BuiltinObject {Fn : BuiltinFuncFirst},
    "last"   : &object.BuiltinObject {Fn : BuiltinFuncLast},
    "rest"   : &object.BuiltinObject {Fn : BuiltinFuncRest},
    "push"   : &object.BuiltinObject {Fn : BuiltinFuncPush},
}

func Eval(node ast.Node, env *object.Environment) object.Object {
    switch node := node.(type) {
        case *ast.Program:
            return evalProgram(node.Statements, env)
        case *ast.BlockStatement:
            return evalBlockStatement(node.Statements, env)
        case *ast.LetStatement:
            val := Eval(node.Value, env)
            if isErrorObejct(val) {
                return val
            }

            env.Set(node.Name.Value, val)

            // return &object.ReturnObject{Value : val}
        case *ast.ReturnStatement:
            val := Eval(node.ReturnValue, env)
            if isErrorObejct(val) {
                return val
            }

            return &object.ReturnObject{Value : val}
        case *ast.ExpressionStatement:
            return Eval(node.Expression, env)
        case *ast.Identifier:
            return evalIdentifier(node, env)
        case *ast.IntegerLiteral:
            return &object.Integer{Value : node.Value}
        case *ast.StringLiteral:
            return &object.String{Value : node.Value}
        case *ast.Boolean:
            return nativaBool2BooleanObject(node.Value)
        case *ast.PrefixExpression:
            right := Eval(node.Right, env)
            if isErrorObejct(right) {
                return right
            }

            return evalPrefixExpression(node.Operator, right)
        case *ast.InfixExpression:
            left  := Eval(node.Left, env)
            if isErrorObejct(left) {
                return left
            }

            right := Eval(node.Right, env)
            if isErrorObejct(right) {
                return right
            }

            return evalInfixExpression(node.Operator, left, right)
        case *ast.IfExpression:
            return evalIfExpression(node, env)
        case *ast.FunctionLiteral:
            params := node.Parameters
            body := node.Body
            return &object.FunctionObject{Parameters : params, Body : body, Env : env}
        case *ast.CallExpression:
            function := Eval(node.Function, env)
            if isErrorObejct(function) {
                return function
            }

            args := evalExpressions(node.Arguments, env)
            if len(args) == 1 && isErrorObejct(args[0]) {
                return args[0]
            }

            return applyFunction(function, args)
        case *ast.ArrayLiteral:
            elements := evalExpressions(node.Elements, env)
            if len(elements) == 1 && isErrorObejct(elements[0]) {
                return elements[0]
            }
            return &object.ArrayObject{Elements : elements}
        case *ast.IndexExpression:
            left  := Eval(node.Left, env)
            if isErrorObejct(left) {
                return left
            }

            right := Eval(node.Index, env)
            if isErrorObejct(right) {
                return right
            }

            return evalIndexExpression(left, right)
    }

    return nil
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
    var result object.Object

    for _, statement := range stmts {
        result = Eval(statement, env)
    
        switch result := result.(type) {
            case *object.ReturnObject:
                return result.Value
            case *object.ErrorObject:
                return result
        }
    }

    return result
}

func evalBlockStatement(stmts []ast.Statement, env *object.Environment) object.Object {
    var result object.Object

    for _, statement := range stmts {
        result = Eval(statement, env)
    
        if result != nil {
            rt := result.Type()
            if rt == object.RETURN_VALUE_OBJ ||
               rt == object.ERROR_OBJ {
                return result
            }
        }
    }

    return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
    val, ok := env.Get(node.Value)
    if ok {
        return val
    }

    builtin, ok := builtins[node.Value]
    if ok {
        return builtin
    }

    return newErrorObejct("identifier not found: " + node.Value)
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
    switch operator {
        case "!":
            return evalBangOperatorExpression(right)
        case "-":
            return evalMinusPrefixOperatorExpression(right)
        default:
            return newErrorObejct("unknown operator: %s%s", operator, right.Type())
    }
}

func evalInfixExpression(
    operator string, 
    left, right object.Object) object.Object {
    switch {
        case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
            return evalIntegerInfixExpression(operator, left, right)
        case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
            return evalStringInfixExpression(operator, left, right)
        case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
            return evalBooleanInfixExpression(operator, left, right)
        case left.Type() != right.Type():
            return newErrorObejct("type mismatch: %s %s %s", 
                left.Type(), operator, right.Type())
        default:
            return newErrorObejct("unknown operator: %s %s %s", 
                left.Type(), operator, right.Type())
    }
}

func evalIndexExpression(left, index object.Object) object.Object {
    switch {
        case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
            return evalArrayIndexExpression(left, index)
        default:
            return newErrorObejct("index operator not supported: %s", left.Type())
    }
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
    arrayObject := array.(*object.ArrayObject)
    idx := index.(*object.Integer).Value
    max := int64(len(arrayObject.Elements) - 1)

    if idx < 0 || idx > max {
        return NULL
    }

    return arrayObject.Elements[idx]
}

func evalBangOperatorExpression(right object.Object) object.Object {
    switch right {
        case TRUE:
            return FALSE
        case FALSE:
            return TRUE
        case NULL:
            return TRUE
        default:
            return FALSE
    }
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
    if right.Type() != object.INTEGER_OBJ {
        return newErrorObejct("unknown operator: -%s", right.Type())
    }

    value := right.(*object.Integer).Value
    return &object.Integer{Value : -value}
}

func evalIntegerInfixExpression(
    operator string, 
    left, right object.Object,
) object.Object {
    
    leftValue  := left.(*object.Integer).Value
    rightValue := right.(*object.Integer).Value
    
    switch operator {
        case "+":
            return &object.Integer{Value : leftValue + rightValue}
        case "-":
            return &object.Integer{Value : leftValue - rightValue}
        case "*":
            return &object.Integer{Value : leftValue * rightValue}
        case "/":
            return &object.Integer{Value : leftValue / rightValue}
        case "<":
            return nativaBool2BooleanObject(leftValue < rightValue)
        case ">":
            return nativaBool2BooleanObject(leftValue > rightValue)
        case "==":
            return nativaBool2BooleanObject(leftValue == rightValue)
        case "!=":
            return nativaBool2BooleanObject(leftValue != rightValue)
        default:
            return newErrorObejct("unknown operator: %s %s %s", 
                left.Type(), operator, right.Type())
    }
}

func evalStringInfixExpression (
    operator string,
    left, right object.Object,
) object.Object {
    leftValue  := left.(*object.String).Value
    rightValue := right.(*object.String).Value

    switch operator {
        case "+":
            return &object.String{Value : leftValue + rightValue}
        default:
            return newErrorObejct("unknown operator: %s %s %s", 
                left.Type(), operator, right.Type())
    }
}

func evalBooleanInfixExpression(
    operator string, 
    left, right object.Object,
) object.Object {
    
    leftValue  := left.(*object.Boolean).Value
    rightValue := right.(*object.Boolean).Value
    
    switch operator {
        case "==":
            return nativaBool2BooleanObject(leftValue == rightValue)
        case "!=":
            return nativaBool2BooleanObject(leftValue != rightValue)
        default:
            return newErrorObejct("unknown operator: %s %s %s", 
                left.Type(), operator, right.Type())
    }
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
    condition := Eval(ie.Condition, env)
    if isErrorObejct(condition) {
        return condition
    }

    if isTruthy(condition) {
        return Eval(ie.Consequence, env)
    } else if ie.Alternative != nil {
        return Eval(ie.Alternative, env)
    } else {
        return NULL
    }
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
    var result []object.Object

    for _, e := range exps {
        evaluated := Eval(e, env)
        if isErrorObejct(evaluated) {
            return []object.Object{evaluated}
        }

        result = append(result, evaluated)
    }

    return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
    switch fn := fn.(type) {
        case *object.FunctionObject:
            extendedEnv := extendFunctionEnv(fn, args)
            evaluated := Eval(fn.Body, extendedEnv)
            return unwrapReturnValue(evaluated)
        case *object.BuiltinObject:
            return fn.Fn(args...)
        default:
            return newErrorObejct("not a function: %s", fn.Type())
    }
}

func extendFunctionEnv(fn *object.FunctionObject, args []object.Object) *object.Environment {
    env := object.NewEnclosedEnvironment(fn.Env)

    for idx, param := range fn.Parameters {
        env.Set(param.Value, args[idx])
    }

    return env
}

func unwrapReturnValue(obj object.Object) object.Object {
    returnValue, ok := obj.(*object.ReturnObject)
    if ok {
        return returnValue.Value
    }

    return obj
}

func nativaBool2BooleanObject(input bool) *object.Boolean {
    if input {
        return TRUE
    }
    return FALSE
}

func isTruthy(obj object.Object) bool {
    switch obj {
        case NULL:
            return false
        case TRUE:
            return true
        case FALSE:
            return false
    }

    switch obj := obj.(type) {
        case *object.Integer:
            return isIntegerTruthy(obj)
    }

    return false
}

func isIntegerTruthy(obj *object.Integer) bool {
    intValue := obj.Value
    if intValue == 0 {
        return false
    } else {
        return true
    }
}

func newErrorObejct(format string, a ...interface{}) *object.ErrorObject {
    return &object.ErrorObject{Message : fmt.Sprintf(format, a...)}
}

func isErrorObejct(obj object.Object) bool {
    if obj != nil {
        return obj.Type() == object.ERROR_OBJ
    }
    return false
}