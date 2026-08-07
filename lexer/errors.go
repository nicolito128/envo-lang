package lexer

import "errors"

var (
	ErrBadSymbol          = errors.New("the symbol must be a valid identifier")
	ErrInvalidTerminator  = errors.New("invalid usage for the expression terminator")
	ErrCharBadClose       = errors.New("bad close symbol for character")
	ErrCharEmpty          = errors.New("character must not be empty")
	ErrCharTooLong        = errors.New("character must be one rune long")
	ErrInvalidUnderscore  = errors.New("underscore place is not valid")
	ErrUnterminatedStr    = errors.New("unterminated string")
	ErrUnterminatedRawStr = errors.New("unterminated raw string")
)
