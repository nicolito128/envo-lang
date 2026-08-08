package lexer

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

const (
	DefaultSymbolAlloc = 32
	DefaultRawStrAlloc = 32
	DefaultStrAlloc    = 16
	DefaultCharAlloc   = 2
)

func (lx *Lexer) Lex() (Token, Position, error) {
	for {
		ch, err := lx.scanRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return tokEOF, lx.pos, nil
			}

			return tokUNKNOWN, lx.pos, err
		}

		switch ch {
		case Newln:
			return lx.lexNewline()
		case CommentPrefix:
			return lx.lexComment()
		case SymbolPrefix:
			return lx.lexSymbol()
		case CharacterMark:
			return lx.lexCharacter()
		case StringMark:
			return lx.lexString()
		case RawStringMark:
			return lx.lexRawString()
		case '+':
			return tokADD, lx.pos, nil
		case '*':
			return tokMUL, lx.pos, nil
		case '/':
			return tokDIV, lx.pos, nil
		case '%':
			return tokREM, lx.pos, nil
		case '(':
			return tokLPAREN, lx.pos, nil
		case ')':
			return tokRPAREN, lx.pos, nil
		case '[':
			return tokLBRACK, lx.pos, nil
		case ']':
			return tokRBRACK, lx.pos, nil
		case '{':
			return tokLBRACE, lx.pos, nil
		case '}':
			return tokRBRACE, lx.pos, nil
		case ',':
			return tokCOMMA, lx.pos, nil
		case ':':
			return tokCOLON, lx.pos, nil
		case '.':
			return lx.lexPeriod()
		case ';':
			return lx.lexSemicolon()
		case '!':
			return lx.lexBang()
		case '-':
			return lx.lexMinus()
		case '<':
			return lx.lexLess()
		case '>':
			return lx.lexGreater()
		case '=':
			return lx.lexEql()
		case '&':
			return lx.lexAmpersand()
		case '|':
			return lx.lexPipe()

		default:
			if unicode.IsSpace(ch) {
				continue
			}
			if IsLetter(ch) {
				return lx.lexIdentifier()
			}
			if IsDigit(ch) {
				return lx.lexNumber()
			}

			return tokUNKNOWN, lx.pos, nil
		}
	}
}

func (lx *Lexer) lexNewline() (Token, Position, error) {
	return tokNEWLN, lx.pos, nil
}

func (lx *Lexer) lexSemicolon() (Token, Position, error) {
	startPos := lx.pos

	looked, err := lx.peekRune()
	if err != nil && !errors.Is(err, io.EOF) {
		return tokUNKNOWN, lx.pos, err
	}

	if err == nil && !unicode.IsSpace(looked) {
		return tokUNKNOWN, lx.pos, ErrInvalidTerminator
	}

	return tokSEMICOLON, startPos, nil
}

func (lx *Lexer) lexComment() (Token, Position, error) {
	pos := lx.pos

	peeked, err := lx.peekRune()
	if err != nil && !errors.Is(err, io.EOF) {
		return tokUNKNOWN, lx.pos, err
	}

	if peeked == '{' {
		lx.scanRune()

		var lastRune rune
		for {
			r, err := lx.scanRune()
			if err != nil && !errors.Is(err, io.EOF) {
				return tokUNKNOWN, lx.pos, err
			}

			if lastRune == '}' && r == '#' {
				break
			}
			lastRune = r
		}
	} else {
		for {
			r, err := lx.scanRune()
			if err != nil && !errors.Is(err, io.EOF) {
				return tokUNKNOWN, lx.pos, err
			}

			if r == Newln || errors.Is(err, io.EOF) {
				break
			}
		}
	}

	return tokCOMMENT, pos, nil
}

