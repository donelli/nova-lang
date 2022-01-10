package testing

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Test struct {
	Program string `json:"program"`
	Output  string `json:"output"`
}

func runTest(programName string) (string, error) {

	fmt.Printf("[INFO] Running `%s` test program...\n", programName)

	cmd := exec.Command(os.Args[0], "run", programName, "--simulation")

	stdout, err := cmd.Output()

	return string(stdout), err
}

func cmdCreateNewTest(programName string) string {

	fmt.Println("[INFO] Creating new test. Program: " + programName)

	if _, err := os.Stat(programName); errors.Is(err, os.ErrNotExist) {
		return fmt.Sprintf("File '%s' doesn't exist\n", programName)
	}

	commandOutput, err := runTest(programName)

	if err != nil {
		return fmt.Sprintln("Error running command: ", err)
	}

	programFolder := filepath.Dir(programName)
	programFileName := strings.TrimSuffix(filepath.Base(programName), ".prg")

	testFileName := programFolder + "/" + programFileName + ".json"

	if _, err := os.Stat(testFileName); !errors.Is(err, os.ErrNotExist) {
		return fmt.Sprintf("Test `%s` already exists\n", testFileName)
	}

	fmt.Printf("[INFO] Creating `%v`...\n", testFileName)

	test := Test{
		Program: programName,
		Output:  commandOutput,
	}

	fp, err := os.Create(testFileName)
	if err != nil {
		return fmt.Sprintln("Error creating file: ", err)
	}

	jsonText, err := json.MarshalIndent(test, "", "  ")
	if err != nil {
		return fmt.Sprintln("Error marshalling json: ", err)
	}

	fp.Write(jsonText)
	fp.Close()

	fmt.Println("[INFO] Test created successfully")

	return ""
}

func cmdRunAllTests() {

	files, err := os.ReadDir("./tests")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading tests folder: ", err)
		os.Exit(1)
	}

	testsCount := 0
	totalFailedTests := 0

	for _, file := range files {

		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		testName := "./tests/" + file.Name()

		msg, testPassed := cmdRunTest(testName)

		testsCount++

		if !testPassed {
			fmt.Fprintln(os.Stderr, "[ERROR] ", msg)
			totalFailedTests++
		}

	}

	fmt.Println()
	fmt.Println("[INFO] Tests run:    ", testsCount)
	fmt.Println("[INFO] Tests failed: ", totalFailedTests)
	fmt.Println("[INFO] Tests passed: ", testsCount-totalFailedTests)

}

func cmdUpdateTestOutput(testName string) string {

	jsonTestName := strings.TrimSuffix(testName, ".prg") + ".json"

	if _, err := os.Stat(jsonTestName); errors.Is(err, os.ErrNotExist) {
		return fmt.Sprintf("Test file '%s' doesn't exist\n", jsonTestName)
	}

	contentBytes, err := os.ReadFile(jsonTestName)
	if err != nil {
		return fmt.Sprintln("Error reading file: ", err)
	}

	var test Test
	err = json.Unmarshal(contentBytes, &test)
	if err != nil {
		return fmt.Sprintln("Error unmarshalling json: ", err)
	}

	commandOutput, err := runTest(test.Program)

	if err != nil {
		return fmt.Sprintln("Error running command: ", err)
	}

	test.Output = commandOutput

	fp, err := os.Create(jsonTestName)
	if err != nil {
		return fmt.Sprintln("Error creating file: ", err)
	}

	jsonText, err := json.MarshalIndent(test, "", "  ")
	if err != nil {
		return fmt.Sprintln("Error marshalling json: ", err)
	}

	fp.Write(jsonText)
	fp.Close()

	fmt.Println("[INFO] Test updated successfully")
	return ""
}

func cmdRunTest(testName string) (string, bool) {

	jsonTestName := testName

	if !strings.HasSuffix(jsonTestName, ".json") {
		jsonTestName = strings.TrimSuffix(testName, ".prg") + ".json"
	}

	if _, err := os.Stat(jsonTestName); errors.Is(err, os.ErrNotExist) {
		return fmt.Sprintf("Test file '%s' doesn't exist\n", jsonTestName), false
	}

	contentBytes, err := os.ReadFile(jsonTestName)
	if err != nil {
		return fmt.Sprintln("Error reading file: ", err), false
	}

	var test Test
	err = json.Unmarshal(contentBytes, &test)
	if err != nil {
		return fmt.Sprintln("Error unmarshalling json: ", err), false
	}

	commandOutput, err := runTest(test.Program)

	if err != nil {
		return fmt.Sprintln("Error running command: ", err), false
	}

	if commandOutput != test.Output {
		fmt.Fprintf(os.Stderr, "[ERROR] Test failed: %s\n", test.Program)
		fmt.Fprintln(os.Stderr, "----------------------------------")
		fmt.Fprintln(os.Stderr, "Expected:")
		fmt.Fprintln(os.Stderr, test.Output)
		fmt.Fprintln(os.Stderr, "----------------------------------")
		fmt.Fprintln(os.Stderr, "Actual:")
		fmt.Fprintln(os.Stderr, commandOutput)
		return "", false
	} else {
		fmt.Println("[INFO] Test passed: ", test.Program)
	}

	return "", true
}

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

		return cmdCreateNewTest(programName)

	} else if testCommand == "run" {

		if len(os.Args) < 4 {
			return "Expected test name\n"
		}

		testName := os.Args[3]

		err, _ := cmdRunTest(testName)
		return err

	} else if testCommand == "update" {

		if len(os.Args) < 4 {
			return "Expected test name\n"
		}

		testName := os.Args[3]

		return cmdUpdateTestOutput(testName)

	} else if testCommand == "all" {

		cmdRunAllTests()
		return ""

	}

	return fmt.Sprintf("Unknown test subcommand '%s'\n", testCommand)
}
