package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nicolito128/envo-lang/lexer"
	"github.com/nicolito128/envo-lang/object"
	"github.com/nicolito128/envo-lang/parser"
	"github.com/nicolito128/envo-lang/repl"
	"github.com/nicolito128/envo-lang/runtime"
)

func main() {
	if len(os.Args) < 2 {
		repl.Start(os.Stdin, os.Stdout)
		return
	}

	filePath := os.Args[1]
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("error trying to open the file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	l := lexer.New(strings.NewReader(string(content)))
	p := parser.New(l)

	program, parseErrors := p.Parse()
	if len(parseErrors) > 0 {
		fmt.Printf("parser error: %s:\n", filePath)
		for _, parseErr := range parseErrors {
			fmt.Printf("  %s\n", parseErr)
		}
		os.Exit(1)
	}

	env := runtime.NewEnvironment()
	result := runtime.Eval(program, env)

	if result != nil && result.Type() != object.NIL_OBJ {
		if result.Type() == object.ERROR_OBJ {
			fmt.Printf("runtime error: %s\n", result.String())
			os.Exit(1)
		}
		fmt.Println(result.String())
	}
}
