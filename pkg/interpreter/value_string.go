package interpreter

import (
	"fmt"
	"nova-lang/pkg/shared"
	"strings"
)

// TODO: use string builder to optimize strings
// https://medium.com/swlh/high-performance-string-building-in-go-golang-3fd99b9ca856

type String struct {
	Value string
	Range *shared.Range
}

func (n *String) Copy() Value {
	return &String{
		Value: n.Value,
		Range: n.Range,
	}
}

func (n *String) PrintRepresentation() string {
	return n.Value
}

func (n *String) Type() ValueType {
	return ValueType_String
}

func (n *String) Add(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_String {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot add a `%v` to a string", value.Type()))
	}

	return NewString(n.Value + value.(*String).Value), nil

}

func (n *String) Subtract(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_String {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot subtract a `%v` from a string", value.Type()))
	}

	return NewString(fmt.Sprintf("%s%s", strings.TrimRight(n.Value, " "), value.(*String).Value)), nil

}

func (n *String) Multiply(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(n.Range, "Cannot multiply a string")
}

func (n *String) Divide(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(n.Range, "Cannot divide a string")
}

func (n *String) Exponential(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(n.Range, "Cannot perform exponentialization on a string")
}

func (n *String) Remainder(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(n.Range, "Cannot perform remainer operation on a string")
}

func (n *String) Equals(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_String {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot compare `%v` and `%v`", n.Type(), value.Type()))
	}

	rightString := value.(*String)

	if len(rightString.Value) > len(n.Value) {
		return NewBoolean(false), nil
	}

	compareLen := len(rightString.Value)

	return NewBoolean(n.Value[0:compareLen-1] == rightString.Value[0:compareLen-1]), nil
}

func (n *String) NotEquals(value Value) (Value, *shared.Error) {

	isEquals, err := n.Equals(value)
	if err != nil {
		return nil, err
	}

	return NewBoolean(!isEquals.(*Boolean).Value), nil
}

func (n *String) And(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot perform 'AND' on: `%v` and `%v`", n.Type(), value.Type()))
}

func (n *String) Or(value Value) (Value, *shared.Error) {
	return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot perform 'AND' on: `%v` and `%v`", n.Type(), value.Type()))
}

func (n *String) IsGreater(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_String {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot compare `%v` with a string", value.Type()))
	}

	return NewBoolean(n.Value > value.(*String).Value), nil
}

func (n *String) IsGreaterEquals(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_String {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot compare `%v` with a string", value.Type()))
	}

	return NewBoolean(n.Value >= value.(*String).Value), nil
}

func (n *String) IsLess(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_String {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot compare `%v` with a string", value.Type()))
	}

	return NewBoolean(n.Value < value.(*String).Value), nil
}

func (n *String) IsLessEquals(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_String {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot compare `%v` with a string", value.Type()))
	}

	return NewBoolean(n.Value <= value.(*String).Value), nil
}

func (n *String) EqualsEquals(value Value) (Value, *shared.Error) {

	if value.Type() != ValueType_String {
		return nil, shared.NewRuntimeErrorRange(n.Range, fmt.Sprintf("Cannot compare `%v` and `%v`", n.Type(), value.Type()))
	}

	rightString := value.(*String)

	return NewBoolean(n.Value == rightString.Value), nil
}

func (n *String) IsTrue() (bool, *shared.Error) {
	return false, shared.NewRuntimeErrorRange(n.Range, "Strings cannot be used as boolean")
}

func (n *String) UpdateRange(valueRange *shared.Range) Value {
	n.Range = valueRange
	return n
}

func NewString(value string) *String {
	return &String{
		Value: value,
		Range: nil,
	}
}
