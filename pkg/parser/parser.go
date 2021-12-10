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

func NewParser(lexerResult *lexer.LexerResult) *Parser {
	parser := &Parser{
		LexerResult:       lexerResult,
		CurrentTokenIndex: -1,
	}
	parser.advance()
	return parser
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

func (p *Parser) parseStatement() *ParseResult {

	return nil
}

func (p *Parser) parseStatements() *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start.Copy()
	statements := make([]*Node, 10)

	for p.CurrentToken.Type == lexer.TokenType_NewLine {
		res.RegisterAdvancement()
		p.advance()
	}

	statement := res.Register(p.parseStatement())

	if res.Err != nil {
		return res
	}

	statements = append(statements, statement)
	moreStatements := true

	for {

		newlineCount := 0
		for p.CurrentToken.Type == lexer.TokenType_NewLine {
			res.RegisterAdvancement()
			p.advance()
			newlineCount += 1
		}

		if newlineCount == 0 {
			moreStatements = false
		}

	}

	return res.Success(NewListNode(statements, startPos, &p.CurrentToken.Range.End))
}
