package lexer

import (
	"recital_lsp/pkg/shared"
)

type LexerTokenType int

const (
	// TokenType_EOF represents the end of file token.
	TokenType_EOF = iota + 1

	// TokenType_Identifier represents an identifier token (variable, function name, ...)
	TokenType_Identifier

	// TokenType_Number represents a number token (integer and float)
	TokenType_Number

	// TokenType_String represents a string token
	TokenType_String

	// TokenType_Operator represents a plus operator token
	TokenType_Plus

	// TokenType_Minus represents a minus operator token
	TokenType_Minus

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
)

type LexerToken struct {
	Range shared.Range
	Type  LexerTokenType
	Value string
}

// String token

type StringLexerToken struct {
	LexerToken
}

func (token *StringLexerToken) String() string {
	return "str:" + token.Value + token.Range.String()
}

func NewStringToken(tokenRange shared.Range, value string) *StringLexerToken {
	return &StringLexerToken{
		LexerToken: LexerToken{
			Range: tokenRange,
			Type:  TokenType_String,
			Value: value,
		},
	}
}
