package main

import (
	"flag"
	"fmt"
	"os"

	"n128.xyz/n128/envo/lexer"
)

var (
	fileFlag = flag.String("f", "./example.env", "Filepath to the Envo program")
)

func main() {
	flag.Parse()

	file, err := os.Open(*fileFlag)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lx := lexer.NewLexer(file)

	header := "Pos\t\tKind\t\tLiteral\n"
	fmt.Println(header)

	for {
		tok, pos, err := lx.Lex()
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\t\t%s\t\t%s\n", pos, tok.Kind, tok.Literal)
		if tok.Kind == lexer.EOF {
			break
		}
	}
}
