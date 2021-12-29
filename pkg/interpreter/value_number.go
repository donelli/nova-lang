package interpreter

import (
	"fmt"
	"math"
	"recital_lsp/pkg/shared"
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

	if value.Type() != ValueType_Number {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot add a `%v` to a number", value.Type()))
	}

	return NewNumber(n.Value+value.(*Number).Value, n.Range), nil

}

func (n *Number) Subtract(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Number {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot subtract a `%v` from a number", value.Type()))
	}

	return NewNumber(n.Value-value.(*Number).Value, n.Range), nil

}

func (n *Number) Multiply(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Number {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot multiply a `%v` with a number", value.Type()))
	}

	return NewNumber(n.Value*value.(*Number).Value, n.Range), nil

}

func (n *Number) Divide(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Number {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot divide a `%v` with a number", value.Type()))
	}

	if value.(*Number).Value == 0 {
		return nil, shared.NewRuntimeErrorRange(n.Range, "Cannot divide by zero")
	}

	return NewNumber(n.Value/value.(*Number).Value, n.Range), nil

}

func (n *Number) Exponential(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Number {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot perform exponentialization with `%v` and a number", value.Type()))
	}

	return NewNumber(math.Pow(n.Value, value.(*Number).Value), n.Range), nil

}

func (n *Number) Remainder(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_Number {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot perform remainer operation with `%v` and a number", value.Type()))
	}

	return NewNumber(math.Mod(n.Value, value.(*Number).Value), n.Range), nil

}

func (n *Number) UpdateRange(valueRange *shared.Range) Value {
	n.Range = valueRange
	return n
}

func NewNumber(value float64, numberRange *shared.Range) *Number {
	return &Number{
		Value: value,
		Range: numberRange,
	}
}
