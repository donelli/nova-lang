package parser

import "nova-lang/pkg/shared"

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
	Node_Macro
)

//go:generate stringer -type=ParserNodeType -trimprefix=ParserNodeType_

type Node interface {
	StartPos() shared.Position
	EndPos() shared.Position
	ToHTML() string
	Type() ParserNodeType
	Range() *shared.Range
}
