package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"nova-lang/pkg/interpreter"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/parser"
	"nova-lang/pkg/shared"
	"nova-lang/pkg/utils"
	"os"
	"time"
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
	fmt.Println("  help                    Print help")
	fmt.Println()
	fmt.Println("  run <file> [options]    Run a program")
	fmt.Println("    Options:")
	fmt.Println("      --time			 Show time taken to run the program")
	fmt.Println()
	fmt.Println("  parse <file> [options]  Parse a file")
	fmt.Println("    Options:")
	fmt.Println("      --json <file>   File to print the result as JSON")
	fmt.Println("      --html <file>   File to print the result as HTML")
	fmt.Println()
}

func main() {

	start := time.Now()

	if len(os.Args) < 2 || os.Args[1] == "help" {
		printUsage()
		return
	}

	buf := bytes.NewBuffer([]byte{})
	command := os.Args[1]

	if command == "parse" {

		if len(os.Args) < 3 {
			printUsage()
			return
		}

		fileName := os.Args[2]

		parseCmd := flag.NewFlagSet("nova parse", flag.ContinueOnError)
		parseCmd.SetOutput(buf)
		outHtml := parseCmd.String("html", "", "")
		outJson := parseCmd.String("json", "", "")

		flagErr := parseCmd.Parse(os.Args[3:])

		if flagErr != nil {
			printUsage()
			fmt.Println("Error: " + flagErr.Error())
		}

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

		runFlagSet := flag.NewFlagSet("nova run", flag.ContinueOnError)
		runFlagSet.SetOutput(buf)
		showTime := runFlagSet.Bool("time", false, "")

		flagErr := runFlagSet.Parse(os.Args[3:])

		if flagErr != nil {
			fmt.Println("Error: " + flagErr.Error())
			fmt.Println()
			printUsage()
			return
		}

		fileContent := readFileContent(fileName)

		errors := []*shared.Error{}
		warnings := []*shared.Warning{}

		lex := lexer.NewLexer(fileName, fileContent)
		res := lex.Parse()

		errors = append(errors, res.Errors...)
		warnings = append(warnings, res.Warnings...)

		var parseRes *parser.ParseResult

		if len(errors) == 0 {

			parser := parser.NewParser(res)
			parseRes = parser.Parse()

			if parseRes.Err != nil {
				errors = append(errors, parseRes.Err)
			}

			warnings = append(warnings, parseRes.Warnings...)

			if len(errors) == 0 {

				interpStart := time.Now()

				interp := interpreter.NewInterpreter()
				res := interp.Start(parseRes.Node)

				if res.Error != nil {
					errors = append(errors, res.Error)
				}

				if *showTime {
					fmt.Printf("\nProgram interpreted in %.2f ms", float32(time.Since(interpStart).Microseconds()/1000))
				}

			}

		}

		if len(errors) > 0 {
			fmt.Printf("Errors: %v\n", errors)
		}

		if len(warnings) > 0 {
			fmt.Printf("Warnings: %v\n", warnings)
		}

		if *showTime {
			fmt.Printf("\nTotal time: %.2f ms\n", float32(time.Since(start).Microseconds()/1000))
		}

	}

}
