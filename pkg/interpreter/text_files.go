package interpreter

import (
	"bufio"
	"os"
)

type OpenedTextFile struct {
	File   *os.File
	Reader *bufio.ReadWriter
}

var InterpreterOpenedFiles map[string]*OpenedTextFile

func OpenTextFile(fileName string, fileMode int) int {

	flags := -1

	if fileMode == 0 {
		flags = os.O_RDONLY
	} else if fileMode == 1 {
		flags = os.O_WRONLY
	} else if fileMode == 2 {
		flags = os.O_RDWR
	} else {
		return -1
	}

	f, err := os.OpenFile(fileName, flags, 0666)
	if err != nil {
		return -1
	}

	InterpreterOpenedFiles[]

}
