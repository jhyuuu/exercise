package parser

import (
	"monkey/token"
    "monkey/lexer"
	"monkey/ast"
)

type Parser struct {
    l *lexer.Lexer

	curTok  token.Token
	peekTok token.Token
}

func New(l* lexer.Lexer) *Parser {
    var p *Parser = &Parser{l:l}

	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) NextToken() {
    p.curTok = p.peekTok
	p.peekTok = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
    var program *ast.Program = &ast.Program{}
	// program.Statements = []ast.Statement{}

	for p.curTok.Type != token.EOF {
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
	case token.LET: return p.ParseLetStatement()
	default:        return nil
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
		return false
	}
}