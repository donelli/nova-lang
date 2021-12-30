package interpreter

import (
	"fmt"
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
		return ".T."
	} else {
		return ".F."
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

func (b *Boolean) Equals(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Boolean {
		return nil, shared.NewRuntimeErrorRange(b.Range, fmt.Sprintf("Cannot compare `%v` and `%v`", b.Type(), value.Type()))
	}

	return NewBoolean(b.Value == value.(*Boolean).Value), nil
}

func (b *Boolean) NotEquals(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Boolean {
		return nil, shared.NewRuntimeErrorRange(b.Range, fmt.Sprintf("Cannot compare `%v` and `%v`", b.Type(), value.Type()))
	}

	return NewBoolean(b.Value != value.(*Boolean).Value), nil
}

func (b *Boolean) And(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Boolean {
		return nil, shared.NewRuntimeErrorRange(b.Range, fmt.Sprintf("Cannot compare `%v` and `%v`", b.Type(), value.Type()))
	}

	return NewBoolean(b.Value && value.(*Boolean).Value), nil
}

func (b *Boolean) Or(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Boolean {
		return nil, shared.NewRuntimeErrorRange(b.Range, fmt.Sprintf("Cannot compare `%v` and `%v`", b.Type(), value.Type()))
	}

	return NewBoolean(b.Value || value.(*Boolean).Value), nil
}

func (b *Boolean) EqualsEquals(value Value) (Value, *shared.Error) {
	return b.Equals(value)
}

func (b *Boolean) IsGreater(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform `>` operation with a boolean/logic value")
}

func (b *Boolean) IsGreaterEquals(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform `>=` operation with a boolean/logic value")
}

func (b *Boolean) IsLess(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform `<` operation with a boolean/logic value")
}

func (b *Boolean) IsLessEquals(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform `<=` operation with a boolean/logic value")
}

func (b *Boolean) UpdateRange(valueRange *shared.Range) Value {
	b.Range = valueRange
	return b
}

func NewBoolean(value bool) *Boolean {
	return &Boolean{
		Value: value,
		Range: nil,
	}
}
