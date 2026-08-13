package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/nicolito128/envo-lang/ast"
	"github.com/nicolito128/envo-lang/lexer"
)

func TestParser(t *testing.T) {
	data := [...]string{
		"{}",
		"{1}+{2}",
		"foo{}",
		"foo{~bar}",
		"foo{bar{},baz{~ok}}",
		"id{x}{x}",
		"ab{a}{b}",
		"fact{0}{1}",
		"fact{n}{n*fact{n-1}}",
		"((foo{(~a)}))",
	}

	expected := [...]ast.Node{
		&ast.RawMessage{Words: []ast.Node{}},
		&ast.BinaryOp{
			Op: "+",
			Left: &ast.RawMessage{
				Words: []ast.Node{&ast.IntegerLiteral{Value: 1}},
			},
			Right: &ast.RawMessage{
				Words: []ast.Node{&ast.IntegerLiteral{Value: 2}},
			},
		},
		&ast.Message{
			Label: &ast.Identifier{Name: "foo"},
			Body:  &ast.RawMessage{Words: []ast.Node{}},
		},
		&ast.Message{
			Label: &ast.Identifier{Name: "foo"},
			Body:  &ast.RawMessage{Words: []ast.Node{&ast.SymbolLiteral{Name: "~bar"}}},
		},
		&ast.Message{
			Label: &ast.Identifier{Name: "foo"},
			Body: &ast.RawMessage{
				Words: []ast.Node{
					&ast.Message{
						Label: &ast.Identifier{Name: "bar"},
						Body:  &ast.RawMessage{Words: []ast.Node{}},
					},
					&ast.Message{
						Label: &ast.Identifier{Name: "baz"},
						Body:  &ast.RawMessage{Words: []ast.Node{&ast.SymbolLiteral{Name: "~ok"}}},
					},
				},
			},
		},
		&ast.Define{
			Receiver: &ast.Identifier{Name: "id"},
			Pattern: &ast.RawMessage{
				Words: []ast.Node{
					&ast.Identifier{Name: "x"},
				},
			},
			Response: &ast.RawMessage{
				Words: []ast.Node{
					&ast.Identifier{Name: "x"},
				},
			},
		},
		&ast.Define{
			Receiver: &ast.Identifier{Name: "ab"},
			Pattern: &ast.RawMessage{
				Words: []ast.Node{
					&ast.Identifier{Name: "a"},
				},
			},
			Response: &ast.RawMessage{
				Words: []ast.Node{
					&ast.Identifier{Name: "b"},
				},
			},
		},
		&ast.Define{
			Receiver: &ast.Identifier{Name: "fact"},
			Pattern: &ast.RawMessage{
				Words: []ast.Node{
					&ast.IntegerLiteral{Value: 0},
				},
			},
			Response: &ast.RawMessage{
				Words: []ast.Node{
					&ast.IntegerLiteral{Value: 1},
				},
			},
		},
		&ast.Define{
			Receiver: &ast.Identifier{Name: "fact"},
			Pattern: &ast.RawMessage{
				Words: []ast.Node{
					&ast.Identifier{Name: "n"},
				},
			},
			Response: &ast.RawMessage{
				Words: []ast.Node{
					&ast.BinaryOp{
						Op:   "*",
						Left: &ast.Identifier{Name: "n"},
						Right: &ast.Message{
							Label: &ast.Identifier{Name: "fact"},
							Body: &ast.RawMessage{
								Words: []ast.Node{
									&ast.BinaryOp{
										Op:    "-",
										Left:  &ast.Identifier{Name: "n"},
										Right: &ast.IntegerLiteral{Value: 1},
									},
								},
							},
						},
					},
				},
			},
		},
		&ast.Message{
			Label: &ast.Identifier{Name: "foo"},
			Body: &ast.RawMessage{
				Words: []ast.Node{
					&ast.SymbolLiteral{Name: "~a"},
				},
			},
		},
	}

	for i, d := range data {
		l := lexer.New(strings.NewReader(d))
		p := New(l)

		node, errors := p.Parse()

		for _, err := range errors {
			t.Errorf("unexpected error: %s", err)
		}

		program, ok := node.(*ast.Program)
		if !ok {
			t.Errorf("expected a program as output of the parser")
		}

		testname := fmt.Sprintf("TestParser[%s]", d)
		t.Run(testname, func(t *testing.T) {
			if len(program.Statements) != 1 {
				t.Fatalf("expected 1 statement, got %d", len(program.Statements))
			}

			for _, stmt := range program.Statements {
				if !ast.Equal(stmt, expected[i]) {
					// print human-friendly strings to help debugging mismatches
					t.Logf("expected str: %s", expected[i].String())
					t.Logf("got str:      %s", stmt.String())
					t.Errorf("expected statement %d to be %#v, got %#v", i, expected[i], stmt)
				}
			}
		})
	}
}
