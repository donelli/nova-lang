package parser

import "recital_lsp/pkg/shared"

// import "recital_lsp/pkg/shared"

// TODO change this to use "polymorphism": https://golangbot.com/polymorphism/

type Node interface {
	StartPos() *shared.Position
	EndPos() *shared.Position
}
