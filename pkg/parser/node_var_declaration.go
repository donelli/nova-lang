package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type VarDeclarationNode struct {
	Modifier string
	VarNames []string
	startPos *shared.Position
	endPos   *shared.Position
}

func NewVarDeclarationNode(modifier string, varNames []string, startPos *shared.Position, endPos *shared.Position) *VarDeclarationNode {
	return &VarDeclarationNode{
		Modifier: modifier,
		VarNames: varNames,
		startPos: startPos,
		endPos:   endPos,
	}
}

func (l *VarDeclarationNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *VarDeclarationNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *VarDeclarationNode) Type() ParserNodeType {
	return Node_VarDeclar
}

func (l *VarDeclarationNode) ToHTML() string {
	return BuildNodeBoxHTML(l.Modifier, "var-decl-node", l.VarNames...)
}

func (l *VarDeclarationNode) String() string {

	return fmt.Sprintf("VarDeclar{modifier: %v, varNames: %v, startPos: %v, endPos: %v}", l.Modifier, l.VarNames, l.startPos, l.endPos)

	// return fmt.Sprintf("VarDeclarationNode{Nodes: %v, startPos: %v, endPos: %v}", l.Nodes, l.startPos, l.endPos)
}
