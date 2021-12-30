package parser

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/shared"
	"reflect"
	"strconv"
)

func TokenToString(tokenType lexer.LexerTokenType, value string) string {
	return fmt.Sprintf("%s%s", tokenType.String(), value)
}

var andTokenString = TokenToString(lexer.TokenType_Keyword, "and")
var orTokenString = TokenToString(lexer.TokenType_Keyword, "or")
var possibleClearArgs = []string{"all", "fcache", "gets", "iostats", "keys", "locks", "memory", "menus", "popups", "program", "prompt", "screen", "typeahead", "window"}
var ifKeywordsToIgnore = []string{"else", "elseif", "endif"}
var doCaseKeywordsToIgnore = []string{"case", "endcase", "otherwise"}

type ParseOption int8

const (
	optTerm ParseOption = iota + 1
	optFactor
	optCall
	optArithExpr
	optCompareExpr
	optOrExpr
)

type Parser struct {
	LexerResult       *lexer.LexerResult
	CurrentTokenIndex int
	CurrentToken      *lexer.LexerToken

	optToFunctionName map[ParseOption]reflect.Value
}

func NewParser(lexerResult *lexer.LexerResult) *Parser {
	parser := &Parser{
		LexerResult:       lexerResult,
		CurrentTokenIndex: -1,
		optToFunctionName: make(map[ParseOption]reflect.Value, 6),
	}
	parser.advance()

	parser.optToFunctionName[optTerm] = reflect.ValueOf(parser).MethodByName("ParseTerm")
	parser.optToFunctionName[optFactor] = reflect.ValueOf(parser).MethodByName("ParseFactor")
	parser.optToFunctionName[optCall] = reflect.ValueOf(parser).MethodByName("ParseCall")
	parser.optToFunctionName[optArithExpr] = reflect.ValueOf(parser).MethodByName("ParseArithmeticExpr")
	parser.optToFunctionName[optCompareExpr] = reflect.ValueOf(parser).MethodByName("ParseCompareExpr")
	parser.optToFunctionName[optOrExpr] = reflect.ValueOf(parser).MethodByName("ParseExpressionOr")

	return parser
}

func GetRangeFromNode(node Node) *shared.Range {
	return shared.NewRange(node.StartPos(), node.EndPos())
}

// func (p *Parser) reverse() {
// 	p.CurrentTokenIndex--
// 	p.updateCurrentToken()
// }

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

func (p *Parser) updateCurrentToken() {
	if p.CurrentTokenIndex >= 0 && p.CurrentTokenIndex < len(p.LexerResult.Tokens) {
		p.CurrentToken = p.LexerResult.Tokens[p.CurrentTokenIndex]
	}
}

func (p *Parser) Parse() *ParseResult {
	res := p.parseStatements([]string{})

	if res.Err == nil && p.CurrentToken.Type != lexer.TokenType_EOF {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Parser finished its work but there are still tokens left. Probably a error with the parser."))
	}

	return res
}

func (p *Parser) ParseArithmeticExpr() *ParseResult {
	return p.parseBinaryOperation(optTerm, optTerm, nil, []lexer.LexerTokenType{
		lexer.TokenType_Plus,
		lexer.TokenType_Minus,
	})
}

func (p *Parser) ParseTerm() *ParseResult {
	return p.parseBinaryOperation(optFactor, optFactor, nil, []lexer.LexerTokenType{
		lexer.TokenType_Star,
		lexer.TokenType_Slash,
		lexer.TokenType_Percent,
	})
}

func (p *Parser) ParseFactor() *ParseResult {

	res := NewParseResult()
	currentToken := p.CurrentToken

	if currentToken.MatchType(lexer.TokenType_Minus) {

		res.RegisterAdvancement()
		p.advance()

		factor := res.Register(p.ParseFactor())

		if res.Err != nil {
			return res
		}

		return res.Success(NewUnaryOperationNode(currentToken, factor))

	}

	if p.CurrentToken.MatchType(lexer.TokenType_Macro) {

		start := p.CurrentToken.Range.Start

		res.RegisterAdvancement()
		p.advance()

		factor := res.Register(p.ParseFactor())

		if res.Err != nil {
			return res
		}

		return res.Success(NewMacroNode(factor, start, factor.EndPos()))

	}

	return p.parsePower()
}

