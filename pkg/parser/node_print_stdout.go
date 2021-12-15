package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type PrintStdoutNode struct {
	expr     Node
	startPos *shared.Position
	endPos   *shared.Position
}

func NewPrintStdoutNode(expr Node) *PrintStdoutNode {
	return &PrintStdoutNode{
		expr:     expr,
		startPos: expr.StartPos(),
		endPos:   expr.EndPos(),
	}
}

func (l *PrintStdoutNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *PrintStdoutNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *PrintStdoutNode) ToHTML() string {
	return BuildNodeBoxHTML("", "print-stdout-node", "print", l.expr.ToHTML())
}

func (l *PrintStdoutNode) String() string {
	return fmt.Sprintf("Print{Expr: %v, startPos: %v, endPos: %v}", l.expr, l.startPos, l.endPos)
}
