package lexer

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type LexerTokenType uint8

const (
	// TokenType_EOF represents the end of file token
	TokenType_EOF LexerTokenType = iota + 1

	// TokenType_EOF represents the end of line token
	TokenType_NewLine

	// TokenType_Identifier represents an identifier token (variable, function name, ...)
	TokenType_Identifier

	// TokenType_Number represents a number token (integer and float)
	TokenType_Number

	// TokenType_String represents a string token
	TokenType_String

	// TokenType_Plus represents a plus operator token
	TokenType_Plus

	// TokenType_PlusPlus represents a double plus operator token
	TokenType_PlusPlus

	// TokenType_Minus represents a minus operator token
	TokenType_Minus

	// TokenType_MinusMinus represents a double minus operator token
	TokenType_MinusMinus

	// TokenType_Multiply represents a multiply operator token
	TokenType_Star

	// TokenType_Divide represents a divide operator token
	TokenType_Slash

	// TokenType_Modulo represents a modulo operator token (%)
	TokenType_Percent

	// TokenType_Ampersand represents an escaped &
	TokenType_Ampersand

	// TokenType_Macro represents an &
	TokenType_Macro

	// TokenType_Equal represents an =
	TokenType_Equals

	// TokenType_EqualsEquals represents an ==
	TokenType_EqualsEquals

	// TokenType_NotEqual represents an != or <>
	TokenType_NotEqual

	// TokenType_LessThan represents an <
	TokenType_LessThan

	// TokenType_GreaterThan represents an >
	TokenType_GreaterThan

	// TokenType_LessThanEqual represents an <=
	TokenType_LessThanEqual

	// TokenType_GreaterThanEqual represents an >=
	TokenType_GreaterThanEqual

	// TokenType_LeftParenthesis represents a (
	TokenType_LeftParenthesis

	// TokenType_RightParenthesis represents a )
	TokenType_RightParenthesis

	// TokenType_Comment represents a comment
	TokenType_Comment

	// TokenType_Comma represents a ,
	TokenType_Comma

	// TokenType_LeftBracket represents a [
	TokenType_LeftBracket

	// TokenType_RightBracket represents a ]
	TokenType_RightBracket

	// TokenType_Keyword represents a keyword (reserved word)
	TokenType_Keyword

	// TokenType_QuestionMark represents a ?
	TokenType_QuestionMark

	// TokenType_DollarSign represents a $
	TokenType_DollarSign

	// TokenType_AtSign represents a @
	TokenType_AtSign

	// TokenType_Boolean represents a boolean (.t. or .f.)
	TokenType_Boolean

	// TokenType_Dot represents a .
	TokenType_Dot

	// TokenType_Not represents a ! or not
	TokenType_Not

	// TokenType_Date represents a date: {} {01/01/21} {01/01/2021}
	TokenType_Date

	// TokenType_Path represents a path: '/home/user/file.txt'
	TokenType_Path

	// TokenType_Not represents a ^ or **
	TokenType_Exponential

	// TokenType_Arrrow represents a ->
	TokenType_Arrrow

	// TokenType_Pipe represents a |
	TokenType_Pipe
)

//go:generate stringer -type=LexerTokenType -trimprefix=TokenType_

type LexerToken struct {
	Range *shared.Range
	Type  LexerTokenType
	Value string
}

func (lexerToken *LexerToken) MatchType(tokenType LexerTokenType) bool {
	return lexerToken.Type == tokenType
}

func (lexerToken *LexerToken) Match(tokenType LexerTokenType, value string) bool {
	return lexerToken.Type == tokenType && lexerToken.Value == value
}

func (lexerToken *LexerToken) MatchMultiple(tokenType LexerTokenType, values []string) bool {

	if lexerToken.Type != tokenType {
		return false
	}

	for _, value := range values {
		if lexerToken.Value == value {
			return true
		}
	}

	return false
}

func (lexerToken *LexerToken) String() string {

	if lexerToken.Value != "" {
		return fmt.Sprintf("Tok{Range: %v, Type: %v, Val: %v}", lexerToken.Range, lexerToken.Type, lexerToken.Value)
	} else {
		return fmt.Sprintf("Tok{Range: %v, Type: %v}", lexerToken.Range, lexerToken.Type)
	}

}

func NewLexerToken(startPos shared.Position, endPos shared.Position, tokenType LexerTokenType, value string) *LexerToken {
	return &LexerToken{
		Range: shared.NewRange(startPos, endPos),
		Type:  tokenType,
		Value: value,
	}
}
