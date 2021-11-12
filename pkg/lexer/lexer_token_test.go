package lexer

import (
	"recital_lsp/pkg/shared"
	"testing"
)

func TestStringToken(t *testing.T) {
	tok := NewStringToken(shared.Range{
		Start: shared.Position{
			Row:    0,
			Column: 0,
		},
		End: shared.Position{
			Row:    0,
			Column: 5,
		},
	}, "\"test\"")

	strExp := `str:"test"[(0,0)-(0,5)]`

	if tok.String() != strExp {
		t.Fatalf("Expected %s got %s", strExp, tok.String())
	}

}
