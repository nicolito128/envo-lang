package lexer

import (
	"unicode"
)

func IsLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func IsDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func IsBinary(r rune) bool {
	return r == '0' || r == '1'
}

func IsOctal(r rune) bool {
	return r >= '0' && r <= '7'
}

func IsHex(r rune) bool {
	return unicode.IsDigit(r) || (r >= 'A' && r <= 'F') || (r >= 'a' && r <= 'f')
}

func IsOperator(s string) bool {
	_, ok := operators[s]
	return ok
}

func IsIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	if !IsLetter(rune(s[0])) {
		return false
	}
	for _, r := range s {
		if !IsLetter(r) && !IsDigit(r) {
			return false
		}
	}
	return true
}

func IsSymbol(s string) bool {
	if len(s) <= 0 {
		return false
	}
	if rune(s[0]) != SymbolPrefix {
		return false
	}
	for _, r := range s[1:] {
		if !IsLetter(r) && !IsDigit(r) {
			return false
		}
	}
	return true
}
