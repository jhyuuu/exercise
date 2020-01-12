package parser

import (
    "fmt"
    "strconv"

    "monkey/token"
    "monkey/lexer"
    "monkey/ast"
)

const (
    _ int = iota
    LOWEST
    EQUALS
    LESSGREATER
    SUM
    PRODUCT
    PREFIX
    CALL
)

type (
    prefixParseFn func() ast.Expression
    infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
    l *lexer.Lexer
    errors []string

    curTok  token.Token
    peekTok token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l* lexer.Lexer) *Parser {
    var p *Parser = &Parser{
        l : l,
        errors : []string {},
    }

    p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
    p.registerPrefix(token.IDENT, p.parseIdentifier)
    p.registerPrefix(token.INT, p.parseIntegerLiteral)

    p.NextToken()
    p.NextToken()

    return p
}

func (p *Parser) Errors() []string {
    return p.errors
}

func (p *Parser) NextToken() {
    p.curTok = p.peekTok
    p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
    var program *ast.Program = &ast.Program{}

    for !p.curTokenIs(token.EOF) {
        stmt := p.ParseStatement()
        if stmt != nil {
            program.Statements = append(program.Statements, stmt)
        }
        p.NextToken()
    }
    return program
}

func (p *Parser) ParseStatement() ast.Statement {
    switch p.curTok.Type {
    case token.LET:    return p.ParseLetStatement()
    case token.RETURN: return p.ParseReturnStatement()
    default:           return p.ParseExpressionStatement()
    }
}

func (p *Parser) ParseLetStatement() *ast.LetStatement {
    stmt := &ast.LetStatement{Token : p.curTok}

    if !p.expectPeek(token.IDENT) {
        return nil
    }

    stmt.Name = &ast.Identifier {Token:p.curTok, Value:p.curTok.Literal}

    for !p.curTokenIs(token.SEMICOLON) {
        p.NextToken()
    }

    return stmt
}

func (p *Parser) ParseReturnStatement() *ast.ReturnStatement {
    stmt := &ast.ReturnStatement{Token : p.curTok}

    p.NextToken()

    for !p.curTokenIs(token.SEMICOLON) {
        p.NextToken()
    }

    return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
    return p.curTok.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
    return p.peekTok.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
    if p.peekTokenIs(t) {
        p.NextToken()
        return true
    } else {
        p.peekError(t)
        return false
    }
}

func (p *Parser) peekError(t token.TokenType) {
    msg := fmt.Sprintf("expected next token to be %s, got %s instead",
        t, p.peekTok.Type)
    p.errors = append(p.errors, msg)
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
    p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix (tokenType token.TokenType, fn infixParseFn) {
    p.infixParseFns[tokenType] = fn
}

func (p *Parser) ParseExpressionStatement() *ast.ExpressionStatement {
    stmt := &ast.ExpressionStatement{Token : p.curTok}

    stmt.Expression = p.parseExpression(LOWEST)

    if p.peekTokenIs(token.SEMICOLON) {
        p.NextToken()
    }

    return stmt
}

func (p *Parser) parseExpression(precedenct int) ast.Expression {
    prefixfn := p.prefixParseFns[p.curTok.Type]
    if prefixfn == nil {
        return nil
    }
    leftExp := prefixfn()
    return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
    return &ast.Identifier{Token : p.curTok, Value : p.curTok.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
    it := &ast.IntegerLiteral{Token : p.curTok}

    value, err := strconv.ParseInt(p.curTok.Literal, 0, 64)
    if err != nil {
        msg := fmt.Sprintf("could not parse %q as integer", p.curTok.Literal)
        p.errors = append(p.errors, msg)
        return nil
    }

    it.Value = value
    return it
}