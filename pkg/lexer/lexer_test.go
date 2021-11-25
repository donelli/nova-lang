package lexer_test

import (
	"recital_lsp/pkg/lexer"
	"testing"
)

func TestNewLexer(t *testing.T) {

	lexerInstance := lexer.NewLexer("test.go", "test")

	if lexerInstance.FileName != "test.go" {
		t.Errorf("FileName should be 'test.go', got %s", lexerInstance.FileName)
		return
	}

	if lexerInstance.CurrentPosition.String() != "(0,-1)" {
		t.Errorf("CurrentPosition should be '0,-1', got %s", lexerInstance.CurrentPosition.String())
		return
	}

	lexerInstance.Advance()

	if lexerInstance.CurrentPosition.String() != "(0,0)" {
		t.Errorf("CurrentPosition should be '0,0', got %s", lexerInstance.CurrentPosition.String())
		return
	}

	if lexerInstance.CurrentChar != 't' {
		t.Errorf("CurrentChar should be 't', got %s", string(lexerInstance.CurrentChar))
		return
	}

}

func TestLexerSimpleStatement(t *testing.T) {

	lexerInstance := lexer.NewLexer("test.go", "lnVal = val(lcPreco) * 0.9")
	_, err := lexerInstance.Parse()

	if err != nil {
		t.Errorf("%s", err)
	}

	// fmt.Printf("%s", res)

}

/*
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

}*/
