package lexer

import (
	"bufio"
	"io"
)

type Lexer struct {
	reader *bufio.Reader
	pos    Position
	lastr  rune
}

func New(r io.Reader) *Lexer {
	return NewLexer(r)
}

func NewLexer(r io.Reader) *Lexer {
	lx := new(Lexer)
	lx.pos = Position{line: 1, column: 0}
	lx.reader = bufio.NewReader(r)
	return lx
}

func (lx *Lexer) advanceLine() {
	lx.pos.line++
	lx.pos.column = 0
}

func (lx *Lexer) scanRune() (rune, error) {
	r, _, err := lx.reader.ReadRune()
	if err != nil {
		return 0, err
	}

	lx.lastr = r

	if r == Newln {
		lx.advanceLine()
	} else {
		lx.pos.column++
	}

	return r, nil
}

func (lx *Lexer) peekRune() (rune, error) {
	r, _, err := lx.reader.ReadRune()
	if err != nil {
		return 0, err
	}

	err = lx.reader.UnreadRune()
	return r, err
}
