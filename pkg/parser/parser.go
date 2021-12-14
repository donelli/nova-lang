package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
	"strconv"
)

func TokenToString(tokenType lexer.LexerTokenType, value string) string {
	return fmt.Sprintf("%s%s", tokenType.String(), value)
}

var andTokenString = TokenToString(lexer.TokenType_Keyword, "and")
var orTokenString = TokenToString(lexer.TokenType_Keyword, "and")

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

	if res.Err == nil && p.CurrentToken.Type != lexer.TokenType_EOF {
		return res.Failure(shared.NewInvalidSyntaxError(p.CurrentToken.Range.Start, p.CurrentToken.Range.End, "Token cannot appear after previous tokens"))
	}

	return res
}

func (p *Parser) parseArithmeticExpr() *ParseResult {
	return p.parseBinaryOperation("term", "term", nil, []lexer.LexerTokenType{
		lexer.TokenType_Plus,
		lexer.TokenType_Minus,
	})
}

func (p *Parser) parseTerm() *ParseResult {
	return p.parseBinaryOperation("factor", "factor", nil, []lexer.LexerTokenType{
		lexer.TokenType_Star,
		lexer.TokenType_Slash,
		lexer.TokenType_Percent,
	})
}

func (p *Parser) parseFactor() *ParseResult {

	res := NewParseResult()
	currentToken := p.CurrentToken

	if currentToken.MatchType(lexer.TokenType_Plus) || currentToken.MatchType(lexer.TokenType_Minus) {

		res.RegisterAdvancement()
		p.advance()

		factor := res.Register(p.parseFactor())

		if res.Err != nil {
			return res
		}

		return res.Success(NewUnaryOperationNode(currentToken, factor))

	}

	// TODO macro here

	return p.parsePower()
}

func (p *Parser) parseCall() *ParseResult {

	res := NewParseResult()

	atom := res.Register(p.parseAtom())

	if res.Err != nil {
		return res
	}

	// Add here funcion calls and arrays calls

	return res.Success(atom)
}

func (p *Parser) parseAtom() *ParseResult {

	res := NewParseResult()
	token := p.CurrentToken

	if token.MatchType(lexer.TokenType_Number) {

		value, err := strconv.ParseFloat(token.Value, 64)

		if err != nil {
			return res.Failure(shared.NewInvalidSyntaxError(token.Range.Start, token.Range.End, "Invalid number"))
		}

		res.RegisterAdvancement()
		p.advance()

		return res.Success(NewNumberNode(value, token))

	} else if token.MatchType(lexer.TokenType_String) {
		// TODO return string node
	} else if token.MatchType(lexer.TokenType_Boolean) {
		// TODO return bool node
	} else if token.MatchType(lexer.TokenType_Identifier) {
		// TODO return var acess node
	} else if token.MatchType(lexer.TokenType_LeftParenthesis) {

		res.RegisterAdvancement()
		p.advance()

		expr := res.Register(p.parseExpression())

		if res.Err != nil {
			return res
		}

		if !p.CurrentToken.MatchType(lexer.TokenType_RightParenthesis) {
			return res.Failure(shared.NewInvalidSyntaxError(p.CurrentToken.Range.Start, p.CurrentToken.Range.End, "Expected ')' after expression"))
		}

		res.RegisterAdvancement()
		p.advance()

		return res.Success(expr)
	}

	return res.Failure(shared.NewInvalidSyntaxError(token.Range.Start, token.Range.End, "Unexpected token"))
}

func (p *Parser) parsePower() *ParseResult {

	return p.parseBinaryOperation("call", "factor", nil, []lexer.LexerTokenType{
		lexer.TokenType_Exponential,
	})
}

func (p *Parser) parseCompareExpr() *ParseResult {

	res := NewParseResult()

	if p.CurrentToken.MatchType(lexer.TokenType_Not) {

		operationToken := p.CurrentToken
		res.RegisterAdvancement()
		p.advance()

		node := res.Register(p.parseCompareExpr())

		if res.Err != nil {
			return res
		}

		return res.Success(NewUnaryOperationNode(operationToken, node))

	}

	node := res.Register(p.parseBinaryOperation("arithExpr", "arithExpr", []string{}, []lexer.LexerTokenType{
		lexer.TokenType_Equals,
		lexer.TokenType_NotEqual,
		lexer.TokenType_LessThan,
		lexer.TokenType_LessThanEqual,
		lexer.TokenType_GreaterThan,
		lexer.TokenType_GreaterThanEqual,
	}))

	if res.Err != nil {
		return res
	}

	return res.Success(node)
}

func (p *Parser) invokeFunction(funcName string) *ParseResult {

	if funcName == "compareExpr" {
		return p.parseCompareExpr()
	} else if funcName == "arithExpr" {
		return p.parseArithmeticExpr()
	} else if funcName == "term" {
		return p.parseTerm()
	} else if funcName == "call" {
		return p.parseCall()
	} else if funcName == "factor" {
		return p.parseFactor()
	}

	panic(fmt.Sprintf("%s is not valid", funcName))
}

