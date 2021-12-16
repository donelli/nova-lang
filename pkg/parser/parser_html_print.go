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
		padding: 3px;
		display: inline-block;
		position: relative;
	}

	.node-box {
		border: 1px solid lightgray;
		padding: 5px;
		display: flex;
	}
	.node-box > .node-box-type {
		width: 10px;
		word-wrap: break-word;
		font-family: monospace;
		font-size: 12px;
		border-right: 1px solid lightgray;
		margin-right: 5px;
		padding-right: 3px;
		text-align: center;
	}
	.node-box > .node-box-subnodes {
		display: flex;
		align-items: center;
	}
	.node-box > .node-box-subnodes > span {
		height: 100%;
		display: flex;
		align-items: center;
		justify-content: center;
	}
	.node-box .node-box-subnodes > span:not(:last-child) {
		border-right: 1px solid lightgray;
		margin-right: 6px;
		padding-right: 6px;
	}
	
	.comment-node {
		border-color: #006b0b;
		background-color: #006b0b33;
	}
	.comment-node.node-box > .node-box-type {
		border-color: #006b0b;
	}

	.print-stdout-node {
		border-color: #00A5CC;
		background-color: #00A5CC33;
	}
	.print-stdout-node.node-box > .node-box-type {
		border-color: #00A5CC;
	}
	.print-stdout-node.node-box > .node-box-subnodes > span {
		border-color: #00A5CC;
	}

	.bin-op-node {
		border-color: #CC1F26;
		background-color: #CC1F2633;
	}
	.bin-op-node.node-box > .node-box-type {
		border-color: #CC1F26;
	}
	.bin-op-node.node-box > .node-box-subnodes > span {
		border-color: #CC1F26;
	}

	.var-assign-node {
		border-color: #FFA500;
		background-color: #FFA50033;
	}
	.var-assign-node.node-box > .node-box-type {
		border-color: #FFA500;
	}
	.var-assign-node.node-box > .node-box-subnodes > span {
		border-color: #FFA500;
	}

	.value-node {
		border-color: #9732a8;
		padding: 3px;
	}
	.value-node.node-box > .node-box-type {
		border-color: #9732a8;
	}

	.if-node > .node-box-subnodes > span {
		flex-direction: column;
	}

	.if-case {
		border-bottom: 1px solid lightgray;
		margin-bottom: 5px;
		padding: 3px;
		width: 100%;
	}

	</style>
	</head>
	<body>
	`)

	fh.WriteString(parseRes.Node.ToHTML())

	fh.WriteString("</body></html>")
}

func buildSubNodesHTML(nodesStr []string) string {
	str := ""

	for i := range nodesStr {
		str += "<span>" + nodesStr[i] + "</span>"
	}

	return str
}

func BuildNodeBoxHTML(name string, className string, nodesStr ...string) string {

	if name != "" {
		name = `<div class="node-box-type">` + name + `</div>`
	}

	return `<div class="node-box ` + className + `">
		` + name + `
		<div class="node-box-subnodes">
			` + buildSubNodesHTML(nodesStr) + `
		</div>
	</div>`
}
