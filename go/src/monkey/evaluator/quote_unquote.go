package evaluator

import (
    "fmt"

    "monkey/token"
    "monkey/ast"
    "monkey/object"
)

func quote(node ast.Node, env *object.Environment) object.Object {
    node = evalUnquoteCalls(node, env)
    return &object.Quote{Node : node}
}

func evalUnquoteCalls(node ast.Node, env *object.Environment) ast.Node {
    modifer := func(node ast.Node) ast.Node {
        if !isUnquoteCall(node) {
            return node
        }

        call, ok := node.(*ast.CallExpression)
        if !ok {
            return node
        }

        if len(call.Arguments) != 1 {
            return node
        }

        var unquoted object.Object = Eval(call.Arguments[0], env)
        return convertObjectToAstNode(unquoted)
    }

    return ast.Modify(node, modifer)
}

func isUnquoteCall(node ast.Node) bool {
    call, ok := node.(*ast.CallExpression)
    if !ok {
        return false
    }

    return call.Function.TokenLiteral() == "unquote"
}

func convertObjectToAstNode(obj object.Object) ast.Node {
    switch obj := obj.(type) {
        case *object.Integer:
            t := token.Token {
                Type:    token.INT,
                Literal: fmt.Sprintf("%d", obj.Value),
            }
            return &ast.IntegerLiteral{Token : t, Value : obj.Value}
        case *object.Boolean:
            var t token.Token
            if obj.Value {
                t = token.Token {Type : token.TRUE, Literal : "true"}
            } else {
                t = token.Token {Type : token.FALSE, Literal : "false"}
            }
            return &ast.Boolean{Token : t, Value : obj.Value}
        case *object.Quote:
            return obj.Node
        default:
            return nil
    }
}