func (p *Parser) ParseCall() *ParseResult {

	res := NewParseResult()

	atom := res.Register(p.parseAtom())

	if res.Err != nil {
		return res
	}

	if p.CurrentToken.MatchType(lexer.TokenType_LeftParenthesis) {

		if atom.Type() != Node_VarAccess {
			return res.Failure(shared.NewInvalidSyntaxError(atom.StartPos(), atom.EndPos(), "Function name expected"))
		}

		endPos := p.CurrentToken.Range.End

		res.RegisterAdvancement()
		p.advance()
		argNodes := []Node{}

		if p.CurrentToken.MatchType(lexer.TokenType_RightParenthesis) {
			res.RegisterAdvancement()
			p.advance()
		} else {

			expr := res.Register(p.parseExpression())
			if res.Err != nil {
				return res
			}

			argNodes = append(argNodes, expr)

			for p.CurrentToken.MatchType(lexer.TokenType_Comma) {

				res.RegisterAdvancement()
				p.advance()

				expr = res.Register(p.parseExpression())
				if res.Err != nil {
					return res
				}

				argNodes = append(argNodes, expr)

			}

			if !p.CurrentToken.MatchType(lexer.TokenType_RightParenthesis) {
				return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected ',' or ')'"))
			}

			res.RegisterAdvancement()
			p.advance()

		}

		return res.Success(NewFunctionCallNode(atom, argNodes, atom.StartPos(), endPos))
	}

	// TODO Add here array acess `[`, `]`

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

}

func (p *Parser) parsePower() *ParseResult {

	return p.parseBinaryOperation(optCall, optFactor, nil, []lexer.LexerTokenType{
		lexer.TokenType_Exponential,
		lexer.TokenType_Pipe,
	})
}

func (p *Parser) ParseCompareExpr() *ParseResult {

	res := NewParseResult()

	if p.CurrentToken.MatchType(lexer.TokenType_Not) {

		operationToken := p.CurrentToken
		res.RegisterAdvancement()
		p.advance()

		node := res.Register(p.ParseCompareExpr())

		if res.Err != nil {
			return res
		}

		return res.Success(NewUnaryOperationNode(operationToken, node))

	}

	node := res.Register(p.parseBinaryOperation(optArithExpr, optArithExpr, []string{}, []lexer.LexerTokenType{
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

func (p *Parser) invokeFunction(funcName ParseOption) *ParseResult {

	// TODO test the performace of this method

	result := p.optToFunctionName[funcName].Call(nil)

	return result[0].Interface().(*ParseResult)

	// if funcName == optCompareExpr {
	// 	return p.parseCompareExpr()
	// } else if funcName == optArithExpr {
	// 	return p.parseArithmeticExpr()
	// } else if funcName == optTerm {
	// 	return p.parseTerm()
	// } else if funcName == optCall {
	// 	return p.parseCall()
	// } else if funcName == optFactor {
	// 	return p.parseFactor()
	// } else if funcName == optOrExpr {
	// 	return p.parseExpressionOr()
	// }

}

func (p *Parser) parseBinaryOperation(leftFuncName ParseOption, rightFuncName ParseOption, typeValueOptions []string, typeOptions []lexer.LexerTokenType) *ParseResult {

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

// I had to split the 'parseExpression' into two functions
// The 'AND' operator has more priority than the 'OR' operator

func (p *Parser) parseExpression() *ParseResult {

	res := NewParseResult()

	node := res.Register(p.parseBinaryOperation(optOrExpr, optOrExpr, []string{andTokenString}, []lexer.LexerTokenType{}))

	if res.Err != nil {
		return res
	}

	return res.Success(node)
}

func (p *Parser) ParseExpressionOr() *ParseResult {

	res := NewParseResult()

	node := res.Register(p.parseBinaryOperation(optCompareExpr, optCompareExpr, []string{orTokenString}, []lexer.LexerTokenType{}))

	if res.Err != nil {
		return res
	}

	return res.Success(node)
}
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

	return res.Success(NewCommandNodeRange(CommandType_Loop, nil, token.Range))
}

func (p *Parser) parseExit() *ParseResult {

	token := p.CurrentToken
	res := NewParseResult()

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected token after 'exit' keyword"))
	}

	return res.Success(NewCommandNodeRange(CommandType_Exit, nil, token.Range))
}

func (p *Parser) parseFunction() *ParseResult {

	res := NewParseResult()
	funcKeywordToken := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Identifier) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected function name"))
	}

	funcName := p.CurrentToken.Value

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected token after function name"))
	}

	node := res.Register(p.parseStatements([]string{"return", "function", "procedure"}))
	if res.Err != nil {
		return res
	}

	statements := node.(*ListNode)

	if !p.CurrentToken.Match(lexer.TokenType_Keyword, "return") {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'return' after function body (got "+p.CurrentToken.Value+")"))
	}

	returnToken := res.Register(p.parseReturn())
	if res.Err != nil {
		return res
	}

	statements.Nodes = append(statements.Nodes, returnToken)

	params := []string{}
	var paramNode Node = nil

	for i := range statements.Nodes {
		if statements.Nodes[i].Type() == Node_VarDeclar {

			if paramNode != nil {
				return res.Failure(shared.NewInvalidSyntaxError(statements.Nodes[i].StartPos(), statements.Nodes[i].EndPos(), "Multiple parameters definitions in function"))
			}

			paramNode = statements.Nodes[i]
			params = append(params, statements.Nodes[i].(*VarDeclarationNode).VarNames...)
		}
	}

	if paramNode != nil {
		if statements.Nodes[0] != paramNode {
			res.Warning(shared.NewWarning(paramNode.StartPos(), paramNode.EndPos(), "Function parameters should be defined in the first line of the function"))
		}
	}

	return res.Success(NewFunctionNode(funcName, statements, params, funcKeywordToken.Range.Start, returnToken.EndPos()))
}

