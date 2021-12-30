package lexer

import (
	"fmt"
	"nova-lang/pkg/shared"
)

type LexerResult struct {
	Tokens      []*LexerToken
	TokensCount uint32

	Errors   []*shared.Error
	Warnings []*shared.Warning
}

func (LexerResult *LexerResult) AddError(error *shared.Error) {
	LexerResult.Errors = append(LexerResult.Errors, error)
}

func (LexerResult *LexerResult) AddToken(token *LexerToken) {
	LexerResult.Tokens = append(LexerResult.Tokens, token)
	LexerResult.TokensCount++
}

func (LexerResult *LexerResult) String() string {

	str := "LexerResult{\n"

	for pointer := range LexerResult.Tokens {
		tok := LexerResult.Tokens[pointer]
		str += "  " + fmt.Sprintf("%+v", tok) + "\n"
	}

	return str + "}"
}

func NewLexerResult() *LexerResult {
	return &LexerResult{
		Tokens:      make([]*LexerToken, 0, startTokenCount),
		TokensCount: 0,
		Errors:      make([]*shared.Error, 0),
		Warnings:    make([]*shared.Warning, 0),
	}
}
