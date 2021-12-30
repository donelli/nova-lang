package interpreter

import "nova-lang/pkg/shared"

type RuntimeResult struct {
	Value               Value
	Error               *shared.Error
	FunctionReturnValue Value
	LoopShouldContinue  bool
	LoopShouldBreak     bool
}

func (runtimeResult *RuntimeResult) reset() {
	runtimeResult.Error = nil
	runtimeResult.Value = nil
	runtimeResult.FunctionReturnValue = nil
	runtimeResult.LoopShouldContinue = false
	runtimeResult.LoopShouldBreak = false
}

func (runtimeResult *RuntimeResult) Success(value Value) *RuntimeResult {
	runtimeResult.reset()
	runtimeResult.Value = value
	return runtimeResult
}

func (runtimeResult *RuntimeResult) Failure(err *shared.Error) *RuntimeResult {
	runtimeResult.reset()
	runtimeResult.Error = err
	return runtimeResult
}

func (runtimeResult *RuntimeResult) SuccessReturn(value Value) *RuntimeResult {
	runtimeResult.reset()
	runtimeResult.FunctionReturnValue = value
	return runtimeResult
}

func (res *RuntimeResult) Register(otherRes *RuntimeResult) Value {

	res.Error = otherRes.Error
	res.FunctionReturnValue = otherRes.FunctionReturnValue
	res.LoopShouldContinue = otherRes.LoopShouldContinue
	res.LoopShouldBreak = otherRes.LoopShouldBreak

	return otherRes.Value

}

func (res *RuntimeResult) ShouldReturn() bool {

	// TODO check if this logic is correct
	// Note: this will allow you to continue and break outside the current function

	return res.Error != nil || res.FunctionReturnValue != nil || res.LoopShouldContinue || res.LoopShouldBreak
}

func NewRuntimeResult() *RuntimeResult {
	return &RuntimeResult{
		Value:               nil,
		Error:               nil,
		FunctionReturnValue: nil,
		LoopShouldContinue:  false,
		LoopShouldBreak:     false,
	}
}
