package lexer

import (
	"fmt"
	"recital_lsp/pkg/shared"
	"recital_lsp/pkg/utils"
	"regexp"
	"strings"
)

// TODO define optimal number to startTokenCount
// Start Token count for the lexer
const startTokenCount = 500

type Lexer struct {
	FileName    string
	FileContent string
	contentLen  int32

	CurrentPosition *shared.Position
	CurrentChar     string
	CurrentRune     rune
	hasCurrentChar  bool
	currentResult   *LexerResult

	dateRegex *regexp.Regexp
}

func (lexer *Lexer) Advance() {
	lexer.CurrentPosition.Advance(lexer.CurrentRune)

	if lexer.CurrentPosition.Index >= lexer.contentLen {
		lexer.hasCurrentChar = false
	} else {
		lexer.CurrentChar = string(lexer.FileContent[lexer.CurrentPosition.Index])
		lexer.CurrentRune = rune(lexer.FileContent[lexer.CurrentPosition.Index])
	}

}

func (lexer *Lexer) PeekNextChar() (rune, bool) {

	index := lexer.CurrentPosition.Index + 1

	if index >= lexer.contentLen {
		return 0, false
	}

	return rune(lexer.FileContent[index]), true

}

func (lexer *Lexer) PeekNextNonEmptyChar() (rune, bool) {

	index := lexer.CurrentPosition.Index + 1

	for lexer.hasCurrentChar && index < lexer.contentLen && strings.Contains(shared.WhitespaceChars, string(lexer.FileContent[index])) {
		lexer.Advance()
	}

	if index >= lexer.contentLen {
		return 0, false
	}

	return rune(lexer.FileContent[index]), true
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

	if lexer.hasCurrentChar && lexer.CurrentRune == '+' {
		lexer.addTokenWithPos(TokenType_PlusPlus, "++", startPos, *lexer.CurrentPosition)
		lexer.Advance()
		return
	}

	lexer.addToken(TokenType_Plus, "+")

}

func (lexer *Lexer) makeMinusToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentRune == '-' {
		lexer.addTokenWithPos(TokenType_MinusMinus, "--", startPos, *lexer.CurrentPosition)
		lexer.Advance()
		return
	}

	lexer.addToken(TokenType_Minus, "-")

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

	for lexer.hasCurrentChar && lexer.CurrentRune != '\n' {
		comment += lexer.CurrentChar
		lexer.Advance()
	}

	lexer.addTokenWithPos(TokenType_Comment, comment, startPos, *lexer.CurrentPosition)

}

func (lexer *Lexer) makeMultiplierOrCommentToken() {

	if lexer.isFirstTokenOfTheLine() {
		lexer.makeComment()
		return
	}

	next, hasNext := lexer.PeekNextChar()

	if hasNext && next == '*' {

		startPos := *lexer.CurrentPosition
		lexer.Advance()

		lexer.addTokenWithPos(TokenType_Exponential, "**", startPos, *lexer.CurrentPosition)
		lexer.Advance()

	} else {
		lexer.addToken(TokenType_Star, "*")
		lexer.Advance()
	}

}

func (lexer *Lexer) reportError(error *shared.Error) {
	lexer.currentResult.AddError(error)
}

func (lexer *Lexer) makeNumber() {
	startPos := *lexer.CurrentPosition
	number := ""
	var dotCount uint8 = 0

	for strings.Contains(shared.DigitsAndDot, lexer.CurrentChar) {

		if lexer.CurrentRune == '.' {
			dotCount++
		}

		number += lexer.CurrentChar
		lexer.Advance()
	}

	if dotCount > 1 {
		lexer.reportError(shared.NewError(startPos, *lexer.CurrentPosition, "Invalid number"))
	}

	lexer.addTokenWithPos(TokenType_Number, number, startPos, *lexer.CurrentPosition)

}

func (lexer *Lexer) matchLastTokenType(tokenType LexerTokenType) bool {

	if lexer.currentResult.TokensCount == 0 {
		return false
	}

	lastToken := lexer.currentResult.Tokens[lexer.currentResult.TokensCount-1]

	return lastToken.Type == tokenType
}

func (lexer *Lexer) matchLastTokenTypeAndValue(tokenType LexerTokenType, value string) bool {

	if lexer.currentResult.TokensCount == 0 {
		return false
	}

	lastToken := lexer.currentResult.Tokens[lexer.currentResult.TokensCount-1]

	return lastToken.Type == tokenType && lastToken.Value == value
}

