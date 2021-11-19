package main

import (
	"fmt"
	"os"
	"recital_lsp/pkg/lexer"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: recital <command>")
		return
	}

	command := os.Args[1]

	if command == "lex" {

		if len(os.Args) < 3 {
			fmt.Println("Usage: recital lex <filename>")
			return
		}

		dat, err := os.ReadFile(os.Args[2])

		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		lex := lexer.NewLexer(os.Args[2], string(dat))

		res, err := lex.Parse()

		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		fmt.Printf("%+v\n", res.String())

	}

}
