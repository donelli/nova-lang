package lexer

import "fmt"

type LexerResult struct {
	Tokens      []*LexerToken
	TokensCount uint32
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
	}
}