func (p *Parser) parseForStatement() *ParseResult {

	res := NewParseResult()
	forKeywordToken := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Identifier) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected variable name after 'for' keyword"))
	}

	varName := p.CurrentToken.Value

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Equals) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected '=' after variable name in 'for' loop"))
	}

	res.RegisterAdvancement()
	p.advance()

	startExpr := res.Register(p.parseExpression())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.Match(lexer.TokenType_Keyword, "to") {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'to' after start value in 'for' loop"))
	}

	res.RegisterAdvancement()
	p.advance()

	endExpr := res.Register(p.parseExpression())
	if res.Err != nil {
		return res
	}

	var stepExpr Node = nil

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "step") {

		res.RegisterAdvancement()
		p.advance()

		stepExpr = res.Register(p.parseExpression())
		if res.Err != nil {
			return res
		}

	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected token after 'for' loop"))
	}

	body := res.Register(p.parseStatements([]string{"next"}))

	if p.CurrentToken.MatchType(lexer.TokenType_EOF) {
		return res.Failure(shared.NewInvalidSyntaxError(forKeywordToken.Range.Start, p.CurrentToken.Range.End, "Unclosed 'for' block"))
	}

	if !p.CurrentToken.Match(lexer.TokenType_Keyword, "next") {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'next' keyword"))
	}

	endPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	return res.Success(NewForNode(varName, startExpr, endExpr, stepExpr, body, forKeywordToken.Range.Start, endPos))
}

func (p *Parser) parseDoStatement() *ParseResult {

	res := NewParseResult()
	nextToken, hasNextToken := p.getNextToken()
	var node Node

	if !hasNextToken {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected something after 'do' keyword"))
	}

	if nextToken.Match(lexer.TokenType_Keyword, "case") {
		node = res.Register(p.parseDoCase())
	} else if nextToken.Match(lexer.TokenType_Keyword, "while") {
		node = res.Register(p.parseDoWhile())
	} else {
		node = res.Register(p.parseDoCommand())
	}

	if res.Err != nil {
		return res
	}

	return res.Success(node)
}

