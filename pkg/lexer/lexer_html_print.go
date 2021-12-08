package lexer

import (
	"html"
	"os"
)

func PrintLexerResultToHTML(lexerResult *LexerResult, outputFile string) {

	fh, err := os.Create(outputFile)

	if err != nil {
		panic(err)
	}

	defer fh.Close()

	fh.WriteString(`
	<html>
	<head>
	<style>
	.token {
		border: 1px solid gray;
		padding: 2px;
		margin-right: 2px;
		font-family: monospace;
	}
	.new-line {
		width: 100%;
		height: 10px;
		margin-top: 3px;
		margin-bottom: 3px;
	}
	.token-str {
		border-color: #00FF00;
		background-color: #00FF0033;
	}
	.token-ident {
		border-color: #4d91ff;
		background-color: #4d91ff33;
	}
	.token-keyword {
		border-color: #9200c7;
		background-color: #9200c733;
	}
	.token-number {
		border-color: #47ffbc;
		background-color: #47ffbc33;
	}
	.token-operator {
		
	}
	.token-comment {
		border-color: #006b0b;
		background-color: #006b0b33;
	}
	.token-bool {
		border-color: #e06900;
		background-color: #e069003b;
	}
	</style>
	</head>
	<body>
	`)

	for _, token := range lexerResult.Tokens {

		htmlStr := ""

		if token.Type == TokenType_NewLine {
			htmlStr = `<div class="new-line"></div>`
		} else if token.Type == TokenType_String {
			htmlStr = `<span class="token token-str">"` + html.EscapeString(token.Value) + `"</span>`
		} else if token.Type == TokenType_Identifier {
			htmlStr = `<span class="token token-ident">` + token.Value + `</span>`
		} else if token.Type == TokenType_Keyword {
			htmlStr = `<span class="token token-keyword">` + token.Value + `</span>`
		} else if token.Type == TokenType_Number {
			htmlStr = `<span class="token token-number">` + token.Value + `</span>`
		} else if token.Type == TokenType_Comment {
			htmlStr = `<span class="token token-comment">` + html.EscapeString(token.Value) + `</span>`
		} else if token.Type == TokenType_Date {
			htmlStr = `<span class="token token-number">{` + token.Value + `}</span>`
		} else if token.Type == TokenType_Boolean {
			htmlStr = `<span class="token token-bool">` + token.Value + `</span>`
		} else {
			htmlStr = `<span class="token token-operator">` + token.Value + `</span>`
		}

		fh.WriteString(htmlStr)
	}

	fh.WriteString("</body></html>")

}
