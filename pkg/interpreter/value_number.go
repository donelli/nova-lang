package interpreter

import (
	"fmt"
	"math"
	"nova-lang/pkg/shared"
)

type Number struct {
	Value float64
	Range *shared.Range
}

func (n *Number) Copy() Value {
	return &Number{
		Value: n.Value,
		Range: n.Range,
	}
}

func (n *Number) PrintRepresentation() string {

	intPart := int64(n.Value)
	decimalPart := n.Value - float64(intPart)

	if decimalPart > 0 {
		return fmt.Sprintf("%v", n.Value)
	} else {
		return fmt.Sprintf("%d", intPart)
	}

}

func (n *Number) Type() ValueType {
	return ValueType_Number
}

func (n *Number) Add(value Value) (Value, *shared.Error) {
	return NewNumber(n.Value + value.(*Number).Value), nil
}

func (n *Number) Subtract(value Value) (Value, *shared.Error) {
	return NewNumber(n.Value - value.(*Number).Value), nil
}

func (n *Number) Multiply(value Value) (Value, *shared.Error) {
	return NewNumber(n.Value * value.(*Number).Value), nil
}

func (n *Number) Divide(value Value) (Value, *shared.Error) {

	if value.(*Number).Value == 0 {
		return nil, shared.NewRuntimeErrorRange(n.Range, "Cannot divide by zero")
	}

	return NewNumber(n.Value / value.(*Number).Value), nil

}

func (n *Number) Exponential(value Value) (Value, *shared.Error) {
	return NewNumber(math.Pow(n.Value, value.(*Number).Value)), nil
}

func (n *Number) Remainder(value Value) (Value, *shared.Error) {
	return NewNumber(math.Mod(n.Value, value.(*Number).Value)), nil
}

func (n *Number) Equals(value Value) (Value, *shared.Error) {
	return NewBoolean(n.Value == value.(*Number).Value), nil
}

func (n *Number) NotEquals(value Value) (Value, *shared.Error) {
	return NewBoolean(n.Value != value.(*Number).Value), nil
}

func (n *Number) And(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(n.Range, "Numbers cannot be used as boolean/logic")
}

func (n *Number) Or(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(n.Range, "Numbers cannot be used as boolean/logic")
}

func (n *Number) IsGreater(value Value) (Value, *shared.Error) {
	return NewBoolean(n.Value > value.(*Number).Value), nil
}

func (n *Number) IsGreaterEquals(value Value) (Value, *shared.Error) {
	return NewBoolean(n.Value >= value.(*Number).Value), nil
}

func (n *Number) IsLess(value Value) (Value, *shared.Error) {
	return NewBoolean(n.Value < value.(*Number).Value), nil
}

func (n *Number) IsLessEquals(value Value) (Value, *shared.Error) {
	return NewBoolean(n.Value <= value.(*Number).Value), nil
}

func (n *Number) EqualsEquals(value Value) (Value, *shared.Error) {
	return n.Equals(value)
}

func (n *Number) IsTrue() (bool, *shared.Error) {
	return false, shared.NewRuntimeErrorRange(n.Range, "Numbers cannot be used as boolean/logic")
}

func (n *Number) UpdateRange(valueRange *shared.Range) Value {
	n.Range = valueRange
	return n
}

func NewNumber(value float64) *Number {
	return &Number{
		Value: value,
		Range: nil,
	}
}