func (p *Parser) parseDoCommand() *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Skeleton) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected file/procedure name after 'do' keyword"))
	}

	args := map[string]interface{}{}

	args["procedure"] = p.CurrentToken.Value

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "with") {

		parameters := []Node{}

		res.RegisterAdvancement()
		p.advance()

		for {

			expr := res.Register(p.parseExpression())
			if res.Err != nil {
				return res
			}

			parameters = append(parameters, expr)

			if !p.CurrentToken.MatchType(lexer.TokenType_Comma) {
				break
			}

			res.RegisterAdvancement()
			p.advance()

		}

		args["parameters"] = parameters

	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected expression after 'do' command"))
	}

	return res.Success(NewCommandNode(CommandType_Do, args, startPos, p.CurrentToken.Range.End))
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

	return res.Success(NewDoWhileNode(condition, statements, startPos, endPos))
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
				res.Warning(shared.NewWarning(caseToken.Range.Start, condition.EndPos(), "Empty case body"))
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

	return res.Success(NewCaseNode(cases, otherwiseCase, startPos, endPos))
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

	return res.Success(NewIfNode(ifCases, elseCase, startPos, endifToken.Range.End))
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

	return res.Success(NewVarDeclarationNode(modifier, varNames, startPos, endPos))
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

	return res.Success(NewVarAssignmentNode(token.Value, expr, token.Range.Start, expr.EndPos()))
}

func (p *Parser) parseReturn() *ParseResult {

	res := NewParseResult()
	token := p.CurrentToken
	toMaster := false

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "to") {

		res.RegisterAdvancement()
		p.advance()

		if !p.CurrentToken.Match(lexer.TokenType_Keyword, "master") {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'master' keyword"))
		}

		toMaster = true

		res.RegisterAdvancement()
		p.advance()

	}

	if p.CurrentToken.Type == lexer.TokenType_NewLine || p.CurrentToken.Type == lexer.TokenType_EOF {
		return res.Success(NewReturnNode(nil, toMaster, token.Range.Start, token.Range.End))
	}

	expr := res.Register(p.parseExpression())

	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after return expression"))
	}

	return res.Success(NewReturnNode(expr, toMaster, token.Range.Start, expr.EndPos()))
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
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected expression at end of command"))
	}

	return res.Success(NewPrintStdoutNode(expr))
}

func (p *Parser) parseErase() *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start
	endPos := p.CurrentToken.Range.End

	res.RegisterAdvancement()
	p.advance()

	args := map[string]interface{}{}

	if p.CurrentToken.MatchType(lexer.TokenType_Skeleton) {
		args["file"] = p.CurrentToken.Value
		endPos = p.CurrentToken.Range.End

		res.RegisterAdvancement()
		p.advance()

	} else {

		expr := res.Register(p.parseExpression())
		if res.Err != nil {
			return res
		}

		args["expr"] = expr
		endPos = expr.EndPos()

	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected expression at end of command"))
	}

	return res.Success(NewCommandNode(CommandType_Erase, args, startPos, endPos))
}

func (p *Parser) parseCount() *ParseResult {

	res := NewParseResult()
	token := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Success(NewCommandNode(CommandType_Count, nil, token.Range.Start, token.Range.End))
	}

	args := map[string]interface{}{}

	if p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, []string{"all", "rest"}) {
		args["scope"] = p.CurrentToken.Value
		res.RegisterAdvancement()
		p.advance()
	}

	for {

		if p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, []string{"while", "for"}) {

			res.RegisterAdvancement()
			p.advance()

			whileExpr := res.Register(p.parseExpression())
			if res.Err != nil {
				return res
			}

			// forExpr / whileExpr
			args[p.CurrentToken.Value+"Expr"] = whileExpr
			continue

		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "title") {

			res.RegisterAdvancement()
			p.advance()

			titleExpr := res.Register(p.parseExpression())
			if res.Err != nil {
				return res
			}

			args["titleExpr"] = titleExpr
			continue

		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "to") {

			res.RegisterAdvancement()
			p.advance()

			if !p.CurrentToken.MatchType(lexer.TokenType_Identifier) {
				return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected variable name"))
			}

			args["varName"] = p.CurrentToken.Value

			res.RegisterAdvancement()
			p.advance()
			continue

		}

		break
	}

	return res.Success(NewCommandNode(CommandType_Count, args, token.Range.Start, p.CurrentToken.Range.Start))
}

