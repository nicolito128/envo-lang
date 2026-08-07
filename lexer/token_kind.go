package lexer

type TokenKind int

func (tk TokenKind) String() string {
	s := tokens[UNKNOWN]
	if 0 <= tk && tk < TokenKind(len(tokens)) {
		s = tokens[tk]
	}
	return s
}

const (
	UNKNOWN TokenKind = iota

	EOF
	NEWLN
	COMMENT

	literalBegin
	// Identifiers and type literals
	IDENT  // main
	CHAR   // 'a'
	SYMBOL // &ok
	STRING // "abc"
	RAWSTR // `\txyz\n`
	INT    // 1234
	FLOAT  // 0.12345
	IMAG   // 1.23i
	BOOL   // true, false
	literalEnd

	operatorBegin
	// Operators
	ADD // +
	SUB // -
	MUL // *
	DIV // /
	REM // %

	EQL // ==
	NEQ // !=
	LSS // <
	LEQ // <=
	GTR // >
	GEQ // >=

	AND // &&
	OR  // ||
	NOT // !

	ELLIPSIS // ...

	LPAREN // (
	LBRACK // [
	LBRACE // {

	RPAREN // )
	RBRACK // ]
	RBRACE // }

	COMMA  // ,
	PERIOD // .

	COLON     // :
	SEMICOLON // ;
	operatorEnd
)

func (tk TokenKind) IsLiteral() bool  { return literalBegin < tk && tk < literalEnd }
func (tk TokenKind) IsOperator() bool { return operatorBegin < tk && tk < operatorEnd }

var tokens = [...]string{
	UNKNOWN: "Unknown",

	EOF:     "EOF",
	NEWLN:   "Newln",
	COMMENT: "Comment",

	IDENT:  "Ident",
	CHAR:   "Char",
	SYMBOL: "Symbol",
	STRING: "String",
	RAWSTR: "RawString",
	INT:    "Int",
	FLOAT:  "Float",
	IMAG:   "Imag",

	ADD:       "+",
	SUB:       "-",
	MUL:       "*",
	DIV:       "/",
	REM:       "%",
	EQL:       "==",
	NEQ:       "!=",
	LSS:       "<",
	LEQ:       "<=",
	GTR:       ">",
	GEQ:       ">=",
	AND:       "&&",
	OR:        "||",
	NOT:       "!",
	ELLIPSIS:  "...",
	LPAREN:    "(",
	LBRACK:    "[",
	LBRACE:    "{",
	RPAREN:    ")",
	RBRACK:    "]",
	RBRACE:    "}",
	COMMA:     ",",
	PERIOD:    ".",
	COLON:     ":",
	SEMICOLON: ";",
}

var (
	operators map[string]TokenKind
)

func init() {
	operators = make(map[string]TokenKind)
	for i := operatorBegin + 1; i < operatorEnd; i++ {
		operators[tokens[i]] = i
	}
}

func LookupOperator(s string) TokenKind {
	op, ok := operators[s]
	if !ok {
		return UNKNOWN
	}
	return op
}

const (
	Newln rune = '\n'

	CommentPrefix rune = '#'
	SymbolPrefix  rune = '~'

	CharacterMark rune = '\''
	StringMark    rune = '"'
	RawStringMark rune = '`'

	ImaginaryPart rune = 'i'
)
