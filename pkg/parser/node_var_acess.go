package parser

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"recital_lsp/pkg/shared"
)

type VarAccessNode struct {
	VarName  string
	startPos *shared.Position
	endPos   *shared.Position
}

func NewVarAccessNode(lexerToken *lexer.LexerToken) *VarAccessNode {
	return &VarAccessNode{
		VarName:  lexerToken.Value,
		startPos: &lexerToken.Range.Start,
		endPos:   &lexerToken.Range.End,
	}
}

func (l *VarAccessNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *VarAccessNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *VarAccessNode) ToHTML() string {
	return BuildNodeBoxHTML("", "value-node", l.VarName)
}

func (l *VarAccessNode) String() string {
	return fmt.Sprintf("VarAssign{Var: %v, startPos: %v, endPos: %v}", l.VarName, l.startPos, l.endPos)
}
