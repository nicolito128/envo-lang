package lexer

import "fmt"

type Token struct {
	Kind    TokenKind
	Literal string
}

func (tok Token) String() string {
	return fmt.Sprintf("Token{%s, %s}", tok.Kind, tok.Literal)
}

// Common tokens
var (
	tokUNKNOWN = Token{UNKNOWN, "Unknown"}

	tokEOF     = Token{EOF, "EOF"}
	tokNEWLN   = Token{NEWLN, "ln"}
	tokCOMMENT = Token{COMMENT, "Comment"}

	tokADD = Token{ADD, "+"}
	tokSUB = Token{SUB, "-"}
	tokMUL = Token{MUL, "*"}
	tokDIV = Token{DIV, "/"}
	tokREM = Token{REM, "%"}

	tokEQL = Token{EQL, "=="}
	tokNEQ = Token{NEQ, "!="}
	tokLSS = Token{LSS, "<"}
	tokLEQ = Token{LEQ, "<="}
	tokGTR = Token{GTR, ">"}
	tokGEQ = Token{GEQ, ">="}

	tokAND = Token{AND, "&&"}
	tokOR  = Token{OR, "||"}
	tokNOT = Token{NOT, "!"}

	tokELLIPSIS = Token{ELLIPSIS, "..."}

	tokLPAREN    = Token{LPAREN, "("}
	tokLBRACK    = Token{LBRACK, "["}
	tokLBRACE    = Token{LBRACE, "{"}
	tokRPAREN    = Token{RPAREN, ")"}
	tokRBRACK    = Token{RBRACK, "]"}
	tokRBRACE    = Token{RBRACE, "}"}
	tokCOMMA     = Token{COMMA, ","}
	tokPERIOD    = Token{PERIOD, "."}
	tokCOLON     = Token{COLON, ":"}
	tokSEMICOLON = Token{SEMICOLON, ";"}
)
