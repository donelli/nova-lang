package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"recital_lsp/pkg/interpreter"
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

func printUsage() {
	fmt.Println("Usage: nova <subcommand>")
	fmt.Println("Subcommands:")
	fmt.Println("   run <file>          	   run a program")
	fmt.Println("   parse <file> [options]  	parse a file")
	fmt.Println("   	 Options:")
	fmt.Println("   	 	-json <file>   File to print the result as JSON")
	fmt.Println("   	 	-html <file>   File to print the result as HTML")
}

func main() {

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	if command == "parse" {

		if len(os.Args) < 3 {
			printUsage()
			return
		}

		fileName := os.Args[2]

		parseCmd := flag.NewFlagSet("nova parse", flag.ExitOnError)
		outHtml := parseCmd.String("html", "", "")
		outJson := parseCmd.String("json", "", "")

		parseCmd.Parse(os.Args[3:])

		fileContent := readFileContent(fileName)

		errors := []*shared.Error{}
		warnings := []*shared.Warning{}

		lex := lexer.NewLexer(fileName, fileContent)
		res := lex.Parse()

		errors = append(errors, res.Errors...)
		warnings = append(warnings, res.Warnings...)

		var parseRes *parser.ParseResult

		if len(res.Errors) == 0 {

			parser := parser.NewParser(res)
			parseRes = parser.Parse()

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

		} else {

			if len(errors) > 0 {
				fmt.Printf("Errors: %v\n", errors)
			}

			if len(warnings) > 0 {
				fmt.Printf("Warnings: %v\n", warnings)
			}

		}

		if *outHtml != "" && parseRes != nil {
			parser.PrintParseResultToHTML(parseRes, *outHtml)
		}

	} else if command == "run" {

		if len(os.Args) < 3 {
			printUsage()
			return
		}

		fileName := os.Args[2]

		fileContent := readFileContent(fileName)

		errors := []*shared.Error{}
		warnings := []*shared.Warning{}

		lex := lexer.NewLexer(fileName, fileContent)
		res := lex.Parse()

		errors = append(errors, res.Errors...)
		warnings = append(warnings, res.Warnings...)

		var parseRes *parser.ParseResult

		if len(res.Errors) == 0 {

			parser := parser.NewParser(res)
			parseRes = parser.Parse()

			if parseRes.Err != nil {
				errors = append(errors, parseRes.Err)
			}

			warnings = append(warnings, parseRes.Warnings...)

			if len(res.Errors) == 0 {

				interp := interpreter.NewInterpreter()
				res := interp.Visit(parseRes.Node)

				if res.Error != nil {
					errors = append(errors, res.Error)
				}

			}

		}

		if len(errors) > 0 {
			fmt.Printf("Errors: %v\n", errors)
		}

		// TODO print errors and warnings

	}

}
