package runtime_test

import (
	"strings"
	"testing"

	"github.com/nicolito128/envo-lang/lexer"
	"github.com/nicolito128/envo-lang/object"
	"github.com/nicolito128/envo-lang/parser"
	"github.com/nicolito128/envo-lang/runtime"
)

func TestFactorialEvaluation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected object.Object
	}{
		{
			name: "Base case: fact{0}",
			input: `
				fact{0}{1}
				fact{0}
			`,
			expected: &object.RawMessage{Words: []object.Object{&object.Integer{Value: 1}}},
		},
		{
			name: "Recursive case: fact{5}",
			input: `
				fact{0}{1}
				fact{n}{n * fact{n - 1}}
				fact{5}
			`,
			expected: &object.RawMessage{Words: []object.Object{&object.Integer{Value: 120}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(strings.NewReader(tt.input))
			p := parser.New(l)
			program, parseErrors := p.Parse()

			if len(parseErrors) > 0 {
				t.Fatalf("unexpected parser error: %v", parseErrors[0])
			}

			env := runtime.NewEnvironment()
			evaluated := runtime.Eval(program, env)

			testIntegerObject(t, evaluated, tt.expected)
		})
	}
}

func TestComparisons(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "less than", input: "1 < 2", expected: true},
		{name: "less than or equal", input: "2 <= 2", expected: true},
		{name: "greater than", input: "3 > 2", expected: true},
		{name: "greater than or equal", input: "2 >= 2", expected: true},
		{name: "not equal", input: "1 != 2", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(strings.NewReader(tt.input))
			p := parser.New(l)
			program, parseErrors := p.Parse()

			if len(parseErrors) > 0 {
				t.Fatalf("unexpected parser error: %v", parseErrors[0])
			}

			env := runtime.NewEnvironment()
			evaluated := runtime.Eval(program, env)

			testBoolObject(t, evaluated, tt.expected)
		})
	}
}

func testIntegerObject(t *testing.T, obj object.Object, expected object.Object) bool {
	if obj == nil {
		t.Fatalf("result returned by Eval was nil")
		return false
	}

	result, ok := obj.(*object.Integer)
	if !ok {
		if errObj, isErr := obj.(*object.Error); isErr {
			t.Fatalf("unexpected runtime error: %s", errObj.Message)
		}
		t.Fatalf("object type is not *object.Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if !object.Equal(result, expected) {
		t.Errorf("object has wrong integer value. got=%d, want=%v", result.Value, expected)
		return false
	}

	return true
}

func testBoolObject(t *testing.T, obj object.Object, expected bool) bool {
	if obj == nil {
		t.Fatalf("result returned by Eval was nil")
		return false
	}

	result, ok := obj.(*object.Bool)
	if !ok {
		if errObj, isErr := obj.(*object.Error); isErr {
			t.Fatalf("unexpected runtime error: %s", errObj.Message)
		}
		t.Fatalf("object type is not *object.Bool. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong boolean value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}
