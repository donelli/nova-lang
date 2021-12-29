package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type DoCaseCase struct {
	CaseExpr Node
	Body     Node
}

type CaseNode struct {
	Cases         []DoCaseCase
	OtherwiseCase Node
	nodeRange     *shared.Range
}

func NewDoCaseCase(condition Node, body Node) DoCaseCase {
	return DoCaseCase{
		CaseExpr: condition,
		Body:     body,
	}
}

func (l *DoCaseCase) String() string {
	return fmt.Sprintf("DoCaseCase{CaseExpr: %v, Body: %v}", l.CaseExpr, l.Body)
}

func NewCaseNode(DoCaseCases []DoCaseCase, otherwiseCase Node, startPos shared.Position, endPos shared.Position) *CaseNode {
	return &CaseNode{
		Cases:         DoCaseCases,
		OtherwiseCase: otherwiseCase,
		nodeRange:     shared.NewRange(startPos, endPos),
	}
}

func (l *CaseNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *CaseNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *CaseNode) Type() ParserNodeType {
	return Node_DoCase
}

func (l *CaseNode) ToHTML() string {

	content := ""

	for i := range l.Cases {
		content += `<div class="if-case">case:`
		content += l.Cases[i].CaseExpr.ToHTML()
		content += "<div>then:</div>"

		if l.Cases[i].Body == nil {
			content += "NO BODY"
		} else {
			content += l.Cases[i].Body.ToHTML()
		}

		content += `</div>`
	}

	if l.OtherwiseCase != nil {
		content += `<div class="if-case"><div>otherwise:</div>`
		content += l.OtherwiseCase.ToHTML()
		content += `</div>`
	}

	return BuildNodeBoxHTML("DO&nbsp;CASE", "if-node", content)
}

func (l *CaseNode) String() string {
	return fmt.Sprintf("CaseNode{Cases: %v, OtherwiseCase: %v, range: %v}", l.Cases, l.OtherwiseCase, l.nodeRange)
}

func (l *CaseNode) Range() *shared.Range {
	return l.nodeRange
}
