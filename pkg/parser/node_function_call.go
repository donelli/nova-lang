package parser

import (
	"fmt"
	"nova-lang/pkg/shared"
)

type FunctionCallNode struct {
	FunctionName Node
	Args         []Node
	nodeRange    *shared.Range
}

func NewFunctionCallNode(funcName Node, args []Node, startPos shared.Position, endPos shared.Position) *FunctionCallNode {
	return &FunctionCallNode{
		FunctionName: funcName,
		Args:         args,
		nodeRange:    shared.NewRange(startPos, endPos),
	}
}

func (l *FunctionCallNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *FunctionCallNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *FunctionCallNode) Type() ParserNodeType {
	return Node_FuncCall
}

func (l *FunctionCallNode) ToHTML() string {

	args := ""
	for i := range l.Args {
		args += "<div>" + l.Args[i].ToHTML() + "</div>"
	}

	return BuildNodeBoxHTML("CALL", "func-call-node",
		"<div>Func:<br>"+l.FunctionName.ToHTML()+"</div>"+
			"<div>Args:<br>"+args+"</div>")
}

func (l *FunctionCallNode) String() string {
	argsStr := ""

	for i := range l.Args {
		if i < len(l.Args)-1 {
			argsStr += ", "
		}
		argsStr += fmt.Sprintf("%s", l.Args[i])
	}

	return fmt.Sprintf("Call{func: %v, args: %s, Range: %v}", l.FunctionName, argsStr, l.nodeRange)
}

func (l *FunctionCallNode) Range() *shared.Range {
	return l.nodeRange
}
