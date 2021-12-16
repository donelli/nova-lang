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
var orTokenString = TokenToString(lexer.TokenType_Keyword, "or")

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

func GetRangeFromNode(node Node) *shared.Range {
	return shared.NewRange(*node.StartPos(), *node.EndPos())
}

func (p *Parser) advance() {
	p.CurrentTokenIndex++
	p.updateCurrentToken()
}

func (p *Parser) getNextToken() (*lexer.LexerToken, bool) {

	index := p.CurrentTokenIndex + 1

	if index < len(p.LexerResult.Tokens) {
		return p.LexerResult.Tokens[index], true
	}

	return nil, false
}

// func (p *Parser) reverseAmount(amount int) {
// 	p.CurrentTokenIndex -= amount
// 	p.updateCurrentToken()
// }

func (p *Parser) updateCurrentToken() {
	if p.CurrentTokenIndex >= 0 && p.CurrentTokenIndex < len(p.LexerResult.Tokens) {
		p.CurrentToken = p.LexerResult.Tokens[p.CurrentTokenIndex]
	}
}

func (p *Parser) Parse() *ParseResult {
	res := p.parseStatements([]string{})

	if res.Err == nil && p.CurrentToken.Type != lexer.TokenType_EOF {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Token cannot appear after previous tokens"))
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

		res.RegisterAdvancement()
		p.advance()

		return res.Success(NewStringNode(token))

	} else if token.MatchType(lexer.TokenType_Boolean) {

		res.RegisterAdvancement()
		p.advance()

		return res.Success(NewBooleanNode(token.Value == ".t.", token))

	} else if token.MatchType(lexer.TokenType_Identifier) {

		res.RegisterAdvancement()
		p.advance()

		return res.Success(NewVarAccessNode(token))

	} else if token.MatchType(lexer.TokenType_LeftParenthesis) {

		res.RegisterAdvancement()
		p.advance()

		expr := res.Register(p.parseExpression())

		if res.Err != nil {
			return res
		}

		if !p.CurrentToken.MatchType(lexer.TokenType_RightParenthesis) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected ')' after expression"))
		}

		res.RegisterAdvancement()
		p.advance()

		return res.Success(expr)
	}

	return res.Failure(shared.NewInvalidSyntaxError(token.Range.Start, token.Range.End, fmt.Sprintf("Expected number, string, bool, identifier or parenthesis, found %s", token.Type.String())))

	// return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected expression, found end of line"))
	// return res.Failure(shared.NewInvalidSyntaxError(token.Range.Start, token.Range.End, "Unexpected token"))
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
		lexer.TokenType_EqualsEquals,
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

	panic(fmt.Sprintf("'%s' is not a valid function", funcName))
}

