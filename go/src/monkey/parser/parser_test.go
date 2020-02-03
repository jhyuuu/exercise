package parser

import (
    "fmt"
    "testing"

    "monkey/ast"
    "monkey/lexer"
)

func TestLetStatements (t *testing.T) {
    input :=    `
let x = 5;
let y = 10;
let foobar = 838383;
`
    var l *lexer.Lexer = lexer.New(input)
    var p *Parser       = New(l)

    var program *ast.Program = p.ParseProgram()
    checkParseErrors(t, p)

    if program == nil {
        t.Fatalf("ParseProgram() returned nil")
    }

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements does not contain 3 statements. got=%d",
             len(program.Statements))
    }

    tests := []struct {
         expectedIdentifier string
    } {
        {"x"},
        {"y"},
        {"foobar"},
    }

    for i, tt := range tests {
        stmt := program.Statements[i]
        if !testLetStatement(t, stmt, tt.expectedIdentifier) {
            return
        }
    }
}

func TestReturnStatements(t *testing.T) {
    input := `
return 5;
return 10;
return 993322;
`
    l := lexer.New(input)
    p := New(l)

    program := p.ParseProgram()
    checkParseErrors(t, p)

    if len(program.Statements) != 3 {
        t.Fatalf("program.Statements does not contain 3 statements. got=%d",
             len(program.Statements))
    }

    for _, stmt := range program.Statements {
        returnStmt, ok := stmt.(*ast.ReturnStatement)
        if !ok {
            t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
            continue
        }

        if returnStmt.TokenLiteral() != "return" {
            t.Errorf("returnStmt.TokenLiteral not 'return', got %q",
                returnStmt.TokenLiteral())
        }
    }
}

func TestIdentifierExpression(t *testing.T) {
    input := "foobar;"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParseErrors(t, p)
    
    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d",
            len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
            program.Statements[0])
    }

    ident, ok := stmt.Expression.(*ast.Identifier)
    if !ok {
        t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
    }

    if ident.Value != "foobar" {
        t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
    }

    if ident.TokenLiteral() != "foobar" {
        t.Errorf("ident.TokenLiteral not %s. got %s", "foobar", ident.TokenLiteral())
    }
}

func TestIntegerLiteralExpression(t *testing.T) {
    input := "5;"

    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParseErrors(t, p)
    
    if len(program.Statements) != 1 {
        t.Fatalf("program has not enough statements. got=%d",
            len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
            program.Statements[0])
    }

    literal, ok := stmt.Expression.(*ast.IntegerLiteral)
    if !ok {
        t.Fatalf("expression not *ast.IntegerLiteral. got=%T", stmt.Expression)
    }

    if literal.Value != 5 {
        t.Errorf("literal.Value not %s. got=%q", "5", literal.Value)
    }

    if literal.TokenLiteral() != "5" {
        t.Errorf("literal.TokenLiteral not %s. got %s", "5", literal.TokenLiteral())
    }
}

func TestBoolLiteralExpression(t *testing.T) {
    tests := []struct {
        input           string
        expectedBoolean bool
    } {
        {"true;", true},
        {"false;", false},
    }

    for _, tt := range tests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParseErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program has not enough statements. got=%d",
                len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                program.Statements[0])
        }

        boolean, ok := stmt.Expression.(*ast.Boolean)
        if !ok {
            t.Fatalf("expression not *ast.Boolean. got=%T", stmt.Expression)
        }

        if boolean.Value != tt.expectedBoolean {
            t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
                boolean.Value)
        }
    }
}

func TestParsingPrefixExpression(t *testing.T) {
    prefixTests := []struct {
        input        string
        operator     string
        value        interface{}
    } {
        {"!5", "!", 5},
        {"-15", "-", 15},
        {"!foobar;", "!", "foobar"},
        {"-foobar;", "-", "foobar"},
        {"!true;", "!", true},
        {"!false;", "!", false},
    }

    for _, tt := range prefixTests {
        l := lexer.New(tt.input)
        p := New(l)

        program := p.ParseProgram()
        checkParseErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program has not enough statements. got=%d",
                len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                program.Statements[0])
        }

        if !testPrefixExpression(t, stmt.Expression, tt.operator, tt.value) {
            return
        }
    }
}