func (lexer *Lexer) makeIdentifierOrKeyword() {
	startPos := *lexer.CurrentPosition
	identifier := ""

	for strings.Contains(shared.VariableChars, lexer.CurrentChar) {
		identifier += lexer.CurrentChar
		lexer.Advance()
	}

	identifierLower := strings.ToLower(identifier)

	if realKeword, ok := shared.KeywordsMap[identifierLower]; ok {
		lexer.addTokenWithPos(TokenType_Keyword, realKeword, startPos, *lexer.CurrentPosition)
		return
	}

	lexer.addTokenWithPos(TokenType_Identifier, identifier, startPos, *lexer.CurrentPosition)

}

func (lexer *Lexer) makeEqualsToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentRune == '=' {
		lexer.addTokenWithPos(TokenType_EqualsEquals, "==", startPos, *lexer.CurrentPosition)
		lexer.Advance()
	} else {
		lexer.addTokenWithPos(TokenType_Equals, "=", startPos, startPos)
	}

}

func (lexer *Lexer) makeString() {

	startPos := *lexer.CurrentPosition
	stringValue := ""

	startChar := lexer.CurrentRune
	endChar := startChar
	if startChar == '[' {
		endChar = ']'
	}

	lexer.Advance()

	lastChar := ' '

	for lexer.hasCurrentChar {

		if lexer.CurrentRune == endChar {
			break
		}

		if lexer.CurrentRune == '\n' {

			if lastChar != ';' {
				break
			}

			stringValue = strings.TrimSuffix(stringValue, " ")
			stringValue = strings.TrimSuffix(stringValue, ";")

			stringValue += " "

			lastChar = lexer.CurrentRune
			lexer.Advance()
			continue

		}

		if lexer.CurrentRune != ' ' {
			lastChar = lexer.CurrentRune
		}

		stringValue += lexer.CurrentChar
		lexer.Advance()
	}

	if lexer.CurrentRune != endChar {
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
		lexer.addTokenWithPos(TokenType_LeftBracket, "[", pos, pos)
		return
	}

	lexer.makeString()

}

func (lexer *Lexer) makeLessThenEqualsOrNotToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentRune == '=' {
		lexer.addTokenWithPos(TokenType_LessThanEqual, "<=", startPos, *lexer.CurrentPosition)
		lexer.Advance()
	} else if lexer.hasCurrentChar && lexer.CurrentRune == '>' {
		lexer.addTokenWithPos(TokenType_NotEqual, "<>", startPos, *lexer.CurrentPosition)
		lexer.Advance()
	} else {
		lexer.addTokenWithPos(TokenType_LessThan, "<", startPos, startPos)
	}

}

func (lexer *Lexer) makeGreaterThanEqualsToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentRune == '=' {
		lexer.addTokenWithPos(TokenType_GreaterThanEqual, ">=", startPos, *lexer.CurrentPosition)
		lexer.Advance()
	} else {
		lexer.addTokenWithPos(TokenType_GreaterThan, ">", startPos, startPos)
	}

}

func (lexer *Lexer) makeBoolOrDotToken() {

	nextChars := string(lexer.getNextBytes(2))

	if nextChars == "t." || nextChars == "f." {

		startPos := *lexer.CurrentPosition
		boolValue := lexer.CurrentChar + nextChars
		lexer.Advance()
		lexer.Advance()
		lexer.Advance()
		lexer.addTokenWithPos(TokenType_Boolean, boolValue, startPos, *lexer.CurrentPosition)

		return
	}

	lexer.addToken(TokenType_Dot, ".")
	lexer.Advance()

}

func (lexer *Lexer) makeNotOrNotEqualsToken() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	if lexer.hasCurrentChar && lexer.CurrentRune == '=' {
		lexer.addTokenWithPos(TokenType_NotEqual, "!=", startPos, *lexer.CurrentPosition)
		lexer.Advance()
	} else {
		lexer.addTokenWithPos(TokenType_Not, "!", startPos, startPos)
	}

}

func (lexer *Lexer) makeDate() {

	startPos := *lexer.CurrentPosition
	lexer.Advance()

	dateValue := ""

	for lexer.hasCurrentChar && lexer.CurrentRune != '}' && lexer.CurrentRune != '\n' {

		dateValue += lexer.CurrentChar
		lexer.Advance()

	}

	if !lexer.hasCurrentChar || lexer.CurrentRune != '}' {
		lexer.reportError(shared.NewError(startPos, *lexer.CurrentPosition, "Expected '}' to close date"))
		lexer.Advance()
		return
	}

	if dateValue != "" && !lexer.dateRegex.Match([]byte(dateValue)) {
		lexer.reportError(shared.NewError(startPos, *lexer.CurrentPosition, "Invalid date"))
	}

	lexer.addTokenWithPos(TokenType_Date, dateValue, startPos, *lexer.CurrentPosition)
	lexer.Advance()

}

