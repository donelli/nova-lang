package lexer_test

import (
	"fmt"
	"recital_lsp/pkg/lexer"
	"testing"
)

func TestLexerSimpleStatement(t *testing.T) {

	lexerInstance := lexer.NewLexer("test.go", "lnVal = val(lcPreco) * 0.9")
	res, err := lexerInstance.Parse()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s", res)

}

func TestLexerMultipleLines(t *testing.T) {

	lexerInstance := lexer.NewLexer("test.go", `
	
	lnVal = 1

	@ 01, 01 get lnVal pict "999"
	read
	
	if lnVal = 7
		dial box "sete"
	else
		dialog box "outro"
	endif
	
	`)
	res, err := lexerInstance.Parse()

	if err != nil {
		t.Errorf("%s", err)
	}

	fmt.Printf("%s", res)

}