func (p *Parser) parseBrowse() *ParseResult {

	res := NewParseResult()
	token := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Success(NewCommandNode(CommandType_Browse, nil, token.Range.Start, token.Range.End))
	}

	args := map[string]interface{}{}

	for {

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "for") {

			// TODO show error on multiple for's

			res.RegisterAdvancement()
			p.advance()

			expr := res.Register(p.parseExpression())
			if res.Err != nil {
				return res
			}

			args["forExpr"] = expr

			continue
		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "fields") {

			// TODO show error on multiple fields

			res.RegisterAdvancement()
			p.advance()

			fields := []string{}

			for {

				if !p.CurrentToken.MatchType(lexer.TokenType_Identifier) {
					return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected field name"))
				}

				fields = append(fields, p.CurrentToken.Value)

				res.RegisterAdvancement()
				p.advance()

				if p.CurrentToken.MatchType(lexer.TokenType_Comma) {
					res.RegisterAdvancement()
					p.advance()
					continue
				}

				break

			}

			args["fields"] = fields

			continue
		}

		break
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected expression at end of command"))
	}

	return res.Success(NewCommandNode(CommandType_Browse, args, token.Range.Start, p.CurrentToken.Range.Start))
}

func (p *Parser) parseEject() *ParseResult {

	res := NewParseResult()
	token := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after command"))
	}

	return res.Success(NewCommandNode(CommandType_Eject, nil, token.Range.Start, token.Range.End))
}

func (p *Parser) parseAlias() *ParseResult {

	res := NewParseResult()
	start := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Identifier) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected identifier (alias name)"))
	}

	aliasName := p.CurrentToken.Value

	res.RegisterAdvancement()
	p.advance()

	expr := res.Register(p.parseExpression())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after command"))
	}

	args := map[string]interface{}{"name": aliasName, "expr": expr}
	return res.Success(NewCommandNode(CommandType_Alias, args, start, expr.EndPos()))
}

func (p *Parser) parseSleep() *ParseResult {

	res := NewParseResult()
	start := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	expr := res.Register(p.parseExpression())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after command"))
	}

	args := map[string]interface{}{"expr": expr}
	return res.Success(NewCommandNode(CommandType_Sleep, args, start, expr.EndPos()))
}

func (p *Parser) parseCompile() *ParseResult {

	res := NewParseResult()
	start := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Skeleton) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected path to files (skeleton)"))
	}

	skeletonToken := p.CurrentToken

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after command"))
	}

	args := map[string]interface{}{"skeleton": skeletonToken.Value}
	return res.Success(NewCommandNode(CommandType_Compile, args, start, skeletonToken.Range.End))

}

func (p *Parser) parseDialog() *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.Type != lexer.TokenType_Keyword {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid dialog command. Expected keyword"))
	}

	if p.CurrentToken.Value == "box" {

		// DIALOG BOX <expC1> [LABEL <expC2>]

		res.RegisterAdvancement()
		p.advance()

		var label Node = nil

		expr := res.Register(p.parseExpression())
		if res.Err != nil {
			return res
		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "label") {
			res.RegisterAdvancement()
			p.advance()

			label = res.Register(p.parseExpression())
			if res.Err != nil {
				return res
			}

		}

		if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after command"))
		}

		endPos := expr.EndPos()
		if label != nil {
			endPos = label.EndPos()
		}

		args := map[string]interface{}{"subcommand": "box", "message": expr, "label": label}
		return res.Success(NewCommandNode(CommandType_Dialog, args, startPos, endPos))

	} else if p.CurrentToken.Value == "fields" {
		// DIALOG FIELDS [LABEL <expC>]
		panic("not implemented")
	} else if p.CurrentToken.Value == "files" {
		// DIALOG FILES LIKE <skeleton> [TRIM] [LABEL <expC>]
		panic("not implemented")
	} else if p.CurrentToken.Value == "get" {
		// DIALOG GET <memvar> [PICTURE <expC>] [HELP <expC>]
		// [LABEL <expC>] [TITLE <expC>]
		panic("not implemented")
	} else if p.CurrentToken.Value == "message" {

		// DIALOG MESSAGE <expC>

		res.RegisterAdvancement()
		p.advance()

		expr := res.Register(p.parseExpression())
		if res.Err != nil {
			return res
		}

		if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after `dialog message`"))
		}

		args := map[string]interface{}{"subcommand": "message", "message": expr}
		return res.Success(NewCommandNode(CommandType_Dialog, args, startPos, expr.EndPos()))

	} else if p.CurrentToken.Value == "query" {

		// DIALOG QUERY [LOCK]

		endPos := p.CurrentToken.Range.End
		lock := false

		res.RegisterAdvancement()
		p.advance()

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "lock") {
			lock = true
			endPos = p.CurrentToken.Range.End
			res.RegisterAdvancement()
			p.advance()
		}

		if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after `dialog query`"))
		}

		args := map[string]interface{}{"subcommand": "query", "lock": lock}
		return res.Success(NewCommandNode(CommandType_Dialog, args, startPos, endPos))

	} else if p.CurrentToken.Value == "scope" {

		// DIALOG SCOPE

		endPos := p.CurrentToken.Range.End

		res.RegisterAdvancement()
		p.advance()

		if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after `dialog scope`"))
		}

		args := map[string]interface{}{"subcommand": "scope"}
		return res.Success(NewCommandNode(CommandType_Dialog, args, startPos, endPos))

	} else {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid dialog command"))
	}

}

