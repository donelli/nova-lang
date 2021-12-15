package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type VarAssignmentNode struct {
	VarName  string
	Expr     Node
	startPos *shared.Position
	endPos   *shared.Position
}

func NewVarAssignmentNode(varName string, expr Node, startPos *shared.Position, endPos *shared.Position) *VarAssignmentNode {
	return &VarAssignmentNode{
		VarName:  varName,
		Expr:     expr,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (l *VarAssignmentNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *VarAssignmentNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *VarAssignmentNode) ToHTML() string {
	return BuildNodeBoxHTML("", "var-assign-node", l.VarName, " = ", l.Expr.ToHTML())
}

func (l *VarAssignmentNode) String() string {
	return fmt.Sprintf("VarAssign{Var: %v, Expr: %v, startPos: %v, endPos: %v}", l.VarName, l.Expr, l.startPos, l.endPos)
}
