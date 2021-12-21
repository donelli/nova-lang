package parser

import (
	"fmt"
	"recital_lsp/pkg/shared"
)

type CommandType int8

const (
	CommandType_Close CommandType = iota + 1
	CommandType_Clear
	CommandType_Exit
	CommandType_Loop
	CommandType_Dialog
	CommandType_Compile
	CommandType_Alias
	CommandType_Eject
	CommandType_Sleep
)

//go:generate stringer -type=CommandType -trimprefix=CommandType_

type CommandNode struct {
	CommandType CommandType
	Args        map[string]interface{}
	startPos    *shared.Position
	endPos      *shared.Position
}

func NewCommandNode(commandType CommandType, args map[string]interface{}, startPos *shared.Position, endPos *shared.Position) *CommandNode {
	return &CommandNode{
		CommandType: commandType,
		Args:        args,
		startPos:    startPos,
		endPos:      endPos,
	}
}

func NewCommandNodeRange(commandType CommandType, args map[string]interface{}, commandRange *shared.Range) *CommandNode {
	return &CommandNode{
		CommandType: commandType,
		Args:        args,
		startPos:    &commandRange.Start,
		endPos:      &commandRange.End,
	}
}

func (l *CommandNode) StartPos() *shared.Position {
	return l.startPos
}

func (l *CommandNode) EndPos() *shared.Position {
	return l.endPos
}

func (l *CommandNode) Type() ParserNodeType {
	return Node_Command
}

func (l *CommandNode) ToHTML() string {
	return BuildNodeBoxHTML("", "command-node", l.CommandType.String(), fmt.Sprintf("%v", l.Args))
}

func (l *CommandNode) String() string {
	return fmt.Sprintf("Command{Args: %v, startPos: %v, endPos: %v}", l.Args, l.startPos, l.endPos)
}