func (p *Parser) parseRelease() *ParseResult {

	// RELEASE <memvar list> / ALL [LIKE / EXCEPT <skeleton>]

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start
	endPos := p.CurrentToken.Range.End
	args := map[string]interface{}{}

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "all") {

		args["all"] = true
		endPos = p.CurrentToken.Range.End

		res.RegisterAdvancement()
		p.advance()

	} else {

		if p.CurrentToken.Type != lexer.TokenType_Identifier {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected variable names after `release`"))
		}

		varNames := []string{}
		varNames = append(varNames, p.CurrentToken.Value)
		endPos = p.CurrentToken.Range.End

		res.RegisterAdvancement()
		p.advance()

		for p.CurrentToken.MatchType(lexer.TokenType_Comma) {

			res.RegisterAdvancement()
			p.advance()

			if !p.CurrentToken.MatchType(lexer.TokenType_Identifier) {
				return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected variable name after `,`"))
			}

			varNames = append(varNames, p.CurrentToken.Value)
			endPos = p.CurrentToken.Range.End

			res.RegisterAdvancement()
			p.advance()

		}

		args["variables"] = varNames

	}

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "like") {

		res.RegisterAdvancement()
		p.advance()

		if !p.CurrentToken.MatchType(lexer.TokenType_Skeleton) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected pattern after `like`"))
		}

		args["like"] = p.CurrentToken.Value
		endPos = p.CurrentToken.Range.End

		res.RegisterAdvancement()
		p.advance()

	} else if p.CurrentToken.Match(lexer.TokenType_Keyword, "except") {

		res.RegisterAdvancement()
		p.advance()

		if !p.CurrentToken.MatchType(lexer.TokenType_Skeleton) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected pattern after `like`"))
		}

		args["except"] = p.CurrentToken.Value
		endPos = p.CurrentToken.Range.End

		res.RegisterAdvancement()
		p.advance()

	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after command"))
	}

	return res.Success(NewCommandNode(CommandType_Release, args, startPos, endPos))
}

func (p *Parser) parseStore() *ParseResult {

	// STORE <exp> TO <memvar>

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	expr := res.Register(p.parseExpression())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.Match(lexer.TokenType_Keyword, "to") {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected `to` after expression"))
	}

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.Type != lexer.TokenType_Identifier {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected variable name after `to`"))
	}

	varNames := []string{}
	varNames = append(varNames, p.CurrentToken.Value)

	res.RegisterAdvancement()
	p.advance()

	for p.CurrentToken.MatchType(lexer.TokenType_Comma) {

		res.RegisterAdvancement()
		p.advance()

		if !p.CurrentToken.MatchType(lexer.TokenType_Identifier) {
			return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected variable name after `,`"))
		}

		varNames = append(varNames, p.CurrentToken.Value)

		res.RegisterAdvancement()
		p.advance()

	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid token after `store` variables"))
	}

	args := map[string]interface{}{"varNames": varNames, "value": expr}
	return res.Success(NewCommandNode(CommandType_Store, args, startPos, p.CurrentToken.Range.End))
}

