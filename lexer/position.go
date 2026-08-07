package lexer

import "fmt"

type Position struct {
	line   int
	column int
}

func (p Position) Line() int {
	return p.line
}

func (p Position) Column() int {
	return p.column
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.line, p.column)
}
