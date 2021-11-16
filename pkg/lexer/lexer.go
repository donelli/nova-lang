package lexer

import (
	"errors"
	"recital_lsp/pkg/shared"
)

type Lexer struct {
	FileName    string
	FileContent string

	CurrentPosition *shared.Position
	CurrentChar     rune
}

type LexerResult struct {
}

func (lexer *Lexer) Advance() {
	lexer.CurrentPosition.Advance(lexer.CurrentChar)
	lexer.CurrentChar = rune(lexer.FileContent[lexer.CurrentPosition.Index])
}

func (lexer *Lexer) Parse() (LexerResult, error) {
	return LexerResult{}, errors.New("Not implemented")
}

func NewLexer(fileName string, fileContent string) *Lexer {
	return &Lexer{
		FileName:        fileName,
		FileContent:     fileContent,
		CurrentChar:     ' ',
		CurrentPosition: shared.NewPosition(),
	}
}
