package interpreter

import (
	"recital_lsp/pkg/shared"
)

type Boolean struct {
	Value bool
	Range *shared.Range
}

func (b *Boolean) Copy() Value {
	return &Boolean{
		Value: b.Value,
		Range: b.Range,
	}
}

func (b *Boolean) PrintRepresentation() string {

	if b.Value {
		return ".t."
	} else {
		return ".f."
	}

}

func (b *Boolean) Type() ValueType {
	return ValueType_Boolean
}

func (b *Boolean) Add(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot add a boolean/logic value")
}

func (b *Boolean) Subtract(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot subtract a boolean/logic value")
}

func (b *Boolean) Multiply(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot multiply a boolean/logic value")
}

func (b *Boolean) Divide(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot divide a boolean/logic value")
}

func (b *Boolean) Exponential(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform exponentialization with a boolean/logic value")
}

func (b *Boolean) Remainder(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform remainer operation with a boolean/logic value")
}

func (b *Boolean) UpdateRange(valueRange *shared.Range) Value {
	b.Range = valueRange
	return b
}

func NewBoolean(value bool, BooleanRange *shared.Range) *Boolean {
	return &Boolean{
		Value: value,
		Range: BooleanRange,
	}
}
