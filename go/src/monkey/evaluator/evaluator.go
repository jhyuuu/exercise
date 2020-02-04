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

func Eval(node ast.Node) object.Object {
    switch node := node.(type) {
        case *ast.Program:
            return evalProgram(node.Statements)
        case *ast.BlockStatement:
            return evalBlockStatement(node.Statements)
        case *ast.ReturnStatement:
            val := Eval(node.ReturnValue)
            if isError(val) {
                return val
            }

            return &object.ReturnObject{Value : val}
        case *ast.ExpressionStatement:
            return Eval(node.Expression)
        case *ast.IntegerLiteral:
            return &object.Integer{Value : node.Value}
        case *ast.Boolean:
            return nativaBool2BooleanObject(node.Value)
        case *ast.PrefixExpression:
            right := Eval(node.Right)
            if isError(right) {
                return right
            }

            return evalPrefixExpression(node.Operator, right)
        case *ast.InfixExpression:
            left  := Eval(node.Left)
            if isError(left) {
                return left
            }

            right := Eval(node.Right)
            if isError(right) {
                return right
            }

            return evalInfixExpression(node.Operator, left, right)
        case *ast.IfExpression:
            return evalIfExpression(node)
        
    }

    return nil
}

func evalProgram(stmts []ast.Statement) object.Object {
    var result object.Object

    for _, statement := range stmts {
        result = Eval(statement)
    
        // returnValue, ok := result.(*object.ReturnObject)
        // if ok {
        //     return returnValue.Value
        // }

        switch result := result.(type) {
            case *object.ReturnObject:
                return result.Value
            case *object.ErrorObject:
                return result
        }
    }

    return result
}

func evalBlockStatement(stmts []ast.Statement) object.Object {
    var result object.Object

    for _, statement := range stmts {
        result = Eval(statement)
    
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

func nativaBool2BooleanObject(input bool) *object.Boolean {
    if input {
        return TRUE
    }
    return FALSE
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
    switch operator {
        case "!":
            return evalBangOperatorExpression(right)
        case "-":
            return evalMinusPrefixOperatorExpression(right)
        default:
            return newError("unknown operator: %s%s", operator, right.Type())
    }
}

func evalInfixExpression(
    operator string, 
    left, right object.Object) object.Object {
    switch {
        case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
            return evalIntegerInfixExpression(operator, left, right)
        case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
            return evalBooleanInfixExpression(operator, left, right)
        case left.Type() != right.Type():
            return newError("type mismatch: %s %s %s", 
                left.Type(), operator, right.Type())
        default:
            return newError("unknown operator: %s %s %s", 
                left.Type(), operator, right.Type())
    }
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
        return newError("unknown operator: -%s", right.Type())
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
            return newError("unknown operator: %s %s %s", 
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
            return newError("unknown operator: %s %s %s", 
                left.Type(), operator, right.Type())
    }
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
    condition := Eval(ie.Condition)
    if isError(condition) {
        return condition
    }

    if isTruthy(condition) {
        return Eval(ie.Consequence)
    } else if ie.Alternative != nil {
        return Eval(ie.Alternative)
    } else {
        return NULL
    }
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

func newError(format string, a ...interface{}) *object.ErrorObject {
    return &object.ErrorObject{Message : fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
    if obj != nil {
        return obj.Type() == object.ERROR_OBJ
    }
    return false
}