package interpreter

import "nova-lang/pkg/shared"

type ValueType uint8

const (
	ValueType_Number ValueType = iota + 1
	ValueType_Boolean
	ValueType_String
)

//go:generate stringer -type=ValueType -trimprefix=ValueType_

type Value interface {
	Copy() Value
	Type() ValueType
	UpdateRange(valueRange *shared.Range) Value
	PrintRepresentation() string
	Add(Value) (Value, *shared.Error)
	Subtract(Value) (Value, *shared.Error)
	Multiply(Value) (Value, *shared.Error)
	Divide(Value) (Value, *shared.Error)
	Exponential(Value) (Value, *shared.Error)
	Remainder(Value) (Value, *shared.Error)
	Equals(Value) (Value, *shared.Error)
	EqualsEquals(Value) (Value, *shared.Error)
	NotEquals(Value) (Value, *shared.Error)
	IsGreater(Value) (Value, *shared.Error)
	IsGreaterEquals(Value) (Value, *shared.Error)
	IsLess(Value) (Value, *shared.Error)
	IsLessEquals(Value) (Value, *shared.Error)
	And(Value) (Value, *shared.Error)
	Or(Value) (Value, *shared.Error)
	IsTrue() (bool, *shared.Error)
	IsEmpty() bool
}
