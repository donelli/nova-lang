package testing

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParseAndExecTests() string {

	if len(os.Args) < 3 {
		return "Expected subcommand for test\n"
	}

	testCommand := os.Args[2]

	if testCommand == "new" {

		if len(os.Args) < 4 {
			return "Expected program name\n"
		}

		programName := os.Args[3]

		fmt.Println("Creating new test. Program: " + programName)

		if _, err := os.Stat(programName); errors.Is(err, os.ErrNotExist) {
			return fmt.Sprintf("File '%s' doesn't exist\n", programName)
		}

		programFolder := filepath.Dir(programName)
		programFileName := strings.TrimSuffix(filepath.Base(programName), ".prg")

		testName := programFolder + "/" + programFileName + ".json"

		fmt.Printf("Creating `%v`...\n", testName)

	} else {
		return fmt.Sprintf("Unknown test subcommand '%s'\n", testCommand)
	}

	return ""
}
