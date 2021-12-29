package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type VarAssignmentNode struct {
	VarName   string
	Expr      Node
	nodeRange *shared.Range
}

func NewVarAssignmentNode(varName string, expr Node, startPos shared.Position, endPos shared.Position) *VarAssignmentNode {
	return &VarAssignmentNode{
		VarName:   varName,
		Expr:      expr,
		nodeRange: shared.NewRange(startPos, endPos),
	}
}

func (l *VarAssignmentNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *VarAssignmentNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *VarAssignmentNode) Type() ParserNodeType {
	return Node_VarAssign
}

func (l *VarAssignmentNode) ToHTML() string {
	return BuildNodeBoxHTML("assign", "var-assign-node", l.VarName, l.Expr.ToHTML())
}

func (l *VarAssignmentNode) String() string {
	return fmt.Sprintf("VarAssign{Var: %v, Expr: %v, Range: %v}", l.VarName, l.Expr, l.nodeRange)
}

func (l *VarAssignmentNode) Range() *shared.Range {
	return l.nodeRange
}
