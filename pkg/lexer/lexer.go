package lexer

import (
	"errors"
)

type Lexer struct {
	FileName    string
	FileContent string
}

type LexerResult struct {
}

func (lexer *Lexer) Parse() (LexerResult, error) {
	return LexerResult{}, errors.New("Not implemented")
}

func NewLexer(fileName string, fileContent string) *Lexer {
	return &Lexer{
		FileName:    fileName,
		FileContent: fileContent,
	}
}
