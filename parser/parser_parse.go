package parser

import (
	"fmt"
	"strconv"

	"github.com/nicolito128/envo-lang/ast"
	"github.com/nicolito128/envo-lang/lexer"
)

func (p *Parser) Parse() (ast.Node, []ParseError) {
	return p.parseProgram(), p.Errors()
}

func (p *Parser) parseProgram() *ast.Program {
	node := &ast.Program{
		Statements: make([]ast.Node, 0),
	}

	for !p.curIs(lexer.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			node.Statements = append(node.Statements, stmt)
		}
		p.nextToken()
	}

	return node
}

func (p *Parser) parseStatement() ast.Node {
	kind := p.curToken.Tok.Kind
	switch kind {
	case lexer.IDENT:
		return p.parseIdent()
	default:
		return p.parseExpression(LOWEST)
	}
}

func (p *Parser) parseIdent() ast.Node {
	if p.curToken.Tok.Kind != lexer.IDENT {
		p.errors = append(p.errors, ParseError{
			Pos:     p.curToken.Pos,
			Message: "expected an identifier",
		})
		return nil
	}

	var node ast.Node

	ident := &ast.Identifier{Name: p.curToken.Tok.Literal}
	node = ident

	var pattern *ast.RawMessage
	if p.peekIs(lexer.LBRACE) {
		p.nextToken()

		pattern = p.parseRawMessage()
		if pattern != nil {
			node = &ast.Message{
				Label: ident,
				Body:  pattern,
			}
		}

		if !p.peekIs(lexer.LBRACE) || p.peekIs(lexer.EOF) {
			return &ast.Message{
				Label: ident,
				Body:  pattern,
			}
		}
	}

	var response *ast.RawMessage
	if p.peekIs(lexer.LBRACE) {
		p.nextToken()

		response = p.parseRawMessage()
		if response != nil {
			node = &ast.Define{
				Receiver: ident,
				Pattern:  pattern,
				Response: response,
			}
		}
	}

	return node
}

func (p *Parser) parseRawMessage() *ast.RawMessage {
	if !p.curIs(lexer.LBRACE) {
		p.errors = append(p.errors, ParseError{
			Pos:     p.curToken.Pos,
			Message: "expected '{' to start the raw message",
		})
		return nil
	}

	msg := &ast.RawMessage{
		Words: make([]ast.Node, 0),
	}

	hasPrevWord := false
	for {
		p.nextToken()
		if p.curIs(lexer.RBRACE) || p.curIs(lexer.EOF) {
			break
		}

		if p.curIs(lexer.COMMA) {
			if !hasPrevWord {
				p.errors = append(p.errors, ParseError{
					Pos:     p.curToken.Pos,
					Message: "unexpected ',' in raw message",
				})
			}
			hasPrevWord = false

			continue
		}

		word := p.parseExpression(LOWEST)
		if word != nil {
			msg.Words = append(msg.Words, word)
			hasPrevWord = true
		}
	}

	return msg
}

func (p *Parser) parseExpression(precedence int) ast.Node {
	var node ast.Node

	kind := p.curToken.Tok.Kind
	lit := p.curToken.Tok.Literal

	switch kind {
	case lexer.IDENT:
		node = p.parseIdent()
	case lexer.SYMBOL:
		node = &ast.SymbolLiteral{Name: lit}
	case lexer.STRING:
		node = &ast.StringLiteral{Value: fmt.Sprintf("%s", lit)}
	case lexer.RAWSTR:
		node = &ast.RawStringLiteral{Value: lit}
	case lexer.BOOL:
		if lit == "true" {
			node = &ast.BoolLiteral{Value: true}
		} else {
			node = &ast.BoolLiteral{Value: false}
		}
	case lexer.CHAR:
		node = &ast.CharLiteral{Value: []rune(lit)}
	case lexer.INT:
		if num, err := strconv.ParseInt(lit, 0, 64); err == nil {
			node = &ast.IntegerLiteral{Value: num}
		}
	case lexer.FLOAT:
		if num, err := strconv.ParseFloat(lit, 64); err == nil {
			node = &ast.FloatLiteral{Value: num}
		}
	case lexer.IMAG:
		if num, err := strconv.ParseComplex(lit, 128); err == nil {
			node = &ast.ImaginaryLiteral{Value: num}
		}
	case lexer.LBRACE:
		node = p.parseRawMessage()
	case lexer.LPAREN:
		p.nextToken() // (
		node = p.parseExpression(LOWEST)
		if !p.expect(lexer.RPAREN) { // )
			return nil
		}
	case lexer.NEWLN:
		return nil
	case lexer.SEMICOLON:
		return nil
	case lexer.COMMENT:
		return nil
	default:
		if kind.IsOperator() {
			node = p.parseOperator()
		} else {
			p.errors = append(p.errors, ParseError{
				Pos:     p.curToken.Pos,
				Message: "unexpected token in expression",
			})
		}
	}

	guard := !p.peekIs(lexer.EOF) && !p.peekIs(lexer.RBRACE) && !p.peekIs(lexer.COMMA)
	for guard && precedence < p.peekPrecedence() {
		p.nextToken()
		node = p.parseBinaryExpr(node)
	}

	if node != nil && p.curIs(lexer.RBRACE) {
		return node
	}

	return node
}

func (p *Parser) parseOperator() ast.Node {
	kind := p.curToken.Tok.Kind
	lit := p.curToken.Tok.Literal

	if !kind.IsOperator() {
		return nil
	}
	p.nextToken()

	var node ast.Node

	operand := p.parseExpression(PREFIX)
	if operand == nil {
		p.errors = append(p.errors, ParseError{
			Pos:     p.curToken.Pos,
			Message: "expected operand after prefix operator",
		})
	}

	node = &ast.UnaryOp{
		Op:      lit,
		Operand: operand,
	}

	return node
}

func (p *Parser) parseBinaryExpr(left ast.Node) ast.Node {
	lit := p.curToken.Tok.Literal

	curWeight := p.curPrecedence()
	p.nextToken()

	right := p.parseExpression(curWeight)
	if right == nil {
		p.errors = append(p.errors, ParseError{
			Pos:     p.curToken.Pos,
			Message: "expected right-hand side of operator",
		})
	}

	return &ast.BinaryOp{
		Op:    lit,
		Left:  left,
		Right: right,
	}
}