func (p *Parser) parseBinaryOperation(leftFuncName string, rightFuncName string, typeValueOptions []string, typeOptions []lexer.LexerTokenType) *ParseResult {

	res := NewParseResult()

	leftRes := res.Register(p.invokeFunction(leftFuncName))

	if res.Err != nil {
		return res
	}

	for {

		isValidOption := false

		if len(typeValueOptions) > 0 {

			tokenStr := TokenToString(p.CurrentToken.Type, p.CurrentToken.Value)

			for opt := range typeValueOptions {
				if tokenStr == typeValueOptions[opt] {
					isValidOption = true
					break
				}
			}

		} else {

			for opt := range typeOptions {
				if typeOptions[opt] == p.CurrentToken.Type {
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

var ifKeywordsToIgnore = []string{"else", "elseif", "endif"}
var doCaseKeywordsToIgnore = []string{"case", "endcase", "otherwise"}

func (p *Parser) parseIfCase(caseWord string) (*ParseResult, []IfCase, Node) {

	res := NewParseResult()
	cases := []IfCase{}
	var elseCase Node = nil

	startIfPos := p.CurrentToken.Range.Start

	if !p.CurrentToken.Match(lexer.TokenType_Keyword, caseWord) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, fmt.Sprintf("Expected '%s' keyword", caseWord))), nil, nil
	}

	res.RegisterAdvancement()
	p.advance()

	condition := res.Register(p.parseExpression())
	if res.Err != nil {
		return res, nil, nil
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected new line after expression")), nil, nil
	}

	statements := res.Register(p.parseStatements(ifKeywordsToIgnore))

	if res.Err != nil {
		return res, nil, nil
	}

	cases = append(cases, NewIfCase(condition, statements))

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "endif") {
		res.RegisterAdvancement()
		p.advance()
		return res, cases, elseCase
	}

	if !p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, []string{"else", "elseif"}) {
		return res.Failure(shared.NewInvalidSyntaxError(startIfPos, p.CurrentToken.Range.End, "If block unclosed. Expected 'else', 'elseif' or 'endif' keyword")), nil, nil
	}

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "elseif") {

		elseifErr, elseifCases, elseNode := p.parseIfCase("elseif")

		if elseifErr.Err != nil {
			return elseifErr, nil, nil
		}

		elseCase = elseNode

		if len(elseifCases) > 0 {
			cases = append(cases, elseifCases...)
		}

	} else {

		res.RegisterAdvancement()
		p.advance()

		if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "'else' cannot have conditions")), nil, nil
		}

		startPos := p.CurrentToken.Range.Start

		statements := res.Register(p.parseStatements(ifKeywordsToIgnore))
		if res.Err != nil {
			return res, nil, nil
		}

		elseCase = statements

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "else") {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Only one 'else' can be used by 'if'")), nil, nil
		}

		if !p.CurrentToken.Match(lexer.TokenType_Keyword, "endif") {
			return res.Failure(shared.NewInvalidSyntaxError(startPos, p.CurrentToken.Range.End, "Expected 'endif' keyword")), nil, nil
		}

	}

	return res, cases, elseCase
}

func (p *Parser) parseLoop() *ParseResult {

	token := p.CurrentToken
	res := NewParseResult()

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected token after 'loop' keyword"))
	}

	return res.Success(NewLoopNode(token.Range))
}

func (p *Parser) parseExit() *ParseResult {

	token := p.CurrentToken
	res := NewParseResult()

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected token after 'exit' keyword"))
	}

	return res.Success(NewExitNode(token.Range))
}

func (p *Parser) parseDoStatement() *ParseResult {

	res := NewParseResult()
	nextToken, hasNextToken := p.getNextToken()
	var node Node

	if !hasNextToken {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected something after 'do' keyword"))
	}

	if !nextToken.MatchMultiple(lexer.TokenType_Keyword, []string{"while", "case"}) {
		res.RegisterAdvancement()
		p.advance()
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'case' or 'while' after 'do' keyword"))
	}

	if nextToken.Match(lexer.TokenType_Keyword, "case") {

		node = res.Register(p.parseDoCase())
		if res.Err != nil {
			return res
		}

	} else {

		// While

		node = res.Register(p.parseDoWhile())
		if res.Err != nil {
			return res
		}

	}

	return res.Success(node)
}

func (p *Parser) parseDoWhile() *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	res.RegisterAdvancement()
	p.advance()

	condition := res.Register(p.parseExpression())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected token after 'do while' condition"))
	}

	statements := res.Register(p.parseStatements([]string{"enddo"}))

	if p.CurrentToken.MatchType(lexer.TokenType_EOF) {
		return res.Failure(shared.NewInvalidSyntaxError(startPos, p.CurrentToken.Range.End, "Unclosed 'do while' block"))
	}

	if !p.CurrentToken.Match(lexer.TokenType_Keyword, "enddo") {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'enddo' keyword"))
	}

	endPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	return res.Success(NewDoWhileNode(condition, statements, &startPos, &endPos))
}

