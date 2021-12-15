package main

import (
	"fmt"
	"os"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/parser"
	"recital_lsp/pkg/utils"
	"time"
)

func readFileContent(fileName string) string {

	fileReadStart := time.Now()

	dat, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	dat = utils.NormalizeNewlines(dat)

	fmt.Printf("\n-> File read time: %d ms\n", (time.Since(fileReadStart)).Milliseconds())

	return string(dat)
}

func execLexer(fileName string, fileContent string) *lexer.LexerResult {

	lexerStartTime := time.Now()

	lex := lexer.NewLexer(fileName, fileContent)

	res, err := lex.Parse()

	lexerEndTime := time.Now()

	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}

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

	return res
}

func execParser(lexerRes *lexer.LexerResult, fileName string) *parser.ParseResult {

	startTime := time.Now()

	parser := parser.NewParser(lexerRes)

	res := parser.Parse()

	fmt.Printf("-> Parser time: %d ms\n\n", (time.Since(startTime)).Milliseconds())

	if res.Err != nil {
		fmt.Printf("%s\n", res.Err.StringWithProgram(fileName))
	} else {
		fmt.Printf("ok\n")
	}

	return res

	// if len(res.Errors) == 0 {
	// 	fmt.Printf("-> No errors\n")
	// } else {
	// 	fmt.Printf("-> Found %d errors\n", len(res.Errors))
	// }

	// if len(res.Warnings) == 0 {
	// 	fmt.Printf("-> No warnings\n")
	// } else {
	// 	fmt.Printf("-> Found %d warnings\n", len(res.Warnings))
	// }

	// fmt.Printf("\n-> Lexer result: %d tokens\n", res.TokensCount)

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: rt <command>")
		return
	}

	command := os.Args[1]

	if command == "lex" {

		if len(os.Args) < 4 {
			fmt.Println("Usage: rt lex <filename> <html output filename>")
			return
		}

		fileName := os.Args[2]
		htmlOutputFileName := os.Args[3]

		fileContent := readFileContent(fileName)

		res := execLexer(fileName, fileContent)

		lexer.PrintLexerResultToHTML(res, htmlOutputFileName)

	} else if command == "parse" {

		if len(os.Args) < 4 {
			fmt.Println("Usage: rt parse <filename> <html output filename>")
			return
		}

		fileName := os.Args[2]
		htmlOutputFileName := os.Args[3]

		fileContent := readFileContent(fileName)

		res := execLexer(fileName, fileContent)

		if len(res.Errors) > 0 {
			return
		}

		parseRes := execParser(res, fileName)

		if parseRes.Err == nil {
			parser.PrintParseResultToHTML(parseRes, htmlOutputFileName)
		}

	}

}