func (p *Parser) parseClose() *ParseResult {

	// CLOSE [<alias> / ALL / DATABASES / FORMAT / INDEX / PROCEDURE / ALTERNATE [TO PRINT]]

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Success(NewCommandNode(CommandType_Close, nil, startPos, p.CurrentToken.Range.End))
	}

	args := map[string]interface{}{}

	if p.CurrentToken.MatchType(lexer.TokenType_Identifier) {

		args["alias"] = p.CurrentToken.Value

		res.RegisterAdvancement()
		p.advance()

	} else if p.CurrentToken.MatchType(lexer.TokenType_Keyword) {

		args["arg"] = p.CurrentToken.Value

		res.RegisterAdvancement()
		p.advance()

		if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {

			if args["arg"] == "alternate" {

				if !p.CurrentToken.Match(lexer.TokenType_Keyword, "to") {
					return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'TO' after ALTERNATE"))
				}

				res.RegisterAdvancement()
				p.advance()

				if !p.CurrentToken.Match(lexer.TokenType_Keyword, "print") {
					return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Expected 'PRINT' after TO"))
				}

				res.RegisterAdvancement()
				p.advance()

				args["toPrint"] = true

			} else {
				return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unrecognized token after close argument"))
			}

		}

	} else if p.CurrentToken.MatchType(lexer.TokenType_Number) {

		args["workarea"] = p.CurrentToken.Value

		res.RegisterAdvancement()
		p.advance()

	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unrecognized token after close argument"))
	}

	return res.Success(NewCommandNode(CommandType_Close, args, startPos, p.CurrentToken.Range.End))
}

func (p *Parser) parseClear() *ParseResult {

	// clear [all/fcache/gets/iostats/keys/locks/memory/menus/popups/program/prompt/screen/typeahead/window]

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Success(NewCommandNode(CommandType_Clear, nil, startPos, p.CurrentToken.Range.End))
	}

	if !p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, possibleClearArgs) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Invalid clear argument: `"+p.CurrentToken.Value+"`"))
	}

	arg := p.CurrentToken.Value

	res.RegisterAdvancement()
	p.advance()

	args := map[string]interface{}{arg: true}

	return res.Success(NewCommandNode(CommandType_Clear, args, startPos, p.CurrentToken.Range.End))
}

