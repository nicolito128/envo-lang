package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func TestComments(t *testing.T) {
	data := [...]string{
		"# Inline comment",
		`#{
			Multi-line comment
		}#`,
	}

	expected := [...]Token{
		tokCOMMENT,
		tokCOMMENT,
	}

	for i, d := range data {
		lx := NewLexer(strings.NewReader(d))

		got, _, err := lx.Lex()
		want := expected[i]

		if err != nil {
			t.Errorf("got %v trying to parse the input data", err)
		}

		testname := fmt.Sprintf("TestComments[%d]", i)
		t.Run(testname, func(t *testing.T) {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestSymbols(t *testing.T) {
	data := [...]string{
		"~ok",
		"~_error",
		"~bad_elem",
		"~s123",
	}

	expected := [...]Token{
		{SYMBOL, "~ok"},
		{SYMBOL, "~_error"},
		{SYMBOL, "~bad_elem"},
		{SYMBOL, "~s123"},
	}

	for i, d := range data {
		lx := NewLexer(strings.NewReader(d))

		got, _, err := lx.Lex()
		want := expected[i]

		if err != nil {
			t.Errorf("got %v trying to parse the input data", err)
		}

		testname := fmt.Sprintf("TestSymbols[%d]", i)
		t.Run(testname, func(t *testing.T) {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestCharacters(t *testing.T) {
	data := [...]string{
		"' '",
		"'a'",
		"'\n'",
		"'鞆'",
	}

	expected := [...]Token{
		{CHAR, " "},
		{CHAR, "a"},
		{CHAR, "\n"},
		{CHAR, "鞆"},
	}

	for i, d := range data {
		lx := NewLexer(strings.NewReader(d))

		got, _, err := lx.Lex()
		want := expected[i]

		if err != nil {
			t.Errorf("got %v trying to parse the input data", err)
		}

		testname := fmt.Sprintf("TestCharacters[%d]", i)
		t.Run(testname, func(t *testing.T) {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestStrings(t *testing.T) {
	data := [...]string{
		`""`,
		`"abcd"`,
		`"\"quoted\""`,
		`"a1 b2 c3"`,
		"` A B C \"D\" `",
	}

	expected := [...]Token{
		{STRING, ``},
		{STRING, `abcd`},
		{STRING, `"quoted"`},
		{STRING, `a1 b2 c3`},
		{RAWSTR, ` A B C "D" `},
	}

	for i, d := range data {
		lx := NewLexer(strings.NewReader(d))

		got, _, err := lx.Lex()
		want := expected[i]

		if err != nil {
			t.Errorf("got %v trying to parse the input data", err)
		}

		testname := fmt.Sprintf("TestString[%d]", i)
		t.Run(testname, func(t *testing.T) {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestNumbers(t *testing.T) {
	t.Run("TestInteger", func(t *testing.T) {
		data := [...]string{
			"12345",
			"1_000_000",
			"0",
			"0b1011",
			"0o755",
			"0x1A3F",
		}

		expected := [...]Token{
			{INT, "12345"},
			{INT, "1000000"},
			{INT, "0"},
			{INT, "0b1011"},
			{INT, "0o755"},
			{INT, "0x1A3F"},
			{INT, "0x1A3F"},
		}

		for i, d := range data {
			lx := NewLexer(strings.NewReader(d))

			got, _, err := lx.Lex()
			want := expected[i]

			if err != nil {
				t.Errorf("got %v trying to parse the input data", err)
			}

			testname := fmt.Sprintf("TestInteger[%d]", i)
			t.Run(testname, func(t *testing.T) {
				if got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			})
		}
	})

	t.Run("TestFloat", func(t *testing.T) {
		data := [...]string{
			"1.23456789",
			"0.0001",
			"1.23_45_67",
			"3.1_415",
		}

		expected := [...]Token{
			{FLOAT, "1.23456789"},
			{FLOAT, "0.0001"},
			{FLOAT, "1.234567"},
			{FLOAT, "3.1415"},
		}

		for i, d := range data {
			lx := NewLexer(strings.NewReader(d))

			got, _, err := lx.Lex()
			want := expected[i]

			if err != nil {
				t.Errorf("got %v trying to parse the input data", err)
			}

			testname := fmt.Sprintf("TestFloat[%d]", i)
			t.Run(testname, func(t *testing.T) {
				if got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			})
		}
	})

	t.Run("TestComplex", func(t *testing.T) {
		data := [...]string{
			"1i",
			"12_34i",
			"1.23i",
			"0.000_123i",
		}

		expected := [...]Token{
			{IMAG, "1i"},
			{IMAG, "1234i"},
			{IMAG, "1.23i"},
			{IMAG, "0.000123i"},
		}

		for i, d := range data {
			lx := NewLexer(strings.NewReader(d))

			got, _, err := lx.Lex()
			want := expected[i]

			if err != nil {
				t.Errorf("got %v trying to parse the input data", err)
			}

			testname := fmt.Sprintf("TestComplex[%d]", i)
			t.Run(testname, func(t *testing.T) {
				if got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			})
		}
	})

	t.Run("TestScientificNotation", func(t *testing.T) {
		data := [...]string{
			"1.23e4",
			"654e-2",
			"3.1415e+2",
			"0x1A3Fp2",
			"1.445e-3i",
		}

		expected := [...]Token{
			{FLOAT, "1.23e4"},
			{FLOAT, "654e-2"},
			{FLOAT, "3.1415e+2"},
			{FLOAT, "0x1A3Fp2"},
			{IMAG, "1.445e-3i"},
		}

		for i, d := range data {
			lx := NewLexer(strings.NewReader(d))

			got, _, err := lx.Lex()
			want := expected[i]

			if err != nil {
				t.Errorf("got %v trying to parse the input data", err)
			}

			testname := fmt.Sprintf("TestNotation[%d]", i)
			t.Run(testname, func(t *testing.T) {
				if got != want {
					t.Errorf("got %v, want %v", got, want)
				}
			})
		}
	})
}

func TestIdentifier(t *testing.T) {
	data := [...]string{
		"foo",
		"_Bar",
		"b_az1234",
	}

	expected := [...]Token{
		{IDENT, "foo"},
		{IDENT, "_Bar"},
		{IDENT, "b_az1234"},
	}

	for i, d := range data {
		lx := NewLexer(strings.NewReader(d))

		got, _, err := lx.Lex()
		want := expected[i]

		if err != nil {
			t.Errorf("got %v trying to parse the input data", err)
		}

		testname := fmt.Sprintf("TestIdentifier[%d]", i)
		t.Run(testname, func(t *testing.T) {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestOperators(t *testing.T) {
	data := [...]string{
		"+",
		"-",
		"*",
		"/",
		"%",
		"==",
		"!=",
		"<",
		"<=",
		">",
		">=",
		"...",
		"(",
		"[",
		"{",
		")",
		"]",
		"}",
		",",
		".",
		":",
		";",
		"&&",
		"||",
		"!",
	}

	expected := [...]Token{
		{ADD, "+"},
		{SUB, "-"},
		{MUL, "*"},
		{DIV, "/"},
		{REM, "%"},
		{EQL, "=="},
		{NEQ, "!="},
		{LSS, "<"},
		{LEQ, "<="},
		{GTR, ">"},
		{GEQ, ">="},
		{ELLIPSIS, "..."},
		{LPAREN, "("},
		{LBRACK, "["},
		{LBRACE, "{"},
		{RPAREN, ")"},
		{RBRACK, "]"},
		{RBRACE, "}"},
		{COMMA, ","},
		{PERIOD, "."},
		{COLON, ":"},
		{SEMICOLON, ";"},
		{AND, "&&"},
		{OR, "||"},
		{NOT, "!"},
	}

	for i, d := range data {
		lx := NewLexer(strings.NewReader(d))

		got, _, err := lx.Lex()
		want := expected[i]

		if err != nil {
			t.Errorf("got %v trying to parse the input data", err)
		}

		testname := fmt.Sprintf("TestOperators[%s]", d)
		t.Run(testname, func(t *testing.T) {
			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
