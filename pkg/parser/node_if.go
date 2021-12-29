package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type IfCase struct {
	CaseExpr Node
	Body     Node
}

type IfNode struct {
	IfCases   []IfCase
	ElseCase  Node
	nodeRange *shared.Range
}

func NewIfCase(condition Node, body Node) IfCase {
	return IfCase{
		CaseExpr: condition,
		Body:     body,
	}
}

func (l *IfCase) String() string {
	return fmt.Sprintf("IfCase{CaseExpr: %v, Body: %v}", l.CaseExpr, l.Body)
}

func NewIfNode(ifCases []IfCase, elseCase Node, startPos shared.Position, endPos shared.Position) *IfNode {
	return &IfNode{
		IfCases:   ifCases,
		ElseCase:  elseCase,
		nodeRange: shared.NewRange(startPos, endPos),
	}
}

func (l *IfNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *IfNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *IfNode) Type() ParserNodeType {
	return Node_If
}

func (l *IfNode) ToHTML() string {

	content := ""

	for i := range l.IfCases {
		content += `<div class="if-case">case:`
		content += l.IfCases[i].CaseExpr.ToHTML()
		content += "<div>then:</div>"
		content += l.IfCases[i].Body.ToHTML()
		content += `</div>`
	}

	if l.ElseCase != nil {
		content += `<div class="if-case"><div>else:</div>`
		content += l.ElseCase.ToHTML()
		content += `</div>`
	}

	return BuildNodeBoxHTML("IF", "if-node", content)
}

func (l *IfNode) String() string {
	return fmt.Sprintf("If{Cases: %v, Else: %v, Range: %v}", l.IfCases, l.ElseCase, l.nodeRange)
}

func (l *IfNode) Range() *shared.Range {
	return l.nodeRange
}
