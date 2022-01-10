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
	"nova-lang/pkg/testing"
	"nova-lang/pkg/utils"
	"os"
	"time"
)

func readFileContent(fileName string) string {

	dat, err := os.ReadFile(fileName)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	dat = utils.NormalizeNewlines(dat)

	return string(dat)
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: nova <subcommand>")
	fmt.Fprintln(os.Stderr, "Subcommands:")
	fmt.Fprintln(os.Stderr, "  help                    Print help")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "  run <file> [options]    Run a program")
	fmt.Fprintln(os.Stderr, "    Options:")
	fmt.Fprintln(os.Stderr, "      --time              Show time taken to run the program")
	fmt.Fprintln(os.Stderr, "      --simulation        Run the program in simulation mode (no screen writing)")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "  parse <file> [options]  Parse a file")
	fmt.Fprintln(os.Stderr, "    Options:")
	fmt.Fprintln(os.Stderr, "      --json <file>       File to print the result as JSON")
	fmt.Fprintln(os.Stderr, "      --html <file>       File to print the result as HTML")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "  test <command>          Testing tool for the development of the language")
	fmt.Fprintln(os.Stderr, "    Commands:")
	fmt.Fprintln(os.Stderr, "      new <.prg file>     Creates a new test file for that program")
	fmt.Fprintln(os.Stderr, "      update <test_name>  Updates the test expected output")
	fmt.Fprintln(os.Stderr, "      run <test_name>     Run that test and compares with the expected output")
	fmt.Fprintln(os.Stderr, "      all                 Run all tests and compares each one with the expected output")
	fmt.Fprintln(os.Stderr)
}

func main() {

	start := time.Now()

	if len(os.Args) < 2 || os.Args[1] == "help" {
		printUsage()
		os.Exit(1)
	}

	buf := bytes.NewBuffer([]byte{})
	command := os.Args[1]

	if command == "parse" {

		if len(os.Args) < 3 {
			printUsage()
			os.Exit(1)
		}

		fileName := os.Args[2]

		parseCmd := flag.NewFlagSet("nova parse", flag.ContinueOnError)
		parseCmd.SetOutput(buf)
		outHtml := parseCmd.String("html", "", "")
		outJson := parseCmd.String("json", "", "")

		flagErr := parseCmd.Parse(os.Args[3:])

		if flagErr != nil {
			printUsage()
			fmt.Fprintf(os.Stderr, "Error: "+flagErr.Error())
			os.Exit(1)
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

			parser := parser.NewParser(res, false)
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

		if len(errors) > 0 {
			os.Exit(1)
		}

	} else if command == "run" {

		if len(os.Args) < 3 {
			printUsage()
			os.Exit(1)
		}

		fileName := os.Args[2]

		runFlagSet := flag.NewFlagSet("nova run", flag.ContinueOnError)
		runFlagSet.SetOutput(buf)
		showTime := runFlagSet.Bool("time", false, "")
		simulationMode := runFlagSet.Bool("simulation", false, "")

		flagErr := runFlagSet.Parse(os.Args[3:])

		if flagErr != nil {
			fmt.Println("Error: " + flagErr.Error())
			printUsage()
			os.Exit(1)
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

			parser := parser.NewParser(res, false)
			parseRes = parser.Parse()

			if parseRes.Err != nil {
				errors = append(errors, parseRes.Err)
			}

			warnings = append(warnings, parseRes.Warnings...)

			if len(errors) == 0 {

				interpStart := time.Now()

				interp := interpreter.NewInterpreter()
				res := interp.Start(parseRes.Node, *simulationMode)

				if res.Error != nil {
					errors = append(errors, res.Error)
				}

				if *showTime {
					fmt.Printf("\nProgram interpreted in %d ms", time.Since(interpStart).Milliseconds())
				}

			}

		}

		if len(errors) > 0 {

			for _, err := range errors {
				fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err.StringWithProgram(fileName))
			}

		}

		if len(warnings) > 0 {
			for _, warn := range warnings {
				fmt.Fprintf(os.Stdout, "[WARN] %s\n", warn.Message)
			}
		}

		if *showTime {
			fmt.Printf("\nTotal time: %d ms\n", time.Since(start).Milliseconds())
		}

		if len(errors) > 0 {
			os.Exit(1)
		}

	} else if command == "test" {

		err := testing.ParseAndExecTests()
		if err != "" {
			fmt.Fprintf(os.Stderr, "[ERROR] %s", err)
			os.Exit(1)
		}

	} else {

		fmt.Fprintf(os.Stderr, "[ERROR] Unknown command '%s'\n", command)
		os.Exit(1)

	}

}
