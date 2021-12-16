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
	startPos      *shared.Position
	endPos        *shared.Position
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

func NewCaseNode(DoCaseCases []DoCaseCase, otherwiseCase Node, startPos *shared.Position, endPos *shared.Position) *CaseNode {
	return &CaseNode{
		Cases:         DoCaseCases,
		OtherwiseCase: otherwiseCase,
		startPos:      startPos,
		endPos:        endPos,
	}
}

func (l *CaseNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *CaseNode) EndPos() *shared.Position {
	return l.endPos
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
	panic("not implemented")
}
