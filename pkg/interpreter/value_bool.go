package interpreter

import (
	"nova-lang/pkg/shared"
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
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot add boolean/logic values")
}

func (b *Boolean) Subtract(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot subtract boolean/logic values")
}

func (b *Boolean) Multiply(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot multiply boolean/logic values")
}

func (b *Boolean) Divide(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot divide boolean/logic values")
}

func (b *Boolean) Exponential(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform exponentialization with boolean/logic values")
}

func (b *Boolean) Remainder(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform remainer operation with boolean/logic values")
}

func (b *Boolean) Equals(value Value) (Value, *shared.Error) {
	return NewBoolean(b.Value == value.(*Boolean).Value), nil
}

func (b *Boolean) NotEquals(value Value) (Value, *shared.Error) {
	return NewBoolean(b.Value != value.(*Boolean).Value), nil
}

func (b *Boolean) And(value Value) (Value, *shared.Error) {
	return NewBoolean(b.Value && value.(*Boolean).Value), nil
}

func (b *Boolean) Or(value Value) (Value, *shared.Error) {
	return NewBoolean(b.Value || value.(*Boolean).Value), nil
}

func (b *Boolean) EqualsEquals(value Value) (Value, *shared.Error) {
	return b.Equals(value)
}

func (b *Boolean) IsGreater(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform greater operation with boolean/logic values")
}

func (b *Boolean) IsGreaterEquals(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform greater/equals operation with boolean/logic values")
}

func (b *Boolean) IsLess(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform less operation with boolean/logic values")
}

func (b *Boolean) IsLessEquals(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(b.Range, "Cannot perform less/equals operation with boolean/logic values")
}

func (b *Boolean) IsTrue() (bool, *shared.Error) {
	return b.Value, nil
}

func (b *Boolean) IsEmpty() bool {
	return b.Value == false
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
