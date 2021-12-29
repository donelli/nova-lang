package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type PrintStdoutNode struct {
	expr      Node
	nodeRange *shared.Range
}

func NewPrintStdoutNode(expr Node) *PrintStdoutNode {
	return &PrintStdoutNode{
		expr:      expr,
		nodeRange: expr.Range(),
	}
}

func (l *PrintStdoutNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *PrintStdoutNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *PrintStdoutNode) Type() ParserNodeType {
	return Node_PrintStdout
}

func (l *PrintStdoutNode) ToHTML() string {
	return BuildNodeBoxHTML("print", "print-stdout-node", l.expr.ToHTML())
}

func (l *PrintStdoutNode) String() string {
	return fmt.Sprintf("Print{Expr: %v, Range: %v}", l.expr, l.nodeRange)
}

func (l *PrintStdoutNode) Range() *shared.Range {
	return l.nodeRange
}
