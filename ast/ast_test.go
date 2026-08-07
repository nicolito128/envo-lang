package ast

import "testing"

func TestAstNodeString(t *testing.T) {
	cases := []struct {
		name string
		node Node
		want string
	}{
		{"identifier", &Identifier{Name: "foo"}, "foo"},
		{"symbol", &SymbolLiteral{Name: "~ok"}, "~ok"},
		{"integer", &IntegerLiteral{Value: 42}, "42"},
		{"float", &FloatLiteral{Value: 3.14}, "3.14"},
		{"imaginary", &ImaginaryLiteral{Value: complex(1, 2)}, "1+2i"},
		{"char", &CharLiteral{Value: []rune("a")}, "'a'"},
		{"string", &StringLiteral{Value: "abc"}, "\"abc\""},
		{"raw string", &RawStringLiteral{Value: "abc"}, "`abc`"},
		{"raw message", &RawMessage{Words: []Node{&Identifier{Name: "foo"}, &IntegerLiteral{Value: 1}}}, "{foo,1}"},
		{"message", &Message{Label: &Identifier{Name: "bar"}, Body: &RawMessage{Words: []Node{&Identifier{Name: "baz"}}}}, "bar{baz}"},
		{"define", &Define{Receiver: &Identifier{Name: "fact"}, Pattern: &RawMessage{Words: []Node{&Identifier{Name: "n"}}}, Response: &RawMessage{Words: []Node{&BinaryOp{Left: &Identifier{Name: "n"}, Op: "*", Right: &Identifier{Name: "fact"}}}}}, "fact{n}{n*fact}"},
		{"binary op", &BinaryOp{Left: &IntegerLiteral{Value: 1}, Op: "+", Right: &IntegerLiteral{Value: 2}}, "1+2"},
		{"unary op", &UnaryOp{Op: "-", Operand: &IntegerLiteral{Value: 5}}, "-5"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.node.String()
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
