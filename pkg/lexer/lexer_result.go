package lexer

import "fmt"

type LexerResult struct {
	Tokens []*LexerToken
}

func (this *LexerResult) AddToken(token *LexerToken) {
	this.Tokens = append(this.Tokens, token)
}

func (this *LexerResult) String() string {
	return fmt.Sprintf("LexerResult{Tokens: %v}", this.Tokens)
}

func NewLexerResult() *LexerResult {
	return &LexerResult{
		Tokens: make([]*LexerToken, 0, startTokenCount),
	}
}
