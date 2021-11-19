package lexer

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type LexerTokenType uint8

const (
	// TokenType_EOF represents the end of file token.
	TokenType_EOF LexerTokenType = iota + 1

	// TokenType_EOF represents the end of line token.
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

	// TokenType_Modulo represents a modulo operator token
	TokenType_Percent

	// TokenType_Ampersand represents an @
	TokenType_Ampersand

	// TokenType_Equal represents an =
	TokenType_Equal

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
)

//go:generate stringer -type=LexerTokenType -trimprefix=TokenType_

type LexerToken struct {
	Range *shared.Range
	Type  LexerTokenType
	Value string
}

func (this *LexerToken) String() string {

	valueString := ""
	if this.Value != "" {
		valueString = fmt.Sprintf(", Val: %v", this.Value)
	}

	return fmt.Sprintf("Tok{Range: %v, Type: %v%s}", this.Range, this.Type, valueString)
}

func NewLexerToken(startPos shared.Position, endPos shared.Position, tokenType LexerTokenType, value string) *LexerToken {
	return &LexerToken{
		Range: shared.NewRange(startPos, endPos),
		Type:  tokenType,
		Value: value,
	}
}