func TestParsingInfixExpression(t *testing.T) {
     infixTests := []struct {
         input      string
         leftValue  interface{}
         operator   string
         rightValue interface{}
     } {
        {"5 + 5;", 5, "+", 5},
        {"5 - 5;", 5, "-", 5},
        {"5 * 5;", 5, "*", 5},
        {"5 / 5;", 5, "/", 5},
        {"5 > 5;", 5, ">", 5},
        {"5 < 5;", 5, "<", 5},
        {"5 == 5;", 5, "==", 5},
        {"5 != 5;", 5, "!=", 5},
        {"foobar + barfoo;", "foobar", "+", "barfoo"},
        {"foobar - barfoo;", "foobar", "-", "barfoo"},
        {"foobar * barfoo;", "foobar", "*", "barfoo"},
        {"foobar / barfoo;", "foobar", "/", "barfoo"},
        {"foobar > barfoo;", "foobar", ">", "barfoo"},
        {"foobar < barfoo;", "foobar", "<", "barfoo"},
        {"foobar == barfoo;", "foobar", "==", "barfoo"},
        {"foobar != barfoo;", "foobar", "!=", "barfoo"},
        {"true == true", true, "==", true},
        {"true != false", true, "!=", false},
        {"false == false", false, "==", false},
     }

     for _, tt := range infixTests {
         l := lexer.New(tt.input)
         p := New(l)

         program := p.ParseProgram()
         checkParseErrors(t, p)

        if len(program.Statements) != 1 {
            t.Fatalf("program has not enough statements. got=%d",
                len(program.Statements))
        }

        stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
        if !ok {
            t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
                program.Statements[0])
        }

        if !testInfixExpression(t, stmt.Expression, tt.leftValue,
            tt.operator, tt.rightValue) {
            return
        }
    }
}

func TestOperatorPrecedenceParsing(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {
            "-a * b",
            "((-a) * b)",
        },
        {
            "!-a",
            "(!(-a))",
        },
        {
            "a + b + c",
            "((a + b) + c)",
        },
        {
            "a + b - c",
            "((a + b) - c)",
        },
        {
            "a * b * c",
            "((a * b) * c)",
        },
        {
            "a * b / c",
            "((a * b) / c)",
        },
        {
            "a + b / c",
            "(a + (b / c))",
        },
        {
            "a + b * c + d / e - f",
            "(((a + (b * c)) + (d / e)) - f)",
        },
        {
            "3 + 4; -5 * 5",
            "(3 + 4)((-5) * 5)",
        },
        {
            "5 > 4 == 3 < 4",
            "((5 > 4) == (3 < 4))",
        },
        {
            "5 < 4 != 3 > 4",
            "((5 < 4) != (3 > 4))",
        },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
         },
        {
            "3 + 4 * 5 == 3 * 1 + 4 * 5",
            "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
        },
        {
            "true",
            "true",
        },
        {
            "false",
            "false",
        },
        {
            "3 > 5 == false",
            "((3 > 5) == false)",
        },
        {
            "3 < 5 == true",
            "((3 < 5) == true)",
        },
        {
            "1 + (2 + 3) + 4",
            "((1 + (2 + 3)) + 4)",
        },
        {
            "(5 + 5) * 2",
            "((5 + 5) * 2)",
        },
        {
            "2 / (5 + 5)",
            "(2 / (5 + 5))",
        },
        {
            "(5 + 5) * 2 * (5 + 5)",
            "(((5 + 5) * 2) * (5 + 5))",
        },
        {
            "-(5 + 5)",
            "(-(5 + 5))",
        },
        {
            "!(true == true)",
            "(!(true == true))",
        },
    }

    for _, tt := range tests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParseErrors(t, p)

        actual := program.String()
        if actual != tt.expected {
            t.Errorf("expected=%q, got=%q", tt.expected, actual)
        }
    }
}

// directory: ~/exercise/go/src/monkey
// cmd: go test -v -run TestOperatorPrecedenceParsing_verbose ./parser/
func TestOperatorPrecedenceParsing_verbose(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {
            "-1 + 2 * 3",
            "((-1) + (2 * 3))",
        },
    }

    for _, tt := range tests {
        l := lexer.New(tt.input)
        p := New(l)
        program := p.ParseProgram()
        checkParseErrors(t, p)

        actual := program.String()
        if actual != tt.expected {
            t.Errorf("expected=%q, got=%q", tt.expected, actual)
        }
    }
}

