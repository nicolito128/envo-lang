package parser

import (
	"fmt"

	"n128.xyz/n128/envo/lexer"
)

type Token struct {
	Tok lexer.Token
	Pos lexer.Position
}

type ParseError struct {
	Pos     lexer.Position
	Message string
}

func (e ParseError) Error() string {
	return fmt.Sprintf("[%d:%d] syntax error: %s", e.Pos.Line(), e.Pos.Column(), e.Message)
}

type Lexer interface {
	Lex() (lexer.Token, lexer.Position, error)
}

type Parser struct {
	lexer Lexer

	curToken  Token
	peekToken Token

	errors []ParseError
}

func New(l Lexer) *Parser {
	p := &Parser{lexer: l}
	p.errors = make([]ParseError, 0)

	p.nextToken()
	p.nextToken()

	return p
}

// nextToken advances the parser to the next token in the input stream.
func (p *Parser) nextToken() {
	tok, pos, err := p.lexer.Lex()
	if err != nil {
		p.errors = append(p.errors, ParseError{
			Pos:     pos,
			Message: err.Error(),
		})
	}
	p.curToken = p.peekToken

	p.peekToken = Token{
		Tok: tok,
		Pos: pos,
	}
}

func (p *Parser) Errors() []ParseError {
	return p.errors
}

// curIs checks if the current token matches the given kind.
func (p *Parser) curIs(kind lexer.TokenKind) bool {
	return p.curToken.Tok.Kind == kind
}

// peekIs checks if the next token matches the given kind.
func (p *Parser) peekIs(kind lexer.TokenKind) bool {
	return p.peekToken.Tok.Kind == kind
}

// expect checks if the next token matches the given kind
// and advances the parser if it does. If not, it records an error.
func (p *Parser) expect(kind lexer.TokenKind) bool {
	if p.peekIs(kind) {
		p.nextToken()
		return true
	}
	p.peekError(kind)
	return false
}

// peekError records an error indicating that the next token does not match the expected kind.
func (p *Parser) peekError(kind lexer.TokenKind) {
	msg := fmt.Sprintf("expected '%v', but got '%v'", kind, p.peekToken.Tok.Kind)
	p.errors = append(p.errors, ParseError{Pos: p.peekToken.Pos, Message: msg})
}
