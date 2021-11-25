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

		if len(res.Errors) == 0 {
			fmt.Printf("\n-> No errors\n")
		} else {
			fmt.Printf("\n-> Found %d errors\n", len(res.Errors))

			for error := range res.Errors {
				fmt.Println(res.Errors[error])
			}

		}

		if len(res.Warnings) == 0 {
			fmt.Printf("\n-> No warnings\n")
		} else {
			fmt.Printf("\n-> Found %d warnings\n", len(res.Warnings))
			// TODO show warnings
		}

		fmt.Printf("\n-> Lexer result:")
		fmt.Printf("\n%+v\n", res.String())

	}

}
