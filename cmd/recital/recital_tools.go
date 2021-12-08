package main

import (
	"fmt"
	"os"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/lsp"
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

		if len(os.Args) < 4 {
			fmt.Println("Usage: recital lex <filename> <html output filename>")
			return
		}

		fileName := os.Args[2]
		htmlOutputFileName := os.Args[3]

		fileReadStart := time.Now()

		dat, err := os.ReadFile(fileName)

		fileReadEnd := time.Now()

		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		dat = utils.NormalizeNewlines(dat)

		lexerStartTime := time.Now()

		lex := lexer.NewLexer(fileName, string(dat))

		res, err := lex.Parse()

		lexerEndTime := time.Now()

		if err != nil {
			fmt.Printf("err: %v\n", err)
			return
		}

		fmt.Printf("\n-> File read time: %d ms\n", (fileReadEnd.Sub(fileReadStart)).Milliseconds())
		fmt.Printf("-> Lexer time: %d ms\n\n", (lexerEndTime.Sub(lexerStartTime)).Milliseconds())

		if len(res.Errors) == 0 {
			fmt.Printf("-> No errors\n")
		} else {
			fmt.Printf("-> Found %d errors\n", len(res.Errors))
		}

		if len(res.Warnings) == 0 {
			fmt.Printf("-> No warnings\n")
		} else {
			fmt.Printf("-> Found %d warnings\n", len(res.Warnings))
		}

		fmt.Printf("\n-> Lexer result: %d tokens\n", res.TokensCount)

		lexer.PrintLexerResultToHTML(res, htmlOutputFileName)

	} else if command == "lsp" {

		lsp.CreateServer()

	}

}
