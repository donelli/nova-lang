package main

import (
	"fmt"
	"os"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/utils"
	"time"
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

		fileReadStart := time.Now()

		dat, err := os.ReadFile(os.Args[2])

		fileReadEnd := time.Now()

		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		dat = utils.NormalizeNewlines(dat)

		lexerStartTime := time.Now()

		lex := lexer.NewLexer(os.Args[2], string(dat))

		res, err := lex.Parse()

		lexerEndTime := time.Now()

		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		fmt.Printf("\n-> File read time: %d ms\n", (fileReadEnd.Sub(fileReadStart)).Milliseconds())
		fmt.Printf("-> Lexer time: %d ms\n", (lexerEndTime.Sub(lexerStartTime)).Milliseconds())

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
