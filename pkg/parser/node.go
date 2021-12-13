package parser

import "recital_lsp/pkg/shared"

type Node interface {
	StartPos() *shared.Position
	EndPos() *shared.Position
}
