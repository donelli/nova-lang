package interpreter

import "recital_lsp/pkg/shared"

type ValueType uint8

const (
	ValueType_Number ValueType = iota + 1
	ValueType_Boolean
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
}