func (lexer *Lexer) makeDividerOrCommentToken() {

	nextChar, hasChar := lexer.PeekNextChar()

	if hasChar && nextChar == '/' {
		lexer.makeComment()
		return
	}

	lexer.addToken(TokenType_Slash, "/")
	lexer.Advance()

}

func (lexer *Lexer) makeCommentOrMacro() {

	nextRune, hasNextRune := lexer.PeekNextChar()

	if hasNextRune && lexer.CurrentRune == nextRune {
		lexer.makeComment()
	} else {
		lexer.addToken(TokenType_Macro, "&")
		lexer.Advance()
	}

}

func (lexer *Lexer) Parse() (*LexerResult, error) {

	result := NewLexerResult()
	lexer.currentResult = result

	for {

		if !lexer.hasCurrentChar {
			break
		}

		if strings.Contains(shared.WhitespaceChars, lexer.CurrentChar) {
			lexer.Advance()
			continue
		}

		switch lexer.CurrentRune {
		case '=':
			lexer.makeEqualsToken()
			continue
		case '(':
			lexer.addToken(TokenType_LeftParenthesis, "(")
		case ')':
			lexer.addToken(TokenType_RightParenthesis, ")")
		case '+':
			lexer.makePlusToken()
			continue
		case ',':
			lexer.addToken(TokenType_Comma, ",")
		case '"', '\'':
			lexer.makeString()
			continue
		case '!':
			lexer.makeNotOrNotEqualsToken()
			continue
		case '[':
			lexer.makeStringOrBracket()
			continue
		case ']':
			lexer.addToken(TokenType_RightBracket, "]")
		case '$':
			lexer.addToken(TokenType_DollarSign, "$")
		case '@':
			lexer.addToken(TokenType_AtSign, "@")
		case '&':
			lexer.makeCommentOrMacro()
			continue
		case '<':
			lexer.makeLessThenEqualsOrNotToken()
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
		case '/':
			lexer.makeDividerOrCommentToken()
			continue
		case '\\':

			nextRune, hasNextRune := lexer.PeekNextChar()

			if hasNextRune && nextRune == '&' {

				startPos := *lexer.CurrentPosition
				lexer.Advance()

				lexer.addTokenWithPos(TokenType_Ampersand, "\\&", startPos, *lexer.CurrentPosition)

				lexer.Advance()
				continue
			}

			lexer.reportError(shared.NewError(*lexer.CurrentPosition, *lexer.CurrentPosition, "Unexpected '\\'"))
			lexer.Advance()

		case '?':
			lexer.addToken(TokenType_QuestionMark, "?")
		case '\n':
			lexer.addToken(TokenType_NewLine, "")
		case '{':
			lexer.makeDate()
			continue
		case '%':
			lexer.addToken(TokenType_Percent, "%")
		case ';':

			lexer.Advance()
			startPos := *lexer.CurrentPosition
			ignoreUntilNewLine := false

			for lexer.hasCurrentChar && lexer.CurrentRune != '\n' {

				if !ignoreUntilNewLine && !strings.Contains(shared.WhitespaceChars, lexer.CurrentChar) {

					nextRune, hasNextRune := lexer.PeekNextChar()

					if lexer.CurrentRune == '&' && hasNextRune && nextRune == '&' {
						ignoreUntilNewLine = true
						continue
					}

					lexer.reportError(shared.NewError(startPos, startPos, "Invalid sintax for ';'"))
					break
				}

				lexer.Advance()
			}

			if lexer.hasCurrentChar && lexer.CurrentRune == '\n' {
				lexer.Advance()
			}

			continue
		case '^':
			lexer.addToken(TokenType_Exponential, "^")
		default:

			if strings.Contains(shared.Digits, lexer.CurrentChar) {
				lexer.makeNumber()
				continue
			}

			if strings.Contains(shared.LettersAndUnderline, lexer.CurrentChar) {
				lexer.makeIdentifierOrKeyword()
				continue
			}

			lexer.reportError(shared.NewError(*lexer.CurrentPosition, *lexer.CurrentPosition, fmt.Sprintf("Invalid character: %s", lexer.CurrentChar)))
		}

		lexer.Advance()

	}

	lexer.addToken(TokenType_EOF, "")

	return lexer.currentResult, nil
}

func NewLexer(fileName string, fileContent string) *Lexer {

	dateRegex, err := regexp.Compile(`\d{0,2}\/\d{0,2}\/\d{2,4}`)
	utils.Assert(err == nil, "Error compiling date regex")

	shared.LoadKeywords()

	return &Lexer{
		FileName:    fileName,
		FileContent: fileContent,
		contentLen:  int32(len(fileContent)),

		CurrentRune:     ' ',
		CurrentChar:     "",
		hasCurrentChar:  true,
		CurrentPosition: shared.NewPosition(),
		currentResult:   nil,

		dateRegex: dateRegex,
	}
}
