package parser

import (
	"fmt"
	"nova-lang/pkg/lexer"
	"nova-lang/pkg/shared"
)

type VarAccessNode struct {
	VarName   string
	nodeRange *shared.Range
}

func NewVarAccessNode(lexerToken *lexer.LexerToken) *VarAccessNode {
	return &VarAccessNode{
		VarName:   lexerToken.Value,
		nodeRange: lexerToken.Range,
	}
}

func (l *VarAccessNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *VarAccessNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *VarAccessNode) Type() ParserNodeType {
	return Node_VarAccess
}

func (l *VarAccessNode) ToHTML() string {
	return BuildNodeBoxHTML("var", "value-node", l.VarName)
}

func (l *VarAccessNode) String() string {
	return fmt.Sprintf("VarAccess{Var: %v, Range: %v}", l.VarName, l.nodeRange)
}

func (l *VarAccessNode) Range() *shared.Range {
	return l.nodeRange
}
