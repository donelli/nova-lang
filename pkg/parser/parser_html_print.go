package parser

import "os"

func PrintParseResultToHTML(parseRes *ParseResult, outputFile string) {

	fh, err := os.Create(outputFile)

	if err != nil {
		panic(err)
	}

	defer fh.Close()

	fh.WriteString(`
	<html>
	<head>
	<style>
	.node {
		border: 1px solid lightgray;
		padding: 5px;
	}
	.node-comment {
		border-color: #006b0b;
		background-color: #006b0b33;
	}
	.node h3 {
		margin-top: 0;
		margin-bottom: 3px;
		border-bottom: 1px solid lightgray;
		width: auto;
	}
	.node {
		display: inline-block;
	}
	</style>
	</head>
	<body>
	`)

	fh.WriteString(parseRes.Node.ToHTML())

	fh.WriteString("</body></html>")
}
