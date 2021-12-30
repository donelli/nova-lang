package parser

import (
	"fmt"
	"nova-lang/pkg/shared"
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
	CommandType_Store
	CommandType_Release
	CommandType_Browse
	CommandType_Count
	CommandType_Do
	CommandType_Erase
)

//go:generate stringer -type=CommandType -trimprefix=CommandType_

type CommandNode struct {
	CommandType CommandType
	Args        map[string]interface{}
	nodeRange   *shared.Range
}

func NewCommandNode(commandType CommandType, args map[string]interface{}, startPos shared.Position, endPos shared.Position) *CommandNode {
	return &CommandNode{
		CommandType: commandType,
		Args:        args,
		nodeRange:   shared.NewRange(startPos, endPos),
	}
}

func NewCommandNodeRange(commandType CommandType, args map[string]interface{}, commandRange *shared.Range) *CommandNode {
	return &CommandNode{
		CommandType: commandType,
		Args:        args,
		nodeRange:   commandRange,
	}
}

func (l *CommandNode) StartPos() shared.Position {
	return l.nodeRange.Start
}

func (l *CommandNode) EndPos() shared.Position {
	return l.nodeRange.End
}

func (l *CommandNode) Type() ParserNodeType {
	return Node_Command
}

func (l *CommandNode) ToHTML() string {
	return BuildNodeBoxHTML("", "command-node", l.CommandType.String(), fmt.Sprintf("%v", l.Args))
}

func (l *CommandNode) String() string {
	return fmt.Sprintf("Command{Args: %v, Range: %v}", l.Args, l.nodeRange)
}

func (l *CommandNode) Range() *shared.Range {
	return l.nodeRange
}