func TestIfExpression(t *testing.T) {
    input := `if (x < y) { x }`
    
    l := lexer.New(input)
    p := New(l)
    program := p.ParseProgram()
    checkParseErrors(t, p)

    if len(program.Statements) != 1 {
        t.Fatalf("program.Body does not contain %d statements. got=%d\n",
            1, len(program.Statements))
    }

    stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
        program.Statements[0])
    }

    exp, ok := stmt.Expression.(*ast.IfExpression)
    if !ok {
        t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
            stmt.Expression)
    }

    if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
        return
    }

    if len(exp.Consequence.Statements) != 1 {
        t.Errorf("consequence is not 1 statements. got=%d\n",
            len(exp.Consequence.Statements))
    }

    consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
    if !ok {
        t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
            exp.Consequence.Statements[0])
    }

    if !testIdentifier(t, consequence.Expression, "x") {
        return
    }

    if exp.Alternative != nil {
        t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
    }
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
    integer, ok := il.(*ast.IntegerLiteral)
    if !ok {
        t.Errorf("il is not *ast.IntegerLitetal. got %T", il)
        return false
    }

    if integer.Value != value {
        t.Errorf("integer.Value not %d. got %d", value, integer.Value)
        return false
    }

    if integer.TokenLiteral() != fmt.Sprintf("%d", value) {
        t.Errorf("integer.TokenLiteral is not %d. got %s", value, integer.TokenLiteral())
        return false
    }
    return true
}

func testLiteralExpression(
    t *testing.T,
    exp ast.Expression,
    expected interface{},
) bool {
    switch v := expected.(type) {
        case int:
            return testIntegerLiteral(t, exp, int64(v))
        case int64:
            return testIntegerLiteral(t, exp, v)
        case string:
            return testIdentifier(t, exp, v)
        case bool:
            return testBooleanLiteral(t, exp, v)
    }
    t.Errorf("type of exp not handled. got=%T", exp)
    return false
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
    ident, ok := exp.(*ast.Identifier)
    if !ok {
        t.Errorf("exp not *ast.Identifier. got=%T", exp)
        return false
    }

    if ident.Value != value {
        t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
        return false
    }

    if ident.TokenLiteral() != value {
        t.Errorf("ident.TokenLiteral not %s. got=%s", value,
        ident.TokenLiteral())
        return false
    }

    return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
    bo, ok := exp.(*ast.Boolean)
    if !ok {
        t.Errorf("exp not *ast.Boolean. got=%T", exp)
        return false
    }

    if bo.Value != value {
        t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
        return false
    }

    if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
        t.Errorf("bo.TokenLiteral not %t. got=%s",
            value, bo.TokenLiteral())
        return false
    }

    return true
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
    if s.TokenLiteral() != "let" {
        t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
        return false   
    }

    letStmt, ok := s.(*ast.LetStatement) // 类型断言
    if !ok {
        t.Errorf("s not *ast.LetStatement. got=%T", s)
        return false
    }

    if letStmt.Name.Value != name {
        t.Errorf("letStmt.Name.Value not '%s'. got=%s",    name, letStmt.Name.Value)
        return false
    }

    if letStmt.Name.TokenLiteral() != name {
        t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
        return false
    }

    return true
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, value interface{}) bool {
    opExp, ok := exp.(*ast.PrefixExpression)
    if !ok {
        t.Fatalf("stmt is not *ast.PrefixExpression. got=%T(%s)", exp, exp)
    }

    if opExp.Operator != operator {
        t.Fatalf("exp.Operator is not '%s'. got=%s",
            operator, opExp.Operator)
    }

    if !testLiteralExpression(t, opExp.Right, value) {
        return false
    }
    
    return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
    operator string, right interface{}) bool {

    opExp, ok := exp.(*ast.InfixExpression)
    if !ok {
        t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
        return false
    }

    if !testLiteralExpression(t, opExp.Left, left) {
        return false
    }

    if opExp.Operator != operator {
        t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
        return false
    }

    if !testLiteralExpression(t, opExp.Right, right) {
        return false
    }

    return true
}

func checkParseErrors(t *testing.T, p *Parser) {
    errors := p.errors
    if len(errors) == 0 {
        return
    }

    t.Errorf("parser has %d errors", len(errors))
    for _, msg := range errors {
        t.Errorf("parser error: %q", msg)
    }
    t.FailNow()
}