func (lx *Lexer) lexSymbol() (Token, Position, error) {
	tok := tokUNKNOWN
	startPos := lx.pos

	var s strings.Builder
	s.Grow(DefaultSymbolAlloc)

	s.WriteRune(SymbolPrefix)
	for {
		r, err := lx.scanRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return tok, lx.pos, err
		}

		s.WriteRune(r)

		lr, err := lx.peekRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return tok, lx.pos, err
		}

		if (!IsLetter(lr) && !IsDigit(lr)) || unicode.IsSpace(lr) {
			break
		}
	}

	lit := s.String()
	if !IsSymbol(lit) {
		return tok, lx.pos, ErrBadSymbol
	}

	tok = Token{SYMBOL, lit}
	return tok, startPos, nil
}

func (lx *Lexer) lexCharacter() (Token, Position, error) {
	startPos := lx.pos
	tok := tokUNKNOWN

	var s strings.Builder
	s.Grow(DefaultCharAlloc)

	r, err := lx.scanRune()
	if err != nil && !errors.Is(err, io.EOF) {
		return tok, lx.pos, err
	}

	if r == CharacterMark {
		return tok, lx.pos, ErrCharEmpty
	}

	s.WriteRune(r)

	if r == '\\' {
		r2, err2 := lx.scanRune()
		if err2 != nil && !errors.Is(err2, io.EOF) {
			return tok, lx.pos, err2
		}
		s.WriteRune(r2)
	}

	closeRune, err := lx.scanRune()
	if err != nil && !errors.Is(err, io.EOF) {
		return tok, lx.pos, err
	}

	if closeRune != CharacterMark {
		return tok, lx.pos, ErrCharBadClose
	}

	tok = Token{CHAR, s.String()}
	return tok, startPos, nil
}

func (lx *Lexer) lexString() (Token, Position, error) {
	tok := tokUNKNOWN
	startPos := lx.pos

	var s strings.Builder
	s.Grow(DefaultStrAlloc)

	for {
		r, err := lx.scanRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return tok, lx.pos, ErrUnterminatedStr
			}
			return tok, lx.pos, err
		}

		if r == '\\' {
			s.WriteRune('\\')

			r2, err := lx.scanRune()
			if err != nil && !errors.Is(err, io.EOF) {
				return tok, lx.pos, err
			}

			s.WriteRune(r2)
			continue
		}

		if r == StringMark {
			break
		}

		s.WriteRune(r)
	}

	lit := s.String()
	if parsed, err := strconv.Unquote(`"` + s.String() + `"`); err == nil {
		tok = Token{STRING, parsed}
	} else {
		tok = Token{STRING, lit}
	}

	return tok, startPos, nil
}

func (lx *Lexer) lexRawString() (Token, Position, error) {
	startPos := lx.pos

	var s strings.Builder
	s.Grow(DefaultRawStrAlloc)

	for {
		r, err := lx.scanRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return Token{Kind: UNKNOWN, Literal: s.String()}, lx.pos, ErrUnterminatedRawStr
			}
			return Token{Kind: UNKNOWN, Literal: s.String()}, lx.pos, err
		}

		if r == '`' {
			break
		}

		s.WriteRune(r)
	}

	return Token{Kind: RAWSTR, Literal: s.String()}, startPos, nil
}

func (lx *Lexer) lexIdentifier() (Token, Position, error) {
	startPos := lx.pos
	tok := tokUNKNOWN

	var s strings.Builder
	s.Grow(DefaultSymbolAlloc)

	s.WriteRune(lx.lastr)
	for {
		r, err := lx.peekRune()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return tok, lx.pos, err
		}

		if !IsLetter(r) && !IsDigit(r) {
			break
		}

		lx.scanRune()
		s.WriteRune(r)
	}

	lit := s.String()
	switch lit {
	case "true", "false":
		return Token{BOOL, lit}, startPos, nil
	}

	return Token{IDENT, lit}, startPos, nil
}