func (p *Parser) parseDoCase() *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.Match(lexer.TokenType_Keyword, "case") {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'case' after 'do' keyword"))
	}

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected token after 'do case' keyword"))
	}

	res.RegisterAdvancement()
	p.advance()

	var cases []DoCaseCase = []DoCaseCase{}
	var otherwiseCase Node = nil

	for {

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "endcase") {
			break
		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "otherwise") {

			res.RegisterAdvancement()
			p.advance()

			if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
				return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected expression after 'otherwise'"))
			}

			otherwiseCase = res.Register(p.parseStatements([]string{"endcase"}))
			continue

		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "case") {

			caseToken := p.CurrentToken

			res.RegisterAdvancement()
			p.advance()

			if p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
				return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected conditional expression after 'case' keyword"))
			}

			condition := res.Register(p.parseExpression())

			if res.Err != nil {
				return res
			}

			if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
				return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected expression after 'do case' condition expression"))
			}

			res.RegisterAdvancement()
			p.advance()

			statements := res.Register(p.parseStatements(doCaseKeywordsToIgnore))

			if res.Err != nil {
				return res
			}

			if statements == nil {
				res.Warning(shared.NewWarning(caseToken.Range.Start, *condition.EndPos(), "Empty case body"))
			}

			doCase := NewDoCaseCase(condition, statements)
			cases = append(cases, doCase)

			continue
		}

		if p.CurrentToken.MatchType(lexer.TokenType_EOF) {
			return res.Failure(shared.NewInvalidSyntaxError(startPos, p.CurrentToken.Range.End, "Unclosed 'do case' block"))
		}

		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'case', 'otherwise' or 'endcase' keyword, found '"+p.CurrentToken.Value+"'"))
	}

	if !p.CurrentToken.Match(lexer.TokenType_Keyword, "endcase") {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'endcase' keyword"))
	}

	endPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	return res.Success(NewCaseNode(cases, otherwiseCase, &startPos, &endPos))
}

func (p *Parser) parseIfStatement() *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	ifRes, ifCases, elseCase := p.parseIfCase("if")

	res.Register(ifRes)
	if res.Err != nil {
		return res
	}

	endifToken := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	return res.Success(NewIfNode(ifCases, elseCase, &startPos, &endifToken.Range.End))
}

func (p *Parser) parseVariableDeclaration() *ParseResult {

	res := NewParseResult()
	modifier := p.CurrentToken.Value
	varNames := []string{}
	startPos := p.CurrentToken.Range.Start
	endPos := p.CurrentToken.Range.End

	res.RegisterAdvancement()
	p.advance()

	for !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {

		if !p.CurrentToken.MatchType(lexer.TokenType_Identifier) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected variable name"))
		}

		varNames = append(varNames, p.CurrentToken.Value)
		endPos = p.CurrentToken.Range.End

		res.RegisterAdvancement()
		p.advance()

		if p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
			break
		}

		if !p.CurrentToken.MatchType(lexer.TokenType_Comma) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected ',' after variable name"))
		}

		res.RegisterAdvancement()
		p.advance()

	}

	return res.Success(NewVarDeclarationNode(modifier, varNames, &startPos, &endPos))
}

func (p *Parser) parseVariableAssignment() *ParseResult {

	res := NewParseResult()
	token := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Equals) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected '=' after variable name"))
	}

	res.RegisterAdvancement()
	p.advance()

	expr := res.Register(p.parseExpression())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after variable assignment"))
	}

	return res.Success(NewVarAssignmentNode(token.Value, expr, &token.Range.Start, expr.EndPos()))
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

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after return expression"))
	}

	return res.Success(NewReturnNode(expr, &token.Range.Start, expr.EndPos()))
}

func (p *Parser) parsePrintStdout() *ParseResult {

	// ? <expr>

	res := NewParseResult()

	res.RegisterAdvancement()
	p.advance()

	expr := res.Register(p.parseExpression())

	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected end of command"))
	}

	return res.Success(NewPrintStdoutNode(expr))
}

