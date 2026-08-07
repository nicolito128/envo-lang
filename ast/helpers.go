package ast

import "strings"

func joinNodes(nodes []Node, sep string) string {
	var b strings.Builder
	for i, node := range nodes {
		if i > 0 {
			b.WriteString(sep)
		}
		if node == nil {
			continue
		}
		b.WriteString(node.String())
	}
	return b.String()
}

func Equal(a, b Node) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.Type() != b.Type() {
		return false
	}

	switch a := a.(type) {
	case *Program:
		b := b.(*Program)
		if len(a.Statements) != len(b.Statements) {
			return false
		}
		for i := range a.Statements {
			if !Equal(a.Statements[i], b.Statements[i]) {
				return false
			}
		}
		return true

	case *Identifier:
		b := b.(*Identifier)
		return a.Name == b.Name

	case *Message:
		b := b.(*Message)
		return Equal(a.Label, b.Label) && Equal(a.Body, b.Body)

	case *RawMessage:
		b := b.(*RawMessage)
		if len(a.Words) != len(b.Words) {
			return false
		}

		for i := range a.Words {
			if !Equal(a.Words[i], b.Words[i]) {
				return false
			}
		}
		return true

	case *SymbolLiteral:
		b := b.(*SymbolLiteral)
		return a.Name == b.Name

	case *IntegerLiteral:
		b := b.(*IntegerLiteral)
		return a.Value == b.Value

	case *FloatLiteral:
		b := b.(*FloatLiteral)
		return a.Value == b.Value

	case *ImaginaryLiteral:
		b := b.(*ImaginaryLiteral)
		return a.Value == b.Value

	case *Literal:
		b := b.(*Literal)
		return a.Value == b.Value

	case *BinaryOp:
		b := b.(*BinaryOp)
		return a.Op == b.Op && Equal(a.Left, b.Left) && Equal(a.Right, b.Right)

	case *UnaryOp:
		b := b.(*UnaryOp)
		return a.Op == b.Op && Equal(a.Operand, b.Operand)

	case *Define:
		b := b.(*Define)
		return Equal(a.Receiver, b.Receiver) && Equal(a.Pattern, b.Pattern) && Equal(a.Response, b.Response)

	default:
		return false
	}
}
