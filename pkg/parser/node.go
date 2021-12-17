package parser

import "recital_lsp/pkg/shared"

type ParserNodeType uint8

const (
	Node_BinOp ParserNodeType = iota + 1
	Node_Bool
	Node_Comment
	Node_DoCase
	Node_DoWhile
	Node_ForLoop
	Node_Function
	Node_If
	Node_List
	Node_Number
	Node_PrintStdout
	Node_Return
	Node_Set
	Node_String
	Node_UnaryOp
	Node_VarAccess
	Node_VarAssign
	Node_VarDeclar
	Node_FuncCall
	Node_Command
)

type Node interface {
	StartPos() *shared.Position
	EndPos() *shared.Position
	ToHTML() string
	Type() ParserNodeType
}