func (p *Parser) parseSet() *ParseResult {

	// Types of set:
	// - set <keyword> to <value> ?
	// - set <keyword> <value> ?

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start.Copy()

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Keyword) {
		return res.Failure(shared.NewInvalidSyntaxError(*startPos, p.CurrentToken.Range.End, "Expected valid configuration name (keyword)"))
	}

	configName := p.CurrentToken.Value
	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "to") {

		res.RegisterAdvancement()
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
			res.RegisterAdvancement()
			p.advance()

			return res.Success(NewBoolSetNode(configName, "on", startPos, endPos))
		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "off") {
			fmt.Println("set off")

			endPos := p.CurrentToken.Range.End.Copy()
			res.RegisterAdvancement()
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

func (p *Parser) parseStatement(keywordsToIgnore []string) *ParseResult {

	res := NewParseResult()
	var successNode Node = nil

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "set") {

		setRes := res.Register(p.parseSet())
		if res.Err != nil {
			return res
		}

		successNode = setRes

	} else if p.CurrentToken.MatchType(lexer.TokenType_QuestionMark) {

		printRes := res.Register(p.parsePrintStdout())
		if res.Err != nil {
			return res
		}

		successNode = printRes

	} else if p.CurrentToken.MatchType(lexer.TokenType_Comment) {

		token := p.CurrentToken
		res.RegisterAdvancement()
		p.advance()

		successNode = NewCommentNode(token)

	} else if p.CurrentToken.Match(lexer.TokenType_Keyword, "return") {

		returnRes := res.Register(p.parseReturn())
		if res.Err != nil {
			return res
		}

		successNode = returnRes

	} else if p.CurrentToken.MatchType(lexer.TokenType_Identifier) {

		varAssignRes := res.Register(p.parseVariableAssignment())
		if res.Err != nil {
			return res
		}

		successNode = varAssignRes

	} else if p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, []string{"private", "public", "local"}) {

		varDeclarNode := res.Register(p.parseVariableDeclaration())
		if res.Err != nil {
			return res
		}

		successNode = varDeclarNode

	} else if p.CurrentToken.Match(lexer.TokenType_Keyword, "if") {

		ifRes := res.Register(p.parseIfStatement())
		if res.Err != nil {
			return res
		}

		successNode = ifRes

	} else if p.CurrentToken.Match(lexer.TokenType_Keyword, "do") {

		doRes := res.Register(p.parseDoStatement())
		if res.Err != nil {
			return res
		}

		successNode = doRes

	} else if p.CurrentToken.Match(lexer.TokenType_Keyword, "exit") {

		doRes := res.Register(p.parseExit())
		if res.Err != nil {
			return res
		}

		successNode = doRes

	} else if p.CurrentToken.Match(lexer.TokenType_Keyword, "loop") {

		doRes := res.Register(p.parseLoop())
		if res.Err != nil {
			return res
		}

		successNode = doRes

	}

	if successNode != nil {
		return res.Success(successNode)
	}

	if p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, keywordsToIgnore) {
		return res
	}

	return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid statement for start of line: '"+p.CurrentToken.Value+"'"))
}

func (p *Parser) parseStatements(keywordsToIgnore []string) *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start.Copy()
	statements := make([]Node, 0, 10)

	for p.CurrentToken.Type == lexer.TokenType_NewLine {
		res.RegisterAdvancement()
		p.advance()
	}

	statement := res.Register(p.parseStatement(keywordsToIgnore))

	if statement == nil || res.Err != nil {
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

		statement := res.Register(p.parseStatement(keywordsToIgnore))

		if res.Err != nil {
			return res
		}

		if statement == nil {
			break
		}

		// if statement == nil {

		// 	fmt.Println("reverse: ", len(statements))

		// 	p.reverseAmount(res.ToReverseCount)
		// 	moreStatements = false
		// 	continue
		// }

		statements = append(statements, statement)

	}

	return res.Success(NewListNode(statements, startPos, &p.CurrentToken.Range.End))
}