func (p *Parser) parseBinaryOperation(leftFuncName string, rightFuncName string, typeValueOptions []string, typeOptions []lexer.LexerTokenType) *ParseResult {

	res := NewParseResult()

	leftRes := res.Register(p.invokeFunction(leftFuncName))

	if res.Err != nil {
		return res
	}

	for {

		isValidOption := false

		if typeValueOptions == nil {

			for opt := range typeOptions {
				if typeOptions[opt] == p.CurrentToken.Type {
					isValidOption = true
					break
				}
			}

		} else {

			tokenStr := TokenToString(p.CurrentToken.Type, p.CurrentToken.Value)

			for opt := range typeValueOptions {
				if tokenStr == typeValueOptions[opt] {
					isValidOption = true
					break
				}
			}

		}

		if !isValidOption {
			break
		}

		operationToken := p.CurrentToken
		res.RegisterAdvancement()
		p.advance()

		rightRes := res.Register(p.invokeFunction(rightFuncName))

		if res.Err != nil {
			return res
		}

		leftRes = NewBinaryOperationNode(leftRes, operationToken, rightRes)

	}

	return res.Success(leftRes)
}

func (p *Parser) parseExpression() *ParseResult {

	res := NewParseResult()

	node := res.Register(p.parseBinaryOperation("compareExpr", "compareExpr", []string{andTokenString, orTokenString}, []lexer.LexerTokenType{}))

	if res.Err != nil {
		return res
	}

	return res.Success(node)
}

func (p *Parser) parseReturn() *ParseResult {

	res := NewParseResult()
	token := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.Type == lexer.TokenType_NewLine || p.CurrentToken.Type == lexer.TokenType_EOF {
		return res.Success(NewReturnNode(nil, &token.Range.Start, &token.Range.End))
	}

	expr := res.Register(p.parseExpression())

	if res.Err != nil {
		return res
	}

	// TODO Check if EOL or Newline

	return res.Success(NewReturnNode(expr, &token.Range.Start, expr.EndPos()))
}

func (p *Parser) parsePrintStdout() *ParseResult {

	// ? <expr>

	res := NewParseResult()
	p.advance()

	expr := res.Register(p.parseExpression())

	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxError(p.CurrentToken.Range.Start, p.CurrentToken.Range.End, fmt.Sprintf("Expected new line (got %s)", p.CurrentToken.Type.String())))
	}

	return res.Success(NewPrintStdoutNode(expr))
}

func (p *Parser) parseSet() *ParseResult {

	// Types of set:
	// - set <keyword> to <value> ?
	// - set <keyword> <value> ?

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start.Copy()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Keyword) {
		return res.Failure(shared.NewInvalidSyntaxError(*startPos, p.CurrentToken.Range.End, "Expected valid configuration name (keyword)"))
	}

	configName := p.CurrentToken.Value
	p.advance()
	res.RegisterAdvancement()

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "to") {

		p.advance()

		if p.CurrentToken.MatchType(lexer.TokenType_NewLine) || p.CurrentToken.MatchType(lexer.TokenType_Comment) {
			fmt.Println("set to empty")
			return res.Success(NewEmptySetNode(configName, startPos, &p.CurrentToken.Range.End))
		}

		// TODO expect file names and paths depending on the config name
		// Example: set procedure to ...

		expr := res.Register(p.parseExpression())

		if expr == nil {
			return res
		}

		return res.Success(NewSetNode(configName, expr, startPos, expr.EndPos()))

	} else {

		if p.CurrentToken.MatchType(lexer.TokenType_NewLine) || p.CurrentToken.MatchType(lexer.TokenType_Comment) {
			fmt.Println("set empty")
			return res.Success(NewEmptySetNode(configName, startPos, &p.CurrentToken.Range.End))
		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "on") {
			fmt.Println("set On")

			endPos := p.CurrentToken.Range.End.Copy()
			p.advance()

			return res.Success(NewBoolSetNode(configName, "on", startPos, endPos))
		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "off") {
			fmt.Println("set off")

			endPos := p.CurrentToken.Range.End.Copy()
			p.advance()

			return res.Success(NewBoolSetNode(configName, "off", startPos, endPos))
		}

		expr := res.Register(p.parseExpression())

		if expr == nil {
			return res
		}

		return res.Success(NewSetNode(configName, expr, startPos, expr.EndPos()))

	}

}

func (p *Parser) parseStatement() *ParseResult {

	res := NewParseResult()

	// TODO change to switch statement

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "set") {

		setRes := res.Register(p.parseSet())
		if res.Err != nil {
			return res
		}

		return res.Success(setRes)

	} else if p.CurrentToken.MatchType(lexer.TokenType_QuestionMark) {

		printRes := res.Register(p.parsePrintStdout())
		if res.Err != nil {
			return res
		}

		return res.Success(printRes)

	} else if p.CurrentToken.MatchType(lexer.TokenType_Comment) {

		token := p.CurrentToken
		res.RegisterAdvancement()
		p.advance()

		return res.Success(NewCommentNode(token))

	} else if p.CurrentToken.Match(lexer.TokenType_Keyword, "return") {

		printRes := res.Register(p.parseReturn())
		if res.Err != nil {
			return res
		}

		return res.Success(printRes)

	}

	panic(fmt.Sprintf("%v is not a valid statement", p.CurrentToken))
}

func (p *Parser) parseStatements() *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start.Copy()
	statements := make([]Node, 0, 10)

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

		if !moreStatements {
			break
		}

		if p.CurrentToken.Type == lexer.TokenType_EOF {
			break
		}

		statement := res.TryRegister(p.parseStatement())

		if statement == nil {
			p.reverseAmount(res.ToReverseCount)
			moreStatements = false
			continue
		}

		statements = append(statements, statement)

	}

	return res.Success(NewListNode(statements, startPos, &p.CurrentToken.Range.End))
}
