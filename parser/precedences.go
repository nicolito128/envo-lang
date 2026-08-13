package parser

import "github.com/nicolito128/envo-lang/lexer"

const (
	LOWEST      int = iota
	COMP            // &&, ||
	EQUALS          // ==, !=
	LESSGREATER     // <, >, <=, >=
	SUM             // +, -
	PRODUCT         // *, /
	PREFIX          // -x
	CALL            // msg{}
)

var precedences = map[lexer.TokenKind]int{
	lexer.AND:    COMP,
	lexer.OR:     COMP,
	lexer.EQL:    EQUALS,
	lexer.NEQ:    EQUALS,
	lexer.LSS:    LESSGREATER,
	lexer.LEQ:    LESSGREATER,
	lexer.GTR:    LESSGREATER,
	lexer.GEQ:    LESSGREATER,
	lexer.ADD:    SUM,
	lexer.SUB:    SUM,
	lexer.MUL:    PRODUCT,
	lexer.DIV:    PRODUCT,
	lexer.REM:    PRODUCT,
	lexer.LBRACE: CALL,
}

func (p *Parser) curPrecedence() int {
	if prec, ok := precedences[p.curToken.Tok.Kind]; ok {
		return prec
	}
	return LOWEST
}

// peekPrecedence returns the precedence of the next token (peekToken).
func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.peekToken.Tok.Kind]; ok {
		return prec
	}
	return LOWEST
}
