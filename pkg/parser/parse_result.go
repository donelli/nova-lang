package parser

import "recital_lsp/pkg/shared"

type ParseResult struct {
	Err                        *shared.Error
	Node                       *Node
	LastRegisteredAdvanceCount int
	AdvanceCount               int
	ToReverseCount             int
}

func NewParseResult() *ParseResult {
	return &ParseResult{
		Err:                        nil,
		Node:                       nil,
		LastRegisteredAdvanceCount: 0,
		AdvanceCount:               0,
		ToReverseCount:             0,
	}
}

func (r *ParseResult) RegisterAdvancement() {
	r.LastRegisteredAdvanceCount = 1
	r.AdvanceCount++
}

func (r *ParseResult) Register(res *ParseResult) *Node {

	r.LastRegisteredAdvanceCount = res.AdvanceCount
	r.AdvanceCount += res.AdvanceCount

	if res.Err != nil {
		r.Err = res.Err
	}

	return res.Node
}

func (r *ParseResult) TryRegister(res *ParseResult) *Node {

	if res.Err != nil {
		r.Err = res.Err
		return nil
	}

	return r.Register(res)
}

func (r *ParseResult) Success(node Node) *ParseResult {
	r.Node = &node
	return r
}

func (r *ParseResult) Failure(err *shared.Error) *ParseResult {
	r.Err = err
	return r
}
