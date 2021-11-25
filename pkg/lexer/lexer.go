package lexer

import (
	"fmt"
	"recital_lsp/pkg/shared"
	"strings"
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

func (lexer *Lexer) isFirstTokenOfTheLine() bool {

	if lexer.currentResult.TokensCount == 0 {
		return true
	}

	lastToken := lexer.currentResult.Tokens[lexer.currentResult.TokensCount-1]

	if lastToken.Range.End.Row != lexer.CurrentPosition.Row {
		return true
	}

	return false
}

func (lexer *Lexer) makeComment() {

	startPos := *lexer.CurrentPosition
	comment := ""

	// TODO some wierd caracter is in the comment

	for lexer.hasCurrentChar && lexer.CurrentChar != '\n' {
		comment += string(lexer.CurrentChar)
		lexer.Advance()
	}

	lexer.addTokenWithPos(TokenType_Comment, comment, startPos, *lexer.CurrentPosition)

}

func (lexer *Lexer) makeMultiplierOrCommentToken() {

	if lexer.isFirstTokenOfTheLine() {
		lexer.makeComment()
		return
	}

	lexer.addToken(TokenType_Star, "")
	lexer.Advance()

}

func (lexer *Lexer) reportError(error shared.Error) {
	lexer.currentResult.AddError(error)
}

func (lexer *Lexer) makeNumber() {
	startPos := *lexer.CurrentPosition
	number := ""
	var dotCount uint8 = 0

	for strings.Contains(shared.DigitsAndDot, string(lexer.CurrentChar)) {

		if lexer.CurrentChar == '.' {
			dotCount++
		}

		number += string(lexer.CurrentChar)
		lexer.Advance()
	}

	if dotCount > 1 {
		lexer.reportError(shared.NewError(startPos, *lexer.CurrentPosition, "Invalid number"))
	}

	lexer.addToken(TokenType_Number, number)

}

func (lexer *Lexer) matchLastTokenType(tokenType LexerTokenType) bool {

	if lexer.currentResult.TokensCount == 0 {
		return false
	}

	lastToken := lexer.currentResult.Tokens[lexer.currentResult.TokensCount-1]

	if lastToken.Type != tokenType {
		return false
	}

	return true
}

func (lexer *Lexer) makeIdentifier() {
	startPos := *lexer.CurrentPosition
	identifier := ""

	for strings.Contains(shared.LettersAndUnderline, string(lexer.CurrentChar)) {
		identifier += string(lexer.CurrentChar)
		lexer.Advance()
	}

	lexer.addTokenWithPos(TokenType_Identifier, identifier, startPos, *lexer.CurrentPosition)

}

func (lexer *Lexer) Parse() (*LexerResult, error) {

	result := NewLexerResult()
	lexer.currentResult = result

	for {

		if !lexer.hasCurrentChar {
			break
		}

		if lexer.CurrentChar == ' ' || lexer.CurrentChar == '\t' || lexer.CurrentChar == '\r' {
			lexer.Advance()
			continue
		}

		switch lexer.CurrentChar {
		case '\n':

			// TODO ignore multiple new lines?

			lexer.addToken(TokenType_NewLine, "")

		case '+':
			lexer.makePlusToken()
			continue
		case '-':
			lexer.makeMinusToken()
			continue
		case '*':
			lexer.makeMultiplierOrCommentToken()
			continue
		default:

			if strings.Contains(shared.Digits, string(lexer.CurrentChar)) {
				lexer.makeNumber()
				continue
			}

			if strings.Contains(shared.LettersAndUnderline, string(lexer.CurrentChar)) {
				lexer.makeIdentifier()
				continue
			}

			fmt.Printf("Unknown char: %+q\n", lexer.CurrentChar)

			lexer.reportError(shared.NewError(*lexer.CurrentPosition, *lexer.CurrentPosition, fmt.Sprintf("Invalid character: %s", string(lexer.CurrentChar))))
		}

		lexer.Advance()

	}

	return lexer.currentResult, nil
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
