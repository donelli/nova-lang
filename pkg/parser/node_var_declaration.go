package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type VarDeclarationNode struct {
	Modifier  string
	VarNames  []string
	nodeRange *shared.Range
}

func NewVarDeclarationNode(modifier string, varNames []string, startPos shared.Position, endPos shared.Position) *VarDeclarationNode {
	return &VarDeclarationNode{
		Modifier:  modifier,
		VarNames:  varNames,
		nodeRange: shared.NewRange(startPos, endPos),
	}
}

func (l *VarDeclarationNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *VarDeclarationNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *VarDeclarationNode) Type() ParserNodeType {
	return Node_VarDeclar
}

func (l *VarDeclarationNode) ToHTML() string {
	return BuildNodeBoxHTML(l.Modifier, "var-decl-node", l.VarNames...)
}

func (l *VarDeclarationNode) String() string {
	return fmt.Sprintf("VarDeclar{Modifier: %v, VarNames: %v, Range: %v}", l.Modifier, l.VarNames, l.nodeRange)
}

func (l *VarDeclarationNode) Range() *shared.Range {
	return l.nodeRange
}
