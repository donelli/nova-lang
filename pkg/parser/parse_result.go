package parser

import "nova-lang/pkg/shared"

type ParseResult struct {
	Err                        *shared.Error
	Warnings                   []*shared.Warning
	Node                       Node
	LastRegisteredAdvanceCount int
	AdvanceCount               int
	ToReverseCount             int
}

func NewParseResult() *ParseResult {
	return &ParseResult{
		Err:                        nil,
		Warnings:                   make([]*shared.Warning, 0),
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

func (r *ParseResult) Register(res *ParseResult) Node {

	r.LastRegisteredAdvanceCount = res.AdvanceCount
	r.AdvanceCount += res.AdvanceCount

	if res.Err != nil {
		r.Err = res.Err
	}

	if len(res.Warnings) > 0 {
		r.Warnings = append(r.Warnings, res.Warnings...)
	}

	return res.Node
}

func (r *ParseResult) TryRegister(res *ParseResult) Node {

	if res.Err != nil {
		res.ToReverseCount = r.AdvanceCount
		return nil
	}

	return r.Register(res)
}

func (r *ParseResult) Success(node Node) *ParseResult {
	r.Node = node
	return r
}

func (r *ParseResult) Failure(err *shared.Error) *ParseResult {

	if r.Err == nil || r.LastRegisteredAdvanceCount == 0 {
		r.Err = err
	}

	return r
}

func (r *ParseResult) Warning(err *shared.Warning) *ParseResult {
	r.Warnings = append(r.Warnings, err)
	return r
}
