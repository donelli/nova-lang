package interpreter

import (
	"fmt"
	"nova-lang/pkg/shared"
)

// TODO: use string builder to optimize strings
// https://medium.com/swlh/high-performance-string-building-in-go-golang-3fd99b9ca856

// Strings in Nova are arrays of bytes because of the special modifications that can be aplied to show them (Reverse, Blink, ...)
type String struct {
	Value []rune
	Range *shared.Range
}

func (n *String) Copy() Value {
	return &String{
		Value: n.Value,
		Range: n.Range,
	}
}

func (n *String) PrintRepresentation() string {
	return string(n.Value)
}

func (n *String) Type() ValueType {
	return ValueType_String
}

func (n *String) Add(value Value) (Value, *shared.Error) {

	otherStr := value.(*String)
	thisLen := len(n.Value)

	newStr := make([]rune, len(n.Value)+len(otherStr.Value))

	for index, r := range n.Value {
		newStr[index] = r
	}
	for index, r := range value.(*String).Value {
		newStr[thisLen+index] = r
	}

	return NewString(newStr), nil

}

func (n *String) Subtract(value Value) (Value, *shared.Error) {

	otherStr := value.(*String)
	firstCopyStrLen := len(n.Value)

	for i := len(n.Value) - 1; i >= 0; i-- {

		char := n.Value[i]

		if char != ' ' {
			break
		}

		firstCopyStrLen--

	}

	newStr := make([]rune, firstCopyStrLen+len(otherStr.Value))

	for i := 0; i < firstCopyStrLen; i++ {
		newStr[i] = n.Value[i]
	}

	for i := 0; i < len(otherStr.Value); i++ {
		newStr[firstCopyStrLen+i] = otherStr.Value[i]
	}

	return NewString(newStr), nil
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

	rightString := value.(*String)

	if len(rightString.Value) > len(n.Value) {
		return NewBoolean(false), nil
	}

	compareLen := len(rightString.Value)

	for i := 0; i < compareLen; i++ {
		if rightString.Value[i] != n.Value[i] {
			return NewBoolean(false), nil
		}
	}

	return NewBoolean(true), nil
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

	return NewBoolean(string(n.Value) > string(value.(*String).Value)), nil
}

func (n *String) IsGreaterEquals(value Value) (Value, *shared.Error) {

	return NewBoolean(string(n.Value) >= string(value.(*String).Value)), nil
}

func (n *String) IsLess(value Value) (Value, *shared.Error) {

	return NewBoolean(string(n.Value) < string(value.(*String).Value)), nil
}

func (n *String) IsLessEquals(value Value) (Value, *shared.Error) {

	return NewBoolean(string(n.Value) <= string(value.(*String).Value)), nil
}

func (n *String) EqualsEquals(value Value) (Value, *shared.Error) {

	rightString := value.(*String)

	return NewBoolean(string(n.Value) == string(rightString.Value)), nil
}

func (n *String) IsTrue() (bool, *shared.Error) {
	return false, shared.NewRuntimeErrorRange(n.Range, "Strings cannot be used as boolean")
}

func (n *String) IsEmpty() bool {
	return len(n.Value) == 0
}

func (n *String) UpdateRange(valueRange *shared.Range) Value {
	n.Range = valueRange
	return n
}

func NewString(value []rune) *String {
	return &String{
		Value: value,
		Range: nil,
	}
}
