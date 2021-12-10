package parser

import "recital_lsp/pkg/shared"

// TODO change this to use "polymorphism": https://golangbot.com/polymorphism/

type Node struct {
	nodeRange *shared.Range
}

func NewNode(startPos *shared.Position, endPos *shared.Position) *Node {
	return &Node{
		nodeRange: shared.NewRange(*startPos, *endPos),
	}
}