func (lx *Lexer) lexNumber() (Token, Position, error) {
	startPos := lx.pos
	tok := tokUNKNOWN

	var s strings.Builder
	s.WriteRune(lx.lastr)

	handler := IsDigit
	isHex := false
	if lx.lastr == '0' {
		nx, err := lx.peekRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return tok, lx.pos, err
		}

		b := false
		switch nx {
		case 'b':
			handler = IsBinary
			b = true
		case 'o':
			handler = IsOctal
			b = true
		case 'x':
			handler = IsHex
			isHex = true
			b = true
		}

		if b {
			flag, err := lx.scanRune()
			if err != nil && !errors.Is(err, io.EOF) {
				return tok, lx.pos, err
			}
			s.WriteRune(flag)
		}
	}

	ok, err := lx.parseInt(&s, handler)
	if err != nil {
		return tok, lx.pos, err
	}
	if ok {
		tok.Kind = INT
	}

	ok, err = lx.parseFloat(&s, handler)
	if err != nil {
		return tok, lx.pos, err
	}
	if ok {
		tok.Kind = FLOAT
	}

	ok, err = lx.parseScientific(&s, handler, isHex)
	if err != nil {
		return tok, lx.pos, err
	}
	if ok {
		tok.Kind = FLOAT
	}

	ok, err = lx.parseImaginary(&s)
	if err != nil {
		return tok, lx.pos, err
	}
	if ok {
		tok.Kind = IMAG
	}

	tok.Literal = s.String()
	return tok, startPos, nil
}

func (lx *Lexer) parseInt(s *strings.Builder, handler func(rune) bool) (bool, error) {
	for {
		r, err := lx.peekRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return false, err
		}
		if r == '_' {
			_, err := lx.scanRune()
			if err != nil && !errors.Is(err, io.EOF) {
				return false, err
			}
			continue
		}
		if !handler(r) {
			break
		}

		lx.scanRune()
		s.WriteRune(r)
	}
	return true, nil
}

func (lx *Lexer) parseFloat(s *strings.Builder, handler func(rune) bool) (bool, error) {
	nx, err := lx.peekRune()
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}

	isFloat := false
	if nx == '.' {
		s.WriteRune('.')
		lx.scanRune()

		r, errpeek := lx.peekRune()
		if errpeek != nil && !errors.Is(errpeek, io.EOF) {
			return false, errpeek
		}

		if handler(r) {
			b := true
			for {
				r, errscan := lx.peekRune()
				if errscan != nil && !errors.Is(errscan, io.EOF) {
					return false, errscan
				}

				if b && r == '_' {
					return false, ErrInvalidUnderscore
				}
				if r == '_' {
					lx.scanRune()
					continue
				}
				if !handler(r) {
					break
				}

				b = false
				s.WriteRune(r)

				_, err := lx.scanRune()
				if err != nil && !errors.Is(err, io.EOF) {
					return false, err
				}
			}
		}
		isFloat = true
	}

	return isFloat, nil
}

func (lx *Lexer) parseImaginary(s *strings.Builder) (bool, error) {
	isImag := false
	if lx.lastr == ImaginaryPart {
		nx, err := lx.peekRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return false, err
		}

		if unicode.IsSpace(nx) || err == io.EOF {
			s.WriteRune(ImaginaryPart)
			isImag = true
		}
	} else {

		nx, err := lx.peekRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return false, err
		}

		if nx == ImaginaryPart {
			_, err := lx.scanRune()
			if err != nil && !errors.Is(err, io.EOF) {
				return false, err
			}

			nx, err := lx.peekRune()
			if err != nil && !errors.Is(err, io.EOF) {
				return false, err
			}

			if unicode.IsSpace(nx) || err == io.EOF {
				s.WriteRune(ImaginaryPart)
				isImag = true
			}
		}
	}

	return isImag, nil
}

