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
	IfCases  []IfCase
	ElseCase Node
	startPos *shared.Position
	endPos   *shared.Position
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

func NewIfNode(ifCases []IfCase, elseCase Node, startPos *shared.Position, endPos *shared.Position) *IfNode {
	return &IfNode{
		IfCases:  ifCases,
		ElseCase: elseCase,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (l *IfNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *IfNode) EndPos() *shared.Position {
	return l.endPos
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
	panic("not implemented")
}
