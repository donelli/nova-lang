package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type FunctionNode struct {
	FuncName   string
	Body       Node
	Parameters []string
	nodeRange  *shared.Range
}

func NewFunctionNode(funcName string, body Node, parameters []string, startPos shared.Position, endPos shared.Position) *FunctionNode {
	return &FunctionNode{
		FuncName:   funcName,
		Body:       body,
		Parameters: parameters,
		nodeRange:  shared.NewRange(startPos, endPos),
	}
}

func (l *FunctionNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *FunctionNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *FunctionNode) Type() ParserNodeType {
	return Node_Function
}

func (l *FunctionNode) ToHTML() string {

	paramsStr := ""
	for index, param := range l.Parameters {

		if index > 0 {
			paramsStr += ", "
		}

		paramsStr += param
	}

	body := fmt.Sprintf("<div>Name:<br>%s</div><div>Params:<br>%s</div><div>Body:<br>%s</div>", l.FuncName, paramsStr, l.Body.ToHTML())
	return BuildNodeBoxHTML("FUNCTION", "function-node", body)
}

func (l *FunctionNode) String() string {
	return fmt.Sprintf("Func{name: %v, params: %v, body %v, Range: %v}", l.FuncName, l.Parameters, l.Body, l.nodeRange)
}

func (l *FunctionNode) Range() *shared.Range {
	return l.nodeRange
}
