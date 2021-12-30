package interpreter

import "nova-lang/pkg/shared"

type RuntimeResult struct {
	Value               Value
	Error               *shared.Error
	FunctionReturnValue Value
	LoopShouldLoop      bool
	LoopShouldExit      bool
}

func (runtimeResult *RuntimeResult) reset() {
	runtimeResult.Error = nil
	runtimeResult.Value = nil
	runtimeResult.FunctionReturnValue = nil
	runtimeResult.LoopShouldLoop = false
	runtimeResult.LoopShouldExit = false
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

func (runtimeResult *RuntimeResult) SuccessExit() *RuntimeResult {
	runtimeResult.reset()
	runtimeResult.LoopShouldExit = true
	return runtimeResult
}

func (runtimeResult *RuntimeResult) SuccessLoop() *RuntimeResult {
	runtimeResult.reset()
	runtimeResult.LoopShouldLoop = true
	return runtimeResult
}

func (res *RuntimeResult) Register(otherRes *RuntimeResult) Value {

	res.Error = otherRes.Error
	res.FunctionReturnValue = otherRes.FunctionReturnValue
	res.LoopShouldLoop = otherRes.LoopShouldLoop
	res.LoopShouldExit = otherRes.LoopShouldExit

	return otherRes.Value

}

func (res *RuntimeResult) ShouldReturn() bool {
	return res.Error != nil || res.FunctionReturnValue != nil || res.LoopShouldExit || res.LoopShouldLoop
}

func NewRuntimeResult() *RuntimeResult {
	return &RuntimeResult{
		Value:               nil,
		Error:               nil,
		FunctionReturnValue: nil,
		LoopShouldLoop:      false,
		LoopShouldExit:      false,
	}
}
