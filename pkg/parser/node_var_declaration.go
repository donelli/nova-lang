package parser

import (
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

func (l *VarDeclarationNode) ToHTML() string {

	panic("not implemented")

	// str := "<div class=\"node node-list\">"

	// for i := range l.Nodes {
	// 	str += l.Nodes[i].ToHTML() + "<hr>"
	// }

	// return str + "</div>"
}

func (l *VarDeclarationNode) String() string {

	panic("not implemented")

	// return fmt.Sprintf("VarDeclarationNode{Nodes: %v, startPos: %v, endPos: %v}", l.Nodes, l.startPos, l.endPos)
}