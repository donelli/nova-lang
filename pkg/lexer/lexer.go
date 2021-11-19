package lexer

import (
	"recital_lsp/pkg/shared"
)

// Start Token count for the lexer
const startTokenCount = 500

type Lexer struct {
	FileName    string
	FileContent string
	contentLen  int32

	CurrentPosition *shared.Position
	CurrentChar     rune
	hasCurrentChar  bool
	currentResult   *LexerResult
}

func (lexer *Lexer) Advance() {
	lexer.CurrentPosition.Advance(lexer.CurrentChar)

	if lexer.CurrentPosition.Index >= lexer.contentLen {
		lexer.hasCurrentChar = false
	} else {
		lexer.CurrentChar = rune(lexer.FileContent[lexer.CurrentPosition.Index])
	}

}

func (lexer *Lexer) addToken(tokenType LexerTokenType, value string) {
	lexer.currentResult.AddToken(NewLexerToken(*lexer.CurrentPosition, *lexer.CurrentPosition, tokenType, value))
}

func (lexer *Lexer) addTokenWithPos(tokenType LexerTokenType, value string, startPos shared.Position, endPos shared.Position) {
	lexer.currentResult.AddToken(NewLexerToken(startPos, endPos, tokenType, value))
}

func (lexer *Lexer) makePlusToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentChar == '+' {
		lexer.addTokenWithPos(TokenType_PlusPlus, "", startPos, *lexer.CurrentPosition)
		lexer.Advance()
		return
	}

	lexer.addToken(TokenType_Plus, "")

}

func (lexer *Lexer) makeMinusToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentChar == '-' {
		lexer.addTokenWithPos(TokenType_MinusMinus, "", startPos, *lexer.CurrentPosition)
		lexer.Advance()
		return
	}

	lexer.addToken(TokenType_Minus, "")

}

func (lexer *Lexer) Parse() (*LexerResult, error) {

	result := NewLexerResult()
	lexer.currentResult = result

	for {

		if !lexer.hasCurrentChar {
			break
		}

		if lexer.CurrentChar == ' ' || lexer.CurrentChar == '\t' {
			lexer.Advance()
			continue
		}

		switch lexer.CurrentChar {
		case '\n':
			lexer.addToken(TokenType_NewLine, "")
		case '+':
			lexer.makePlusToken()
			continue
		case '-':
			lexer.makeMinusToken()
			continue
		}

		lexer.Advance()

	}

	return result, nil
}

func NewLexer(fileName string, fileContent string) *Lexer {
	return &Lexer{
		FileName:        fileName,
		FileContent:     fileContent,
		contentLen:      int32(len(fileContent)),
		CurrentChar:     ' ',
		hasCurrentChar:  true,
		CurrentPosition: shared.NewPosition(),
		currentResult:   nil,
	}
}
