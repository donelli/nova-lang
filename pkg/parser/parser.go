package parser

import (
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type Parser struct {
	LexerResult       *lexer.LexerResult
	CurrentTokenIndex int
	CurrentToken      *lexer.LexerToken
}

func (p *Parser) advance() {
	p.CurrentTokenIndex++
	p.updateCurrentToken()
}

func (p *Parser) reverse() {
	p.CurrentTokenIndex--
	p.updateCurrentToken()
}

func (p *Parser) reverseAmount(amount int) {
	p.CurrentTokenIndex -= amount
	p.updateCurrentToken()
}

func (p *Parser) updateCurrentToken() {
	if p.CurrentTokenIndex >= 0 && p.CurrentTokenIndex < len(p.LexerResult.Tokens) {
		p.CurrentToken = p.LexerResult.Tokens[p.CurrentTokenIndex]
	}
}

func (p *Parser) Parse() *ParseResult {
	res := p.parseStatements()

	if res.Err != nil && p.CurrentToken.Type != lexer.TokenType_EOF {
		return res.Failure(shared.NewInvalidSyntaxError(p.CurrentToken.Range.Start, p.CurrentToken.Range.End, "Token cannot appear after previous tokens"))
	}

	return res
}

func (p *Parser) parseStatements() *ParseResult {

	res := NewParseResult()
	// startPos := p.CurrentToken.Range.Start.Copy()
	// statements := []*StatementNode{}

	return res
	// return res.Success(NewListNode(statements, startPos, p.CurrentToken.Range.End))
}
