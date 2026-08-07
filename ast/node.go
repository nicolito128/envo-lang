package ast

type NodeType int

const (
	ProgramNode NodeType = iota
	IdentifierNode
	LiteralNode
	IntegerLiteralNode
	FloatLiteralNode
	ImaginaryLiteralNode
	BoolLiteralNode
	CharLiteralNode
	StringLiteralNode
	RawStringLiteralNode
	SymbolLiteralNode
	RawMessageNode
	MessageNode
	DefineNode
	BinaryOpNode
	UnaryOpNode
	ExpressionNode
)

type Node interface {
	Node()
	Type() NodeType
	String() string
}

type nodeImpl struct{}

func (n nodeImpl) Node()          {}
func (n nodeImpl) Type() NodeType { return ProgramNode }
func (n nodeImpl) String() string { return "" }
