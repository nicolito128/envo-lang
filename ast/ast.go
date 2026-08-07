package ast

import (
	"fmt"
	"strconv"
	"strings"

	"n128.xyz/n128/envo/lexer"
)

type Program struct {
	nodeImpl

	Statements []Node
}

func (p *Program) String() string {
	if p == nil {
		return ""
	}
	var b strings.Builder
	for i, stmt := range p.Statements {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(stmt.String())
	}
	return b.String()
}

type Identifier struct {
	nodeImpl

	Name string
}

func (id *Identifier) String() string {
	return id.Name
}

func (id *Identifier) Type() NodeType {
	return IdentifierNode
}

type Literal struct {
	nodeImpl

	Kind  lexer.TokenKind
	Value any
}

func (lit *Literal) Type() NodeType {
	return LiteralNode
}

func (lit *Literal) String() string {
	if lit == nil {
		return ""
	}

	if lit.Value == nil {
		return ""
	}

	switch lit.Kind {
	case lexer.INT, lexer.FLOAT, lexer.IMAG:
		return fmt.Sprint(lit.Value)
	case lexer.STRING:
		if s, ok := lit.Value.(string); ok {
			return strconv.Quote(s)
		}
	case lexer.RAWSTR:
		if s, ok := lit.Value.(string); ok {
			if strings.ContainsRune(s, '`') {
				return strconv.Quote(s)
			}
			return "`" + s + "`"
		}
	case lexer.CHAR:
		if s, ok := lit.Value.(string); ok {
			q := strconv.Quote(s)
			return "'" + q[1:len(q)-1] + "'"
		}
	case lexer.SYMBOL:
		if s, ok := lit.Value.(string); ok {
			return s
		}
	}

	return fmt.Sprint(lit.Value)
}

type IntegerLiteral struct {
	nodeImpl

	Value int64
}

func (il *IntegerLiteral) Type() NodeType {
	return IntegerLiteralNode
}

func (il *IntegerLiteral) String() string {
	if il == nil {
		return ""
	}
	return fmt.Sprint(il.Value)
}

type FloatLiteral struct {
	nodeImpl

	Value float64
}

func (fl *FloatLiteral) Type() NodeType {
	return FloatLiteralNode
}

func (fl *FloatLiteral) String() string {
	if fl == nil {
		return ""
	}
	return fmt.Sprint(fl.Value)
}

type BoolLiteral struct {
	nodeImpl

	Value bool
}

func (b *BoolLiteral) Type() NodeType {
	return BoolLiteralNode
}

func (b *BoolLiteral) String() string {
	return fmt.Sprintf("%v", b.Value)
}

type CharLiteral struct {
	nodeImpl

	Value []rune
}

func (charLit *CharLiteral) Type() NodeType {
	return CharLiteralNode
}

func (charLit *CharLiteral) String() string {
	if charLit.Value == nil {
		return ""
	}
	q := string(charLit.Value)
	return "'" + q + "'"
}

type StringLiteral struct {
	nodeImpl

	Value string
}

func (strLit *StringLiteral) Type() NodeType {
	return StringLiteralNode
}

func (strLit *StringLiteral) String() string {
	return strconv.Quote(strLit.Value)
}

type RawStringLiteral struct {
	nodeImpl

	Value string
}

func (rawStrLit *RawStringLiteral) Type() NodeType {
	return RawStringLiteralNode
}

func (rawStrLit *RawStringLiteral) String() string {
	if rawStrLit == nil {
		return ""
	}
	if strings.ContainsRune(rawStrLit.Value, '`') {
		return strconv.Quote(rawStrLit.Value)
	}
	return "`" + rawStrLit.Value + "`"
}

type ImaginaryLiteral struct {
	nodeImpl

	Value complex128
}

func (imagLit *ImaginaryLiteral) Type() NodeType {
	return ImaginaryLiteralNode
}

func (imagLit *ImaginaryLiteral) String() string {
	if imagLit == nil {
		return ""
	}
	return fmt.Sprintf("%g%+gi", real(imagLit.Value), imag(imagLit.Value))
}

type SymbolLiteral struct {
	nodeImpl
	Name string
}

func (sym *SymbolLiteral) Type() NodeType {
	return SymbolLiteralNode
}

func (sym *SymbolLiteral) String() string {
	if sym == nil {
		return ""
	}
	return sym.Name
}

type RawMessage struct {
	nodeImpl

	Words []Node
}

func (rm *RawMessage) Type() NodeType {
	return RawMessageNode
}

func (rm *RawMessage) String() string {
	if rm == nil {
		return ""
	}
	var s strings.Builder
	s.WriteString("{")
	for i, node := range rm.Words {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(node.String())
	}
	s.WriteString("}")
	return s.String()
}

type Message struct {
	nodeImpl

	Label *Identifier
	Body  *RawMessage
}

func (msg *Message) Type() NodeType {
	return MessageNode
}

func (msg *Message) String() string {
	if msg == nil {
		return ""
	}
	if msg.Label == nil {
		return msg.Body.String()
	}
	return msg.Label.String() + msg.Body.String()
}

type Define struct {
	nodeImpl

	Receiver *Identifier
	Pattern  *RawMessage
	Response *RawMessage
}

func (def *Define) Type() NodeType {
	return DefineNode
}

func (def *Define) String() string {
	if def == nil {
		return ""
	}
	return def.Receiver.String() + def.Pattern.String() + def.Response.String()
}

type BinaryOp struct {
	nodeImpl

	Op    string
	Left  Node
	Right Node
}

func (bin *BinaryOp) Type() NodeType {
	return BinaryOpNode
}

func (bin *BinaryOp) String() string {
	if bin == nil {
		return ""
	}
	return bin.Left.String() + bin.Op + bin.Right.String()
}

type UnaryOp struct {
	nodeImpl

	Op      string
	Operand Node
}

func (un *UnaryOp) Type() NodeType {
	return UnaryOpNode
}

func (un *UnaryOp) String() string {
	if un == nil {
		return ""
	}
	return un.Op + un.Operand.String()
}

type Expression struct {
	nodeImpl

	Binary  *BinaryOp
	Unary   *UnaryOp
	Operand Node
}

func (expr *Expression) Type() NodeType {
	return ExpressionNode
}

func (expr *Expression) String() string {
	if expr == nil {
		return ""
	}
	if expr.Binary != nil {
		return expr.Binary.String()
	}
	if expr.Unary != nil {
		return expr.Unary.String()
	}
	if expr.Operand != nil {
		return expr.Operand.String()
	}
	return ""
}
