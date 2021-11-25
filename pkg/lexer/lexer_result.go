package lexer

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type LexerResult struct {
	Tokens      []*LexerToken
	TokensCount uint32

	Errors   []shared.Error
	Warnings []shared.Warning
}

func (this *LexerResult) AddError(error shared.Error) {
	this.Errors = append(this.Errors, error)
}

func (this *LexerResult) AddToken(token *LexerToken) {
	this.Tokens = append(this.Tokens, token)
	this.TokensCount++
}

func (this *LexerResult) String() string {

	str := "LexerResult{\n"

	for pointer := range this.Tokens {
		tok := this.Tokens[pointer]
		str += "  " + fmt.Sprintf("%+v", tok) + "\n"
	}

	return str + "}"
}

func NewLexerResult() *LexerResult {
	return &LexerResult{
		Tokens:      make([]*LexerToken, 0, startTokenCount),
		TokensCount: 0,
		Errors:      make([]shared.Error, 0),
		Warnings:    make([]shared.Warning, 0),
	}
}
