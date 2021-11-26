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

func (lexer *Lexer) getNextBytes(charCount int32) []byte {

	nextBytes := make([]byte, 0, charCount)

	index := lexer.CurrentPosition.Index + 1
	endIndex := index + charCount

	for index < lexer.contentLen && index < endIndex {
		nextBytes = append(nextBytes, lexer.FileContent[index])
		index++
	}

	return nextBytes
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

	return lastToken.Range.End.Row != lexer.CurrentPosition.Row
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

	// TODO check if identifier is a keyword
	// Store keywords in a json, then load and store in a map in the first time ou initialization

	if identifier == "say" || identifier == "if" {
		lexer.addTokenWithPos(TokenType_Keyword, identifier, startPos, *lexer.CurrentPosition)
		return
	}

	lexer.addTokenWithPos(TokenType_Identifier, identifier, startPos, *lexer.CurrentPosition)

}

func (lexer *Lexer) makeEqualsToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentChar == '=' {
		lexer.addTokenWithPos(TokenType_EqualsEquals, "", startPos, *lexer.CurrentPosition)
		lexer.Advance()
	} else {
		lexer.addTokenWithPos(TokenType_Equals, "", startPos, startPos)
	}

}

func (lexer *Lexer) makeString() {

	startPos := *lexer.CurrentPosition
	stringValue := ""

	startChar := lexer.CurrentChar
	endChar := startChar
	if startChar == '[' {
		endChar = ']'
	}

	lexer.Advance()

	lastChar := ' '

	for lexer.hasCurrentChar {

		if lexer.CurrentChar == endChar {
			break
		}

		if lexer.CurrentChar == '\n' {

			if lastChar != ';' {
				break
			}

			stringValue = strings.TrimSuffix(stringValue, " ")
			stringValue = strings.TrimSuffix(stringValue, ";")

			stringValue += " "

			lastChar = lexer.CurrentChar
			lexer.Advance()
			continue

		}

		if lexer.CurrentChar != ' ' {
			lastChar = lexer.CurrentChar
		}

		stringValue += string(lexer.CurrentChar)
		lexer.Advance()
	}

	if lexer.CurrentChar != endChar {
		lexer.reportError(shared.NewError(startPos, *lexer.CurrentPosition, "String not terminated"))
		return
	}

	lexer.addTokenWithPos(TokenType_String, stringValue, startPos, *lexer.CurrentPosition)
	lexer.Advance()

}

func (lexer *Lexer) makeStringOrBracket() {

	if lexer.matchLastTokenType(TokenType_Identifier) {
		pos := *lexer.CurrentPosition
		lexer.Advance()
		lexer.addTokenWithPos(TokenType_LeftBracket, "", pos, pos)
		return
	}

	lexer.makeString()

}

func (lexer *Lexer) makeLessThenEqualsToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentChar == '=' {
		lexer.addTokenWithPos(TokenType_LessThanEqual, "", startPos, *lexer.CurrentPosition)
		lexer.Advance()
	} else {
		lexer.addTokenWithPos(TokenType_LessThan, "", startPos, startPos)
	}

}

func (lexer *Lexer) makeGreaterThanEqualsToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentChar == '=' {
		lexer.addTokenWithPos(TokenType_GreaterThanEqual, "", startPos, *lexer.CurrentPosition)
		lexer.Advance()
	} else {
		lexer.addTokenWithPos(TokenType_GreaterThan, "", startPos, startPos)
	}

}

func (lexer *Lexer) makeBoolOrDotToken() {

	nextChars := string(lexer.getNextBytes(2))

	if nextChars == ".t" || nextChars == ".f" {

		startPos := *lexer.CurrentPosition
		boolValue := string(lexer.CurrentChar) + nextChars
		lexer.Advance()
		lexer.Advance()
		lexer.addTokenWithPos(TokenType_Boolean, boolValue, startPos, *lexer.CurrentPosition)

		return
	}

	lexer.addToken(TokenType_Dot, "")
	lexer.Advance()

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
		case '=':
			lexer.makeEqualsToken()
			continue
		case '(':
			lexer.addToken(TokenType_LeftParenthesis, "")
		case ')':
			lexer.addToken(TokenType_RightParenthesis, "")
		case '+':
			lexer.makePlusToken()
			continue
		case ',':
			lexer.addToken(TokenType_Comma, "")
		case '"', '\'':
			lexer.makeString()
			continue
		case '[':
			lexer.makeStringOrBracket()
			continue
		case ']':
			lexer.addToken(TokenType_RightBracket, "")
		case '$':
			lexer.addToken(TokenType_DollarSign, "")
		case '@':
			lexer.addToken(TokenType_AtSign, "")
		case '&':
			lexer.addToken(TokenType_Ampersand, "")
		case '<':
			lexer.makeLessThenEqualsToken()
			continue
		case '>':
			lexer.makeGreaterThanEqualsToken()
			continue
		case '-':
			lexer.makeMinusToken()
			continue
		case '.':
			lexer.makeBoolOrDotToken()
			continue
		case '*':
			lexer.makeMultiplierOrCommentToken()
			continue
		case '?':
			lexer.addToken(TokenType_QuestionMark, "")
		case '\n':
			lexer.addToken(TokenType_NewLine, "")
		default:

			if strings.Contains(shared.Digits, string(lexer.CurrentChar)) {
				lexer.makeNumber()
				continue
			}

			if strings.Contains(shared.LettersAndUnderline, string(lexer.CurrentChar)) {
				lexer.makeIdentifier()
				continue
			}

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