func (p *Parser) parseSet() *ParseResult {

	// Types of set:
	// - set <keyword> to <value> ?
	// - set <keyword> <value> ?

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start

	res.RegisterAdvancement()
	p.advance()

	if !p.CurrentToken.MatchType(lexer.TokenType_Keyword) {
		return res.Failure(shared.NewInvalidSyntaxError(startPos, p.CurrentToken.Range.End, "Expected valid configuration name (keyword)"))
	}

	configName := p.CurrentToken.Value
	res.RegisterAdvancement()
	p.advance()

	if p.CurrentToken.Match(lexer.TokenType_Keyword, "to") {

		res.RegisterAdvancement()
		p.advance()

		if p.CurrentToken.MatchType(lexer.TokenType_NewLine) || p.CurrentToken.MatchType(lexer.TokenType_Comment) {
			fmt.Println("set to empty")
			return res.Success(NewEmptySetNode(configName, startPos, p.CurrentToken.Range.End))
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
			return res.Success(NewEmptySetNode(configName, startPos, p.CurrentToken.Range.End))
		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "on") {
			fmt.Println("set On")

			endPos := p.CurrentToken.Range.End
			res.RegisterAdvancement()
			p.advance()

			return res.Success(NewBoolSetNode(configName, "on", startPos, endPos))
		}

		if p.CurrentToken.Match(lexer.TokenType_Keyword, "off") {
			fmt.Println("set off")

			endPos := p.CurrentToken.Range.End
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

	if p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, keywordsToIgnore) {
		return res
	}

	if p.CurrentToken.MatchType(lexer.TokenType_QuestionMark) {

		successNode = res.Register(p.parsePrintStdout())

	} else if p.CurrentToken.MatchType(lexer.TokenType_Comment) {

		token := p.CurrentToken
		res.RegisterAdvancement()
		p.advance()

		successNode = NewCommentNode(token)

	} else if p.CurrentToken.MatchType(lexer.TokenType_Identifier) {

		nextToken, hasNextToken := p.getNextToken()

		if hasNextToken && nextToken.MatchType(lexer.TokenType_Equals) {
			successNode = res.Register(p.parseVariableAssignment())
		}

	} else if p.CurrentToken.MatchType(lexer.TokenType_Keyword) {

		// TODO maybe we could use a map(string -> func), to speed up the parsing of commands

		if p.CurrentToken.Value == "if" {

			successNode = res.Register(p.parseIfStatement())

		} else if p.CurrentToken.Value == "do" {

			successNode = res.Register(p.parseDoStatement())

		} else if p.CurrentToken.Value == "exit" {

			successNode = res.Register(p.parseExit())

		} else if p.CurrentToken.Value == "loop" {

			successNode = res.Register(p.parseLoop())

		} else if p.CurrentToken.Value == "for" {

			successNode = res.Register(p.parseForStatement())

		} else if p.CurrentToken.Value == "return" {

			successNode = res.Register(p.parseReturn())

		} else if p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, []string{"function", "procedure"}) {

			successNode = res.Register(p.parseFunction())

		} else if p.CurrentToken.MatchMultiple(lexer.TokenType_Keyword, []string{"private", "public", "local", "parameters"}) {

			successNode = res.Register(p.parseVariableDeclaration())

		} else if p.CurrentToken.Value == "set" {

			successNode = res.Register(p.parseSet())

		} else if p.CurrentToken.Value == "clear" {

			successNode = res.Register(p.parseClear())

		} else if p.CurrentToken.Value == "store" {

			successNode = res.Register(p.parseStore())

		} else if p.CurrentToken.Value == "release" {

			successNode = res.Register(p.parseRelease())

		} else if p.CurrentToken.Value == "close" {

			successNode = res.Register(p.parseClose())

		} else if p.CurrentToken.Value == "dialog" {

			successNode = res.Register(p.parseDialog())

		} else if p.CurrentToken.Value == "sleep" {

			successNode = res.Register(p.parseSleep())

		} else if p.CurrentToken.Value == "compile" {

			successNode = res.Register(p.parseCompile())

		} else if p.CurrentToken.Value == "alias" {

			successNode = res.Register(p.parseAlias())

		} else if p.CurrentToken.Value == "eject" {

			successNode = res.Register(p.parseEject())

		} else if p.CurrentToken.Value == "browse" {

			successNode = res.Register(p.parseBrowse())

		} else if p.CurrentToken.Value == "count" {

			successNode = res.Register(p.parseCount())

		} else if p.CurrentToken.Value == "erase" {

			successNode = res.Register(p.parseErase())

		}

	}

	if res.Err != nil {
		return res
	}

	if successNode != nil {
		return res.Success(successNode)
	}

	expr := res.Register(p.parseExpression())
	if res.Err != nil {
		return res
	}

	if !p.CurrentToken.MatchType(lexer.TokenType_NewLine) {
		return res.Failure(shared.NewInvalidSyntaxErrorRange(p.CurrentToken.Range, "Unexpected token"))
	}

	return res.Success(expr)
}

func (p *Parser) parseStatements(keywordsToIgnore []string) *ParseResult {

	res := NewParseResult()
	startPos := p.CurrentToken.Range.Start
	statements := make([]Node, 0, 10)

	for p.CurrentToken.Type == lexer.TokenType_NewLine {
		res.RegisterAdvancement()
		p.advance()
	}

	statement := res.Register(p.parseStatement(keywordsToIgnore))

	if statement == nil {
		return res.Success(NewListNode(statements, startPos, p.CurrentToken.Range.End))
	}

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

		statement := res.Register(p.parseStatement(keywordsToIgnore))

		if res.Err != nil {
			return res
		}

		if statement == nil {
			break
		}

		statements = append(statements, statement)

	}

	return res.Success(NewListNode(statements, startPos, p.CurrentToken.Range.End))
}
