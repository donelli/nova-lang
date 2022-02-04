package testing

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Test struct {
	Program string `json:"program"`
	Output  string `json:"output"`
}

func runTest(programName string) string {

	fmt.Printf("[INFO] Running `%s` test program...\n", programName)

	novaPath := os.Args[0]
	if !strings.HasSuffix(novaPath, ".exe") {
		novaPath += ".exe"
	}

	stderrWriter := new(strings.Builder)
	cmd := exec.Command(novaPath, "run", programName, "--test", ">", "/dev/null")
	cmd.Stderr = stderrWriter

	_, err := cmd.Output()

	commandOutput := "" // string(stdout)
	if err != nil {
		commandOutput += string(stderrWriter.String()) + fmt.Sprint("Error running command: ", err, " ")
	}

	dat, err := os.ReadFile("test.txt")
	if err != nil {
		commandOutput += fmt.Sprint("Error reading file: ", err, " ")
	} else {
		commandOutput += string(dat)
	}

	os.Remove("test.txt")

	return commandOutput
}

func cmdCreateNewTest(programName string) string {

	fmt.Println("[INFO] Creating new test. Program: " + programName)

	if _, err := os.Stat(programName); errors.Is(err, os.ErrNotExist) {
		return fmt.Sprintf("File '%s' doesn't exist\n", programName)
	}

	commandOutput := runTest(programName)

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
	fmt.Printf("[INFO] Tests results: passed: %d, failed: %d, total: %d\n", testsCount-totalFailedTests, totalFailedTests, testsCount)

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

	commandOutput := runTest(test.Program)

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

	runProgramStart := time.Now()
	commandOutput := runTest(test.Program)

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
		fmt.Println("[INFO] Test passed:", test.Program, "(time", time.Since(runProgramStart).Milliseconds(), "ms)")
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
