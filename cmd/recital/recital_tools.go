package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/parser"
	"recital_lsp/pkg/shared"
	"recital_lsp/pkg/utils"
)

func readFileContent(fileName string) string {

	dat, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	dat = utils.NormalizeNewlines(dat)

	return string(dat)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: rt <command>")
		return
	}

	command := os.Args[1]

	if command == "parse" {

		if len(os.Args) < 3 {
			fmt.Println("Usage: rt parse <file>")
			return
		}

		fileName := os.Args[2]

		parseCmd := flag.NewFlagSet("rt.exe parse", flag.ExitOnError)
		// outHtml := parseCmd.String("html", "", "File to print the result as HTML")
		outJson := parseCmd.String("json", "", "File to print the result as JSON")

		parseCmd.Parse(os.Args[3:])

		fileContent := readFileContent(fileName)

		errors := []*shared.Error{}
		warnings := []*shared.Warning{}

		lex := lexer.NewLexer(fileName, fileContent)
		res := lex.Parse()

		errors = append(errors, res.Errors...)
		warnings = append(warnings, res.Warnings...)

		if len(res.Errors) == 0 {

			parser := parser.NewParser(res)
			parseRes := parser.Parse()

			if parseRes.Err != nil {
				errors = append(errors, parseRes.Err)
			}

			warnings = append(warnings, parseRes.Warnings...)

		}

		if *outJson != "" {
			errorsJson, _ := json.Marshal(errors)
			warningsJson, _ := json.Marshal(warnings)

			fd, _ := os.Create(*outJson)

			fd.WriteString("{\"errors\":")
			fd.Write(errorsJson)
			fd.WriteString(",\"warnings\":")
			fd.Write(warningsJson)
			fd.WriteString("}")

		}

	}

	/*
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

			res, _ := execLexer(fileName, fileContent)

			lexer.PrintLexerResultToHTML(res, htmlOutputFileName)

		} else if command == "parse" {

			if len(os.Args) < 3 {
				// fmt.Println("Usage: rt parse <filename> [<html output filename>]")
			}

			htmlFile := ""
			if len(os.Args) == 4 {
				htmlFile := os.Args[3]
			}

			fileName := os.Args[2]
			htmlFile := os.Args[3]

			fileContent := readFileContent(fileName)

			res, _ := execLexer(fileName, fileContent)

			if len(res.Errors) > 0 {
				return
			}

			parseRes := execParser(res, fileName)

			if htmlFile != "" {
				if parseRes.Err == nil {
					parser.PrintParseResultToHTML(parseRes, htmlFile)
				}
			}

		}*/

}