func (lx *Lexer) parseScientific(s *strings.Builder, handler func(rune) bool, isHex bool) (bool, error) {
	hasNotation := false

	nx, err := lx.peekRune()
	if err != nil && !errors.Is(err, io.EOF) {
		return false, err
	}

	isExp := (!isHex && (nx == 'e' || nx == 'E')) || (isHex && (nx == 'p' || nx == 'P'))
	if isExp {
		lx.scanRune()
		s.WriteRune(nx)

		nx, err := lx.peekRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return false, err
		}

		if nx == '+' || nx == '-' {
			sign, err := lx.scanRune()
			if err != nil && !errors.Is(err, io.EOF) {
				return false, err
			}
			s.WriteRune(sign)
		}

		b := true
		for {
			r, err := lx.scanRune()
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return false, err
			}

			if b && r == '_' {
				return false, ErrInvalidUnderscore
			}
			if r == '_' {
				_, err := lx.scanRune()
				if err != nil && !errors.Is(err, io.EOF) {
					return false, err
				}
				continue
			}
			if !handler(r) {
				break
			}

			b = false
			s.WriteRune(r)
		}
		hasNotation = true
	}

	return hasNotation, nil
}

func (lx *Lexer) lexColon() (Token, Position, error) {
	startPos := lx.pos
	return tokCOLON, startPos, nil // :
}

func (lx *Lexer) lexBang() (Token, Position, error) {
	startPos := lx.pos
	if nxt, _ := lx.peekRune(); nxt == '=' {
		_, err := lx.scanRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return tokUNKNOWN, lx.pos, err
		}

		return tokNEQ, startPos, nil // !=
	}
	return tokNOT, startPos, nil
}

func (lx *Lexer) lexMinus() (Token, Position, error) {
	startPos := lx.pos
	return tokSUB, startPos, nil // -
}

func (lx *Lexer) lexLess() (Token, Position, error) {
	startPos := lx.pos

	nxt, err := lx.peekRune()
	if err != nil && !errors.Is(err, io.EOF) {
		return tokUNKNOWN, lx.pos, err
	}

	switch nxt {
	case '=':
		lx.scanRune()
		return tokLEQ, startPos, nil // <=
	}

	return tokLSS, startPos, nil // <
}

func (lx *Lexer) lexGreater() (Token, Position, error) {
	startPos := lx.pos
	if nxt, _ := lx.peekRune(); nxt == '=' {
		_, err := lx.scanRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return tokUNKNOWN, lx.pos, err
		}

		return tokGEQ, startPos, nil // >=
	}
	return tokGTR, startPos, nil // >
}

func (lx *Lexer) lexPeriod() (Token, Position, error) {
	startPos := lx.pos

	nxt, err := lx.peekRune()
	if err != nil && !errors.Is(err, io.EOF) {
		return tokUNKNOWN, lx.pos, err
	}

	if nxt == '.' {
		_, err := lx.scanRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return tokUNKNOWN, lx.pos, err
		}

		if nxt3, _ := lx.peekRune(); nxt3 == '.' {
			lx.scanRune()
			return tokELLIPSIS, startPos, nil
		}
		return tokUNKNOWN, startPos, fmt.Errorf("unexpected '..', expected '...'")
	}
	return tokPERIOD, startPos, nil
}

func (lx *Lexer) lexEql() (Token, Position, error) {
	startPos := lx.pos
	if nxt, _ := lx.peekRune(); nxt == '=' {
		_, err := lx.scanRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return tokUNKNOWN, lx.pos, err
		}

		return tokEQL, startPos, nil // ==
	}
	return tokUNKNOWN, startPos, fmt.Errorf("unexpected character '=', expected '=='")
}

func (lx *Lexer) lexAmpersand() (Token, Position, error) {
	startPos := lx.pos
	if nxt, _ := lx.peekRune(); nxt == '&' {
		_, err := lx.scanRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return tokUNKNOWN, lx.pos, err
		}
		return tokAND, startPos, nil // &&
	}
	return tokUNKNOWN, startPos, fmt.Errorf("unexpected character '&'")
}

func (lx *Lexer) lexPipe() (Token, Position, error) {
	startPos := lx.pos
	if nxt, _ := lx.peekRune(); nxt == '|' {
		_, err := lx.scanRune()
		if err != nil && !errors.Is(err, io.EOF) {
			return tokUNKNOWN, lx.pos, err
		}
		return tokOR, startPos, nil // ||
	}
	return tokUNKNOWN, startPos, fmt.Errorf("unexpected character '|'")
}
