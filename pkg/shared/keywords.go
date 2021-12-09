package shared

import (
	"fmt"
	"recital_lsp/pkg/utils"
)

const keywordCount = 71

var KeywordsMap = make(map[string]string, keywordCount)

func generateKeyword(keyword string, repeat bool) {
	if repeat && len(keyword) > 4 {
		for i := 4; i <= len(keyword); i++ {
			KeywordsMap[keyword[:i]] = keyword
		}
	} else {
		KeywordsMap[keyword] = keyword
	}
}

func LoadKeywords() {

	if len(KeywordsMap) > 0 {
		return
	}

	// General
	generateKeyword("to", false)
	generateKeyword("do", false)

	// User interaction
	generateKeyword("set", false)
	generateKeyword("say", false)
	generateKeyword("get", false)
	generateKeyword("read", false)
	generateKeyword("picture", true)
	generateKeyword("valid", true)
	generateKeyword("prompt", true)
	generateKeyword("box", true)
	generateKeyword("clear", true)

	// Conditionals
	generateKeyword("if", false)
	generateKeyword("else", false)
	generateKeyword("elseif", false)
	generateKeyword("endif", false)
	generateKeyword("case", false)
	generateKeyword("otherwise", false)
	generateKeyword("or", false)
	generateKeyword("and", false)

	// Loops
	generateKeyword("while", false)
	generateKeyword("enddo", false)
	generateKeyword("for", false)
	generateKeyword("next", false)
	generateKeyword("exit", false)
	generateKeyword("loop", false)

	// Variables and functions
	generateKeyword("procedure", true)
	generateKeyword("function", true)
	generateKeyword("declare", true)
	generateKeyword("private", true)
	generateKeyword("public", true)
	generateKeyword("local", true)
	generateKeyword("parameters", true)
	generateKeyword("return", true)
	generateKeyword("additive", true)

	utils.Assert(len(KeywordsMap) == keywordCount, fmt.Sprintf("KeywordsMap size is %d, expected %d", len(KeywordsMap), keywordCount))

